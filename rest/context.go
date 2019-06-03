package rest

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

// HeaderContentType holds the name of the Content-Type HTTP header
const HeaderContentType = "Content-Type"

// HeaderAccept holds the name of the Accept HTTP header
const HeaderAccept = "Accept"

// ContentTypeJSON holds the application/json content type value
const ContentTypeJSON = "application/json"

// ContentTypeXML holds the application/xml content type value
const ContentTypeXML = "application/xml"

// Context provides functions to help with the lifecycle of a particular request
type Context interface {
	// Request returns the underlying HTTP Request for this Context
	Request() *http.Request

	// Bind reads the body of the HTTP request and deserializes it into the provided interface. The Content-Type header will be used
	// to determine the encoding of the request to deserialize the body. If no Content-Type header is provided,
	// application/json will be assumed.
	Bind(interface{}) error

	// Marshal will serialize the provided response based on the Accept header set in the Request
	Marshal(*Response) ([]byte, error)

	// ContentType returns the Content-Type of the given request
	ContentType() string
}

type ctx struct {
	req          *http.Request
	marshaller   func(interface{}) ([]byte, error)
	unmarshaller func([]byte, interface{}) error
}

// NewContext returns a concrete implementation of the Context interface
func NewContext(req *http.Request) Context {
	ctx := &ctx{
		req:          req,
		marshaller:   json.Marshal,
		unmarshaller: json.Unmarshal,
	}

	if ctx.ContentType() == ContentTypeXML {
		ctx.marshaller = xml.Marshal
		ctx.unmarshaller = xml.Unmarshal
	}

	return ctx
}

func (c *ctx) Request() *http.Request {
	return c.req
}

func (c *ctx) Bind(out interface{}) error {
	body, err := ioutil.ReadAll(c.req.Body)
	if err != nil {
		return err
	}

	return c.unmarshaller(body, out)
}

func (c *ctx) Marshal(res *Response) ([]byte, error) {
	return c.marshaller(res.Body)
}

func (c *ctx) ContentType() string {
	accept := c.req.Header.Get("Accept")
	if accept == ContentTypeXML {
		return ContentTypeXML
	}

	return ContentTypeJSON
}
