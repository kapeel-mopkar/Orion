package apm

import (
	"fmt"
	"net/http"

	"go.elastic.co/apm/internal/apmhttputil"
	"go.elastic.co/apm/model"
)

// Context provides methods for setting transaction and error context.
type Context struct {
	model            model.Context
	request          model.Request
	requestBody      model.RequestBody
	requestSocket    model.RequestSocket
	response         model.Response
	user             model.User
	service          model.Service
	serviceFramework model.Framework
	captureBodyMask  CaptureBodyMode
}

func (c *Context) build() *model.Context {
	switch {
	case c.model.Request != nil:
	case c.model.Response != nil:
	case c.model.User != nil:
	case c.model.Service != nil:
	case len(c.model.Custom) != 0:
	case len(c.model.Tags) != 0:
	default:
		return nil
	}
	return &c.model
}

func (c *Context) reset() {
	modelContext := model.Context{
		// TODO(axw) reuse space for tags
		Custom: c.model.Custom[:0],
	}
	*c = Context{
		model:           modelContext,
		captureBodyMask: c.captureBodyMask,
		request: model.Request{
			Headers: c.request.Headers[:0],
		},
		response: model.Response{
			Headers: c.response.Headers[:0],
		},
	}
}

// SetCustom sets a custom context key/value pair. If the key is invalid
// (contains '.', '*', or '"'), the call is a no-op. The value must be
// JSON-encodable.
//
// If value implements interface{AppendJSON([]byte) []byte}, that will be
// used to encode the value. Otherwise, value will be encoded using
// json.Marshal. As a special case, values of type map[string]interface{}
// will be traversed and values encoded according to the same rules.
func (c *Context) SetCustom(key string, value interface{}) {
	if !validTagKey(key) {
		return
	}
	c.model.Custom.Set(key, value)
}

// SetTag sets a tag in the context. If the key is invalid
// (contains '.', '*', or '"'), the call is a no-op.
func (c *Context) SetTag(key, value string) {
	if !validTagKey(key) {
		return
	}
	value = truncateString(value)
	if c.model.Tags == nil {
		c.model.Tags = map[string]string{key: value}
	} else {
		c.model.Tags[key] = value
	}
}

// SetFramework sets the framework name and version in the context.
//
// This is used for identifying the framework in which the context
// was created, such as Gin or Echo.
//
// If the name is empty, this is a no-op. If version is empty, then
// it will be set to "unspecified".
func (c *Context) SetFramework(name, version string) {
	if name == "" {
		return
	}
	if version == "" {
		// Framework version is required.
		version = "unspecified"
	}
	c.serviceFramework = model.Framework{
		Name:    truncateString(name),
		Version: truncateString(version),
	}
	c.service.Framework = &c.serviceFramework
	c.model.Service = &c.service
}

// SetHTTPRequest sets details of the HTTP request in the context.
//
// This function relates to server-side requests. Various proxy
// forwarding headers are taken into account to reconstruct the URL,
// and determining the client address.
//
// If the request URL contains user info, it will be removed and
// excluded from the URL's "full" field.
//
// If the request contains HTTP Basic Authentication, the username
// from that will be recorded in the context. Otherwise, if the
// request contains user info in the URL (i.e. a client-side URL),
// that will be used.
func (c *Context) SetHTTPRequest(req *http.Request) {
	// Special cases to avoid calling into fmt.Sprintf in most cases.
	var httpVersion string
	switch {
	case req.ProtoMajor == 1 && req.ProtoMinor == 1:
		httpVersion = "1.1"
	case req.ProtoMajor == 2 && req.ProtoMinor == 0:
		httpVersion = "2.0"
	default:
		httpVersion = fmt.Sprintf("%d.%d", req.ProtoMajor, req.ProtoMinor)
	}

	var forwarded *apmhttputil.ForwardedHeader
	if fwd := req.Header.Get("Forwarded"); fwd != "" {
		parsed := apmhttputil.ParseForwarded(fwd)
		forwarded = &parsed
	}
	c.request = model.Request{
		Body:        c.request.Body,
		URL:         apmhttputil.RequestURL(req, forwarded),
		Method:      truncateString(req.Method),
		HTTPVersion: httpVersion,
		Cookies:     req.Cookies(),
	}
	c.model.Request = &c.request

	for k, values := range req.Header {
		if k == "Cookie" {
			// We capture cookies in the request structure.
			continue
		}
		c.request.Headers = append(c.request.Headers, model.Header{
			Key: k, Values: values,
		})
	}

	c.requestSocket = model.RequestSocket{
		Encrypted:     req.TLS != nil,
		RemoteAddress: apmhttputil.RemoteAddr(req, forwarded),
	}
	if c.requestSocket != (model.RequestSocket{}) {
		c.request.Socket = &c.requestSocket
	}

	username, _, ok := req.BasicAuth()
	if !ok && req.URL.User != nil {
		username = req.URL.User.Username()
	}
	c.user.Username = truncateString(username)
	if c.user.Username != "" {
		c.model.User = &c.user
	}
}

// SetHTTPRequestBody sets the request body in context given a (possibly nil)
// BodyCapturer returned by Tracer.CaptureHTTPRequestBody.
func (c *Context) SetHTTPRequestBody(bc *BodyCapturer) {
	if bc == nil || bc.captureBody&c.captureBodyMask == 0 {
		return
	}
	if bc.setContext(&c.requestBody) {
		c.request.Body = &c.requestBody
	}
}

// SetHTTPResponseHeaders sets the HTTP response headers in the context.
func (c *Context) SetHTTPResponseHeaders(h http.Header) {
	for k, values := range h {
		c.response.Headers = append(c.response.Headers, model.Header{
			Key: k, Values: values,
		})
	}
	if len(c.response.Headers) != 0 {
		c.model.Response = &c.response
	}
}

// SetHTTPStatusCode records the HTTP response status code.
func (c *Context) SetHTTPStatusCode(statusCode int) {
	c.response.StatusCode = statusCode
	c.model.Response = &c.response
}

// SetUserID sets the ID of the authenticated user.
func (c *Context) SetUserID(id string) {
	c.user.ID = truncateString(id)
	if c.user.ID != "" {
		c.model.User = &c.user
	}
}

// SetUserEmail sets the email for the authenticated user.
func (c *Context) SetUserEmail(email string) {
	c.user.Email = truncateString(email)
	if c.user.Email != "" {
		c.model.User = &c.user
	}
}

// SetUsername sets the username of the authenticated user.
func (c *Context) SetUsername(username string) {
	c.user.Username = truncateString(username)
	if c.user.Username != "" {
		c.model.User = &c.user
	}
}
