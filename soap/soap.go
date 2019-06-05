package soap

import (
	"encoding/xml"
	"fmt"
	"io"
)

type Encoder struct {
	Writer io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

func (e *Encoder) Encode(res interface{}) error {
	if err, ok := res.(error); ok {
		res = &Fault{
			Code: FaultCode{Value: FaultCodeReceiver},
			Reason: FaultReason{
				[]ReasonText{
					ReasonText{
						Text: "Server error",
						Lang: "en",
					},
				},
			},
			Detail: err.Error(),
		}
	}
	
	env := Envelope{}
	body, err := xml.Marshal(res)
	if err != nil {
		return err
	}
	env.Body = Body{body}
	e.Writer.Write([]byte(xml.Header))
	return xml.NewEncoder(e.Writer).Encode(env)
}

type Decoder struct {
	Reader io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r}
}

func (d *Decoder) Decode(out interface{}) error {
	soapObj := &SOAP{}
	if err := xml.NewDecoder(d.Reader).Decode(soapObj); err != nil {
		return err
	}

	soapBody := soapObj.Envelope.Body.Value
	soapBody = InjectNamespaces(soapBody, soapObj.Namespaces)
	return xml.Unmarshal(soapBody, out)
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
