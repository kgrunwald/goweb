package soap

import (
	"encoding/xml"
)

type SOAP struct {
	Envelope   Envelope
	Namespaces []Namespace
}

type Envelope struct {
	XMLName xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	Body    Body     `xml:"http://www.w3.org/2003/05/soap-envelope Body"`
}

type Body struct {
	Value []byte `xml:",innerxml"`
}

type Fault struct {
	XMLName xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Fault"`
	Code   FaultCode   `xml:"http://www.w3.org/2003/05/soap-envelope Code"`
	Reason FaultReason `xml:"http://www.w3.org/2003/05/soap-envelope Reason"`
	Detail interface{} `xml:"http://www.w3.org/2003/05/soap-envelope Detail,omitempty"`
}

type FaultCode struct {
	Value FaultCodeValue `xml:"http://www.w3.org/2003/05/soap-envelope Value"`
}

type FaultCodeValue string

const (
	FaultCodeReceiver FaultCodeValue = "SOAP-ENV:Receiver"
	FaultCodeSender   FaultCodeValue = "SOAP-ENV:Sender"
)

type FaultReason struct {
	Text []ReasonText `xml:"http://www.w3.org/2003/05/soap-envelope Text"`
}

type ReasonText struct {
	Text string `xml:",innerxml"`
	Lang string `xml:"xml:lang,attr"`
}

type Namespace struct {
	Prefix string
	URL    string
}

func (s *SOAP) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Space == "xmlns" {
			s.Namespaces = append(s.Namespaces, Namespace{attr.Name.Local, attr.Value})
		}
	}

	return d.DecodeElement(&s.Envelope, &start)
}
