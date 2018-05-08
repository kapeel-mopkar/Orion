// protoc-gen-orion is a plugin for the Google protocol buffer compiler to generate
// Orion Go code.  Run it by building this program and putting it in your path with
// the name
// 	protoc-gen-orion
//
// The generated code is documented in the package comment for
// the library.
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/micro/protobuf/protoc-gen-go/generator"
)

const (
	ORION   = "ORION"
	URL     = "URL"
	DELIM   = ":"
	DECODER = "DECODER"
)

type commentsInfo struct {
	Method  string
	Path    string
	Decoder bool
	Encoder bool
}

type data struct {
	FileName    string
	PackageName string
	Services    []*service
}

type service struct {
	ServName       string
	ServiceDescVar string
	Encoders       []*encoder
	Decoders       []*decoder
	Handlers       []*handler
}

type encoder struct {
	SvcName    string
	MethodName string
	Path       string
	Methods    string
}
type decoder struct {
	SvcName    string
	MethodName string
}
type handler struct {
	SvcName    string
	MethodName string
	Path       string
}

var tmpl = `// Code generated by protoc-gen-orion. DO NOT EDIT.
// source: {{ .FileName }}

package {{ .PackageName }}
{{ if .Services }}
import (
	orion "github.com/carousell/Orion/orion"
)

// If you see error please update your orion-protoc-gen by running 'go get -u github.com/carousell/Orion/protoc-gen-orion'
var _ = orion.ProtoGenVersion1_0
{{ end }}
{{ range .Services -}}
// Encoders
{{ range .Encoders }}
// Register{{.SvcName}}{{.MethodName}}Encoder registers the encoder for {{.MethodName}} method in {{.SvcName}}
// it registers HTTP {{ if .Path }} path {{.Path}} {{ end }}with {{.Methods}} methods
func Register{{.SvcName}}{{.MethodName}}Encoder(svr orion.Server, encoder orion.Encoder) {
	orion.RegisterEncoders(svr, "{{.SvcName}}", "{{.MethodName}}", []string{ {{- .Methods -}} }, "{{.Path}}", encoder)
}
{{ end }}
// Handlers
{{ range .Handlers }}
// Register{{.SvcName}}{{.MethodName}}Handler registers the handler for {{.MethodName}} method in {{.SvcName}}
func Register{{.SvcName}}{{.MethodName}}Handler(svr orion.Server, handler orion.HTTPHandler) {
	orion.RegisterHandler(svr, "{{.SvcName}}", "{{.MethodName}}", "{{.Path}}", handler)
}
{{ end }}
// Decoders
{{ range .Decoders }}
// Register{{.SvcName}}{{.MethodName}}Decoder registers the decoder for {{.MethodName}} method in {{.SvcName}}
func Register{{.SvcName}}{{.MethodName}}Decoder(svr orion.Server, decoder orion.Decoder) {
	orion.RegisterDecoder(svr, "{{.SvcName}}", "{{.MethodName}}", decoder)
}
{{ end }}
// Register{{.ServName}}OrionServer registers {{.ServName}} to Orion server
func Register{{.ServName}}OrionServer(srv orion.ServiceFactory, orionServer orion.Server) {
	orionServer.RegisterService(&{{.ServiceDescVar}}, srv)
{{ range .Encoders }}
	Register{{.SvcName}}{{.MethodName}}Encoder(orionServer, nil)
{{- end }}
}{{ end }}
`

// Error reports a problem, including an error, and exits the program.
func logError(err error, msgs ...string) {
	s := strings.Join(msgs, " ") + ":" + err.Error()
	log.Print("protoc-gen-orion: error:", s)
	os.Exit(1)
}

// Fail reports a problem and exits the program.
func logFail(msgs ...string) {
	s := strings.Join(msgs, " ")
	log.Print("protoc-gen-orion: error:", s)
	os.Exit(1)
}

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		logError(err, "reading input")
	}

	request := new(plugin.CodeGeneratorRequest)
	if err := proto.Unmarshal(data, request); err != nil {
		logError(err, "parsing input proto")
	}

	if len(request.FileToGenerate) == 0 {
		logFail("no files to generate")
	}

	response := new(plugin.CodeGeneratorResponse)
	response.File = make([]*plugin.CodeGeneratorResponse_File, 0)

	for _, file := range request.GetProtoFile() {
		// check if file has any service
		if len(file.Service) > 0 {
			f := generateFile(populate(file))
			response.File = append(response.File, f)
		}
	}

	// Send back the results.
	data, err = proto.Marshal(response)
	if err != nil {
		logError(err, "failed to marshal output proto")
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		logError(err, "failed to write output proto")
	}
}

