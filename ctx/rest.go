package ctx

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

type restCtx struct {
	*responseBuilder

	req     *http.Request
	decoder Decoder
	writer  http.ResponseWriter
}

func newRestContext(req *http.Request, w http.ResponseWriter) Context {
	c := &restCtx{
		responseBuilder: &responseBuilder{
			Writer:  w,
			Encoder: json.NewEncoder(w),
		},
		req:     req,
		writer:  w,
		decoder: json.NewDecoder(req.Body),
	}

	if c.ContentType() == ContentTypeXML {
		c.decoder = xml.NewDecoder(req.Body)
	}

	w.Header().Set("Content-Type", c.ContentType())

	return c
}

func (c *restCtx) Request() *http.Request {
	return c.req
}

func (c *restCtx) Bind(out interface{}) error {
	return c.decoder.Decode(out)
}

func (c *restCtx) ContentType() string {
	accept := c.req.Header.Get("Accept")
	if accept == ContentTypeXML {
		return ContentTypeXML
	}

	return ContentTypeJSON
}
