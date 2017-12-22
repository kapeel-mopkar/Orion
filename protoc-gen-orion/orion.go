// protoc-gen-go is a plugin for the Google protocol buffer compiler to generate
// Go code.  Run it by building this program and putting it in your path with
// the name
// 	protoc-gen-go
// That word 'go' at the end becomes part of the option string set for the
// protocol compiler, so once the protocol compiler (protoc) is installed
// you can run
// 	protoc --go_out=output_directory input_directory/file.proto
// to generate Go bindings for the protocol defined by file.proto.
// With that input, the output will be written to
// 	output_directory/file.pb.go
//
// The generated code is documented in the package comment for
// the library.
//
// See the README and documentation for protocol buffers to learn more:
// 	https://developers.google.com/protocol-buffers/
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

// P prints the arguments to the generated output.  It handles strings and int32s, plus
// handling indirections because they may be *string, etc.
func P(g *generator.Generator, str ...interface{}) {
	//g.WriteString(g.indent)
	for _, v := range str {
		switch s := v.(type) {
		case string:
			g.WriteString(s)
		case *string:
			g.WriteString(*s)
		case bool:
			fmt.Fprintf(g, "%t", s)
		case *bool:
			fmt.Fprintf(g, "%t", *s)
		case int:
			fmt.Fprintf(g, "%d", s)
		case *int32:
			fmt.Fprintf(g, "%d", *s)
		case *int64:
			fmt.Fprintf(g, "%d", *s)
		case float64:
			fmt.Fprintf(g, "%g", s)
		case *float64:
			fmt.Fprintf(g, "%g", *s)
		default:
			g.Fail(fmt.Sprintf("unknown type in printer: %T", v))
		}
	}
	g.WriteByte('\n')
}

// Generate the package definition
func generateHeader(g *generator.Generator, file *descriptor.FileDescriptorProto) {
	P(g, "// Code generated by protoc-gen-orion. DO NOT EDIT.")
	P(g, "// source: ", file.Name)
	P(g)

	name := file.GetPackage()

	P(g, "package ", name)
	P(g, "")
	P(g, "import (")
	P(g, "\torion \"github.com/carousell/Orion/orion\"")
	P(g, ")")
}

// Generate the file
func generate(g *generator.Generator, file *descriptor.FileDescriptorProto) {
	generateHeader(g, file)
	for _, svc := range file.GetService() {
		origServName := svc.GetName()
		fullServName := origServName
		if pkg := file.GetPackage(); pkg != "" {
			fullServName = pkg + "." + fullServName
		}
		servName := generator.CamelCase(origServName)
		serviceDescVar := "_" + servName + "_serviceDesc"

		P(g)
		P(g, "func Register", servName, "OrionServer(srv orion.ServiceFactory, orionServer orion.Server) {")
		P(g, "\torionServer.RegisterService(&", serviceDescVar, `, srv)`)
		P(g, "}")
	}
}

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
				} else {
					return &commentsInfo{
						Decoder: true,
					}
				}
			}
		}
	}
	return nil
}

func generateCustomURL(g *generator.Generator, file *descriptor.FileDescriptorProto) {
	comments := extractComments(file)
	for index, svc := range file.GetService() {
		path := fmt.Sprintf("6,%d", index) // 6 means service.
		for i, method := range svc.GetMethod() {
			commentPath := fmt.Sprintf("%s,2,%d", path, i) // 2 means method in a service.
			if loc, ok := comments[commentPath]; ok {
				text := strings.TrimSuffix(loc.GetLeadingComments(), "\n")
				for _, line := range strings.Split(text, "\n") {
					if option := parseComments(line); option != nil {
						if option.Encoder {
							optionsEncoder := false
							methods := strings.Split(option.Method, "/")
							for i := range methods {
								if strings.ToLower(methods[i]) == "options" {
									optionsEncoder = true
								}
								methods[i] = "\"" + methods[i] + "\""
							}
							methodsString := strings.Join(methods, ",")
							P(g, "")
							P(g, "func Register", svc.GetName(), method.GetName(), "Encoder(svr orion.Server, encoder orion.Encoder) {")
							P(g, "\torion.RegisterEncoders(svr, \""+svc.GetName()+"\", \""+method.GetName()+"\", []string{"+methodsString+"}, \""+option.Path+"\", encoder)")
							P(g, "}")
							if optionsEncoder {
								P(g, "")
								P(g, "func Register", svc.GetName(), method.GetName(), "Handler(svr orion.Server, handler orion.HTTPHandler) {")
								P(g, "\torion.RegisterHandler(svr, \""+svc.GetName()+"\", \""+method.GetName()+"\", \""+option.Path+"\", handler)")
								P(g, "}")
							}
						}
						if option.Decoder {
							P(g, "")
							P(g, "func Register", svc.GetName(), method.GetName(), "Decoder(svr orion.Server, decoder orion.Decoder) {")
							P(g, "\torion.RegisterDecoder(svr, \""+svc.GetName()+"\", \""+method.GetName()+"\", decoder)")
							P(g, "}")
						}
					}
				}
			}
		}
	}
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

func main() {
	g := generator.New()

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		g.Error(err, "reading input")
	}

	if err := proto.Unmarshal(data, g.Request); err != nil {
		g.Error(err, "parsing input proto")
	}

	if len(g.Request.FileToGenerate) == 0 {
		g.Fail("no files to generate")
	}

	for _, file := range g.Request.GetProtoFile() {
		g.Reset()
		if len(file.Service) > 0 {
			generate(g, file)
			generateCustomURL(g, file)
			g.Response.File = append(g.Response.File, &plugin.CodeGeneratorResponse_File{
				Name:    proto.String(strings.ToLower(file.GetName()) + ".orion.pb.go"),
				Content: proto.String(g.String()),
			})
		}
	}

	// Send back the results.
	data, err = proto.Marshal(g.Response)
	if err != nil {
		g.Error(err, "failed to marshal output proto")
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		g.Error(err, "failed to write output proto")
	}
}
