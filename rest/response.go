package rest

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

const HeaderContentType = "Content-Type"
const ContentTypeJson = "application/json"
const ContentTypeXml = "application/xml"

type Response struct {
	StatusCode  int
	Body        interface{}
	Marshaller  func(interface{}) ([]byte, error)
	ContentType string
}

type errorResponse struct {
	Error string `xml:"error" json:"error"`
}

func NewResponse(r *http.Request, body interface{}) *Response {
	accept := r.Header.Get("Accept")
	marshaller := json.Marshal
	contentType := ContentTypeJson

	if accept == ContentTypeXml {
		marshaller = xml.Marshal
		contentType = ContentTypeXml
	}

	return &Response{
		StatusCode:  http.StatusOK,
		Body:        body,
		Marshaller:  marshaller,
		ContentType: contentType,
	}
}

func (r *Response) WithStatus(status int) *Response {
	r.StatusCode = status
	return r
}

func (r *Response) Send(w http.ResponseWriter) error {
	bodyBytes, err := r.Marshaller(r.Body)
	if err != nil {
		bodyBytes, err = r.Marshaller(errorResponse{Error: err.Error()})
		if err != nil {
			return err
		}
		r.StatusCode = http.StatusInternalServerError
	}

	w.Header().Set(HeaderContentType, r.ContentType)
	w.WriteHeader(r.StatusCode)
	w.Write(bodyBytes)
	return nil
}
