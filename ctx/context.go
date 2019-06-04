package ctx

import (
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
	if r.Header.Get("SOAPAction") != "" {
		return newSoapContext(r, w)
	}

	return newRestContext(r, w)
}

type Encoder interface {
	Encode(interface{}) error
}

type Decoder interface {
	Decode(interface{}) error
}

type responseBuilder struct {
	Writer  http.ResponseWriter
	Encoder Encoder
}

func (b *responseBuilder) Respond(status int, body interface{}) error {
	b.Writer.WriteHeader(status)
	return b.Encoder.Encode(body)
}

// OK is a helper method that returns a response with a 200 status code
func (b *responseBuilder) OK(body interface{}) error {
	return b.Respond(http.StatusOK, body)
}

// NotFound is a helper method that returns a response with a 404 status code
func (b *responseBuilder) NotFound(body interface{}) error {
	return b.Respond(http.StatusNotFound, body)
}

// Forbidden is a helper method that returns a response with a 403 status code
func (b *responseBuilder) Forbidden(body interface{}) error {
	return b.Respond(http.StatusForbidden, body)
}

// Unauthorized is a helper method that returns a response with a 401 status code
func (b *responseBuilder) Unauthorized(body interface{}) error {
	return b.Respond(http.StatusUnauthorized, body)
}

// BadRequest is a helper method that returns a response with a 400 status code
func (b *responseBuilder) BadRequest(body interface{}) error {
	return b.Respond(http.StatusBadRequest, body)
}
