package ctx

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/kgrunwald/goweb/soap"
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

	// ContentType returns the Content-Type of the given request
	ContentType() string

	Respond(int, interface{}) error
	OK(interface{}) error
	NotFound(interface{}) error
	Unauthorized(interface{}) error
	Forbidden(interface{}) error
	BadRequest(interface{}) error
}

func New(r *http.Request, w http.ResponseWriter) Context {
	c := &ctx{
		req:    r,
		writer: w,
	}

	if c.ContentType() == ContentTypeXML {
		if c.req.Header.Get("SOAPAction") != "" {
			c.decoder = soap.NewDecoder(c.req.Body)
			c.encoder = soap.NewEncoder(c.writer)
		} else {
			c.decoder = xml.NewDecoder(c.req.Body)
			c.encoder = xml.NewEncoder(c.writer)
		}
	} else {
		c.decoder = json.NewDecoder(c.req.Body)
		c.encoder = json.NewEncoder(c.writer)
	}

	return c
}

type Encoder interface {
	Encode(interface{}) error
}

type Decoder interface {
	Decode(interface{}) error
}

type ctx struct {
	req     *http.Request
	writer  http.ResponseWriter
	decoder Decoder
	encoder Encoder
}

func (c *ctx) Request() *http.Request {
	return c.req
}

func (c *ctx) Bind(out interface{}) error {
	return c.decoder.Decode(out)
}

func (c *ctx) ContentType() string {
	accept := c.req.Header.Get("Accept")
	if accept == ContentTypeXML {
		return ContentTypeXML
	}

	return ContentTypeJSON
}

func (c *ctx) Respond(status int, body interface{}) error {
	c.writer.Header().Set("Content-Type", c.ContentType())
	c.writer.WriteHeader(status)
	return c.encoder.Encode(body)
}

// OK is a helper method that returns a response with a 200 status code
func (c *ctx) OK(body interface{}) error {
	return c.Respond(http.StatusOK, body)
}

// NotFound is a helper method that returns a response with a 404 status code
func (c *ctx) NotFound(body interface{}) error {
	return c.Respond(http.StatusNotFound, body)
}

// Forbidden is a helper method that returns a response with a 403 status code
func (c *ctx) Forbidden(body interface{}) error {
	return c.Respond(http.StatusForbidden, body)
}

// Unauthorized is a helper method that returns a response with a 401 status code
func (c *ctx) Unauthorized(body interface{}) error {
	return c.Respond(http.StatusUnauthorized, body)
}

// BadRequest is a helper method that returns a response with a 400 status code
func (c *ctx) BadRequest(body interface{}) error {
	return c.Respond(http.StatusBadRequest, body)
}
