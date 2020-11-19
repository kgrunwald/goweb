package ctx

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/kgrunwald/goweb/apierrors"
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
	requestID() string
	Writer() http.ResponseWriter

	AddValue(interface{}, interface{})
	GetValue(interface{}) interface{}

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
		req:    r,
		writer: w,
		id:     newUID(log),
		log:    log,
	}

	c.Initialize()

	return c
}

func newUID(log ilog.Logger) string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.WithField("error", err).Fatal("Failed to create new ID")
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
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
	id           string
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
	c.writer.Header().Set("requestID", c.requestID())
	c.writer.Header().Set(HeaderContentType, c.responseType)
	c.writer.WriteHeader(status)
	if err, ok := body.(error); ok {
		body = &ErrorMessage{Message: err.Error()}
	}
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
	msg := &ErrorMessage{Message: err.Error()}
	if _, ok := err.(apierrors.BadRequestError); ok {
		return c.BadRequest(msg)
	} else if _, ok := err.(apierrors.UnauthorizedError); ok {
		return c.Unauthorized(msg)
	} else if _, ok := err.(apierrors.ForbiddenError); ok {
		return c.Forbidden(msg)
	} else if _, ok := err.(apierrors.NotFoundError); ok {
		return c.NotFound(msg)
	}

	return c.Respond(http.StatusInternalServerError, msg)
}

func (c *ctx) requestID() string {
	return c.id
}

func (c *ctx) Log() ilog.Logger {
	return c.log.WithField("requestID", c.requestID())
}

func (c *ctx) AddValue(key interface{}, value interface{}) {
	reqCtx := c.req.Context()
	newCtx := context.WithValue(reqCtx, key, value)
	c.req = c.req.WithContext(newCtx)
}

func (c *ctx) GetValue(key interface{}) interface{} {
	return c.req.Context().Value(key.(interface{}))
}

func (c *ctx) Writer() http.ResponseWriter {
	return c.writer
}