func generateFile(d *data) *plugin.CodeGeneratorResponse_File {
	t := template.New("file")
	t, err := t.Parse(tmpl)
	if err != nil {
		logError(err, "failed parsing template")
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, d)
	if err != nil {
		logError(err, "failed parsing template")
	}

	file := new(plugin.CodeGeneratorResponse_File)
	file.Content = proto.String(buf.String())
	file.Name = proto.String(strings.ToLower(d.FileName) + ".orion.pb.go")
	return file
}

func populate(file *descriptor.FileDescriptorProto) *data {
	d := new(data)
	d.FileName = *file.Name
	d.PackageName = file.GetPackage()

	d.Services = make([]*service, 0)
	generate(d, file)

	return d
}

func generate(d *data, file *descriptor.FileDescriptorProto) {
	comments := extractComments(file)
	for index, svc := range file.GetService() {

		origServName := svc.GetName()
		servName := generator.CamelCase(origServName) // use the same logic from go-grpc generator
		serviceDescVar := "_" + servName + "_serviceDesc"

		s := new(service)
		s.Encoders = make([]*encoder, 0)
		s.Handlers = make([]*handler, 0)
		s.Decoders = make([]*decoder, 0)
		s.ServiceDescVar = serviceDescVar
		s.ServName = servName
		d.Services = append(d.Services, s)

		// ** --- START -- Find comments in grpc services
		path := fmt.Sprintf("6,%d", index) // 6 means service.
		for i, method := range svc.GetMethod() {
			commentPath := fmt.Sprintf("%s,2,%d", path, i) // 2 means method in a service.
			if loc, ok := comments[commentPath]; ok {
				text := strings.TrimSuffix(loc.GetLeadingComments(), "\n")
				for _, line := range strings.Split(text, "\n") {
					// ** --- END -- Find comments in grpc services

					if option := parseComments(line); option != nil {
						if option.Encoder {
							methods := strings.Split(option.Method, "/")
							for i := range methods {
								if strings.ToLower(methods[i]) == "options" {
								}
								methods[i] = "\"" + methods[i] + "\""
							}
							methodsString := strings.Join(methods, ", ")

							// populate encoder
							enc := new(encoder)
							enc.SvcName = svc.GetName()
							enc.MethodName = method.GetName()
							enc.Path = option.Path
							enc.Methods = methodsString
							s.Encoders = append(s.Encoders, enc)

							// popluate handler
							han := new(handler)
							han.SvcName = svc.GetName()
							han.MethodName = method.GetName()
							han.Path = option.Path
							s.Handlers = append(s.Handlers, han)
						}

						if option.Decoder {
							// popluate decoder
							dec := new(decoder)
							dec.SvcName = svc.GetName()
							dec.MethodName = method.GetName()
							s.Decoders = append(s.Decoders, dec)
						}
					}
				}
			}
		}
	}
}

func parseComments(line string) *commentsInfo {
	parts := strings.Split(line, DELIM)
	if len(parts) > 1 {
		if ORION == strings.ToUpper(strings.TrimSpace(parts[0])) {
			switch strings.ToUpper(strings.TrimSpace(parts[1])) {
			case URL:
				if len(parts) > 2 {
					values := strings.SplitN(strings.TrimSpace(parts[2]), " ", 2)
					if len(values) == 2 {
						return &commentsInfo{
							Method:  strings.ToUpper(values[0]),
							Path:    values[1],
							Encoder: true,
							Decoder: true,
						}
					}
					return &commentsInfo{
						Method:  strings.ToUpper(values[0]),
						Path:    "",
						Encoder: true,
						Decoder: true,
					}
				}
				return &commentsInfo{
					Decoder: true,
				}
			}
		}
	}
	return nil
}

func extractComments(file *descriptor.FileDescriptorProto) map[string]*descriptor.SourceCodeInfo_Location {
	comments := make(map[string]*descriptor.SourceCodeInfo_Location)
	for _, loc := range file.GetSourceCodeInfo().GetLocation() {
		if loc.LeadingComments == nil {
			continue
		}
		var p []string
		for _, n := range loc.Path {
			p = append(p, strconv.Itoa(int(n)))
		}
		comments[strings.Join(p, ",")] = loc
	}
	return comments
}
