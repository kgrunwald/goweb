package ctx

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/kgrunwald/goweb/ilog"
	"github.com/kgrunwald/goweb/soap"
)

// HeaderContentType holds the name of the Content-Type HTTP header
const HeaderContentType = "Content-Type"

// HeaderAccept holds the name of the Accept HTTP header
const HeaderAccept = "Accept"

// HeaderSOAPAction defines the SOAPAction HTTP header
const HeaderSOAPAction = "Soapaction"

// ContentTypeJSON holds the application/json content type value
const ContentTypeJSON = "application/json"

// ContentTypeXML holds the application/xml content type value
const ContentTypeXML = "application/xml"

// ContentTypeTextXML is the content type required for SOAP 1.2 messages
const ContentTypeTextXML = "text/xml"

// Context provides functions to help with the lifecycle of a particular request
type Context interface {
	// Request returns the underlying HTTP Request for this Context
	Request() *http.Request
	RequestID() string

	// Bind reads the body of the HTTP request and deserializes it into the provided interface. The Content-Type header will be used
	// to determine the encoding of the request to deserialize the body. If no Content-Type header is provided,
	// application/json will be assumed.
	Bind(interface{}) error

	// ContentType returns the Content-Type of the given request
	ContentType() string
	Accept() string

	Respond(int, interface{}) error
	OK(interface{}) error
	NotFound(interface{}) error
	Unauthorized(interface{}) error
	Forbidden(interface{}) error
	BadRequest(interface{}) error
	SendError(error) error
	Log() ilog.Logger
}

func New(r *http.Request, w http.ResponseWriter, log ilog.Logger) Context {
	c := &ctx{
		req:       r,
		writer:    w,
		requestId: uuid.New().String(),
		log:       log,
	}

	c.Initialize()

	return c
}

type Encoder interface {
	Encode(interface{}) error
}

type Decoder interface {
	Decode(interface{}) error
}

type xmlEncoder struct {
	Writer io.Writer
}

func (x *xmlEncoder) Encode(out interface{}) error {
	if _, err := x.Writer.Write([]byte(xml.Header)); err != nil {
		return err
	}

	return xml.NewEncoder(x.Writer).Encode(out)
}

type ctx struct {
	req          *http.Request
	writer       http.ResponseWriter
	log          ilog.Logger
	requestId    string
	responseType string
	encoder      Encoder
	decoder      Decoder
}

type ErrorMessage struct {
	XMLName xml.Name `xml:"error" json:"-"`
	Message string   `xml:",innerxml" json:"error"`
}

func (e *ErrorMessage) Error() string {
	return e.Message
}

func (c *ctx) Request() *http.Request {
	return c.req
}

func (c *ctx) Bind(out interface{}) error {
	return c.decoder.Decode(out)
}

func (c *ctx) ContentType() string {
	contentType := c.req.Header.Get(HeaderContentType)
	if contentType != "" {
		return contentType
	}

	return ContentTypeJSON
}

func (c *ctx) Accept() string {
	accept := c.req.Header.Get(HeaderAccept)
	if accept != "" {
		return accept
	}

	return c.ContentType()
}

func (c *ctx) IsSOAP() bool {
	contentType := c.ContentType()
	return contentType == ContentTypeTextXML && (len(c.req.Header[HeaderSOAPAction]) > 0)
}

func (c *ctx) IsXML() bool {
	contentType := c.ContentType()
	return contentType == ContentTypeTextXML || contentType == ContentTypeXML
}

func (c *ctx) Initialize() {
	if c.IsSOAP() {
		c.encoder = soap.NewEncoder(c.writer)
		c.decoder = soap.NewDecoder(c.req.Body)
		c.responseType = ContentTypeTextXML
		return
	}

	contentType := c.ContentType()
	if contentType == ContentTypeXML || contentType == ContentTypeTextXML {
		c.decoder = xml.NewDecoder(c.req.Body)
	} else {
		c.decoder = json.NewDecoder(c.req.Body)
	}

	accept := c.Accept()
	if accept == ContentTypeXML || accept == ContentTypeTextXML {
		c.encoder = &xmlEncoder{c.writer}
		c.responseType = accept
		return
	}

	if accept == "" && c.IsXML() {
		c.encoder = &xmlEncoder{c.writer}
		c.responseType = c.ContentType()
		return
	}

	c.encoder = json.NewEncoder(c.writer)
	c.responseType = ContentTypeJSON
}

func (c *ctx) Respond(status int, body interface{}) error {
	c.writer.Header().Set("RequestID", c.RequestID())
	c.writer.Header().Set(HeaderContentType, c.responseType)
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

func (c *ctx) SendError(err error) error {
	c.Log().Error(err.Error())
	return c.Respond(http.StatusInternalServerError, &ErrorMessage{Message: err.Error()})
}

func (c *ctx) RequestID() string {
	return c.requestId
}

func (c *ctx) Log() ilog.Logger {
	return c.log.WithField("RequestID", c.RequestID())
}
