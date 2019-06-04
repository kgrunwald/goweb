package ctx

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type soapCtx struct {
	*responseBuilder
	req *http.Request
}

func newSoapContext(req *http.Request, w http.ResponseWriter) Context {
	ctx := &soapCtx{
		responseBuilder: &responseBuilder{
			Writer:  w,
			Encoder: xml.NewEncoder(w),
		},
		req: req,
	}

	w.Header().Set("Content-Type", ctx.ContentType())

	return ctx
}

func (c *soapCtx) Request() *http.Request {
	return c.req
}

func (c *soapCtx) Bind(out interface{}) error {
	soapObj := &SOAP{}
	if err := xml.NewDecoder(c.req.Body).Decode(soapObj); err != nil {
		return err
	}

	soapBody := soapObj.Envelope.Body.Value
	soapBody = InjectNamespaces(soapBody, soapObj.Namespaces)
	return xml.Unmarshal(soapBody, out)
}

func (c *soapCtx) Marshal(res interface{}) ([]byte, error) {
	return MarshalSoap(res)
}

func MarshalSoap(res interface{}) ([]byte, error) {
	env := Envelope{}
	body, err := xml.Marshal(res)
	if err != nil {
		return []byte{}, err
	}

	env.Body = Body{body}

	return xml.Marshal(env)
}

func (c *soapCtx) ContentType() string {
	return ContentTypeXML
}

func InjectNamespaces(body []byte, namespaces []Namespace) []byte {
	i := 0
	for _, char := range body {
		if char == '/' || char == '>' {
			break
		}
		i++
	}

	nsBody := body
	for _, ns := range namespaces {
		nsBytes := fmt.Sprintf(` xmlns:%s="%s"`, ns.Prefix, ns.URL)
		nsBody = append(nsBody[:i], append([]byte(nsBytes), nsBody[i:]...)...)
	}

	return nsBody
}
