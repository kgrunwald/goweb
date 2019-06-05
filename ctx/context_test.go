package ctx

import (
	"bytes"
	"testing"

	"net/http/httptest"

	"github.com/golang/mock/gomock"
	"github.com/kgrunwald/goweb/ilog"
	"github.com/kgrunwald/goweb/ilog/mock_ilog"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
}

var l ilog.Logger

func TestContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	l = mock_ilog.NewMockLogger(ctrl)
	suite.Run(t, new(testSuite))
}

func (s *testSuite) TestNewContext() {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	ctx := New(req, w, l)
	s.EqualValues(ctx.Request(), req)
}

func (s *testSuite) TestBindJSON() {
	type T struct {
		X int `json:"x"`
	}

	req := httptest.NewRequest("GET", "/", bytes.NewBufferString(`{"x": 5}`))
	w := httptest.NewRecorder()
	ctx := New(req, w, l)
	t := T{}
	ctx.Bind(&t)
	s.Equal(5, t.X, "Bind did not parse JSON correctly")
}

func (s *testSuite) TestBindXML() {
	type T struct {
		X int `xml:"x"`
	}

	req := httptest.NewRequest("GET", "/", bytes.NewBufferString(`<T><x>5</x></T>`))
	w := httptest.NewRecorder()
	req.Header.Add(HeaderAccept, ContentTypeXML)
	ctx := New(req, w, l)
	t := T{}
	ctx.Bind(&t)
	s.Equal(5, t.X, "Bind did not parse XML correctly")
}

func (s *testSuite) TestBindSOAP() {
	type T struct {
		X int `xml:"x"`
	}

	req := httptest.NewRequest("GET", "/", bytes.NewBufferString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<Envelope xmlns=\"http://www.w3.org/2003/05/soap-envelope\"><Body xmlns=\"http://www.w3.org/2003/05/soap-envelope\"><T><x>7</x></T></Body></Envelope>"))
	req.Header.Set("SOAPAction", "Action")
	req.Header.Add(HeaderAccept, ContentTypeXML)
	w := httptest.NewRecorder()
	ctx := New(req, w, l)
	t := T{}
	ctx.Bind(&t)
	s.Equal(7, t.X, "Bind did not parse SOAP correctly")
}

func (s *testSuite) TestMarshalJSON() {
	type T struct {
		X int `json:"x"`
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	ctx := New(req, w, l)
	t := T{X: 7}
	ctx.OK(t)
	s.Equal("{\"x\":7}\n", w.Body.String(), "Bind did not marshal JSON correctly")
}

func (s *testSuite) TestMarshalXML() {
	type T struct {
		X int `xml:"x"`
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", ContentTypeXML)
	w := httptest.NewRecorder()
	req.Header.Add(HeaderAccept, ContentTypeXML)
	ctx := New(req, w, l)
	t := T{X: 7}
	ctx.OK(t)
	s.Equal("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<T><x>7</x></T>", w.Body.String(), "Bind did not marshal XML correctly")
}

func (s *testSuite) TestMarshalSOAP() {
	type T struct {
		X int `xml:"x"`
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", ContentTypeXML)
	req.Header.Set("SOAPAction", "Action")
	w := httptest.NewRecorder()
	req.Header.Add(HeaderAccept, ContentTypeXML)
	ctx := New(req, w, l)
	t := T{X: 7}
	ctx.OK(t)
	s.Equal("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<Envelope xmlns=\"http://www.w3.org/2003/05/soap-envelope\"><Body xmlns=\"http://www.w3.org/2003/05/soap-envelope\"><T><x>7</x></T></Body></Envelope>",
		w.Body.String(), "Bind did not marshal SOAP correctly")
}

func (s *testSuite) TestOK() {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	ctx := New(req, w, l)

	ctx.OK(7)
	s.Equal(200, w.Code)
}

func (s *testSuite) TestNotFound() {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	ctx := New(req, w, l)

	ctx.NotFound(7)
	s.Equal(404, w.Code)
}

func (s *testSuite) TestUnauthorized() {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	ctx := New(req, w, l)

	ctx.Unauthorized(7)
	s.Equal(401, w.Code)
}

func (s *testSuite) TestForbidden() {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	ctx := New(req, w, l)

	ctx.Forbidden(7)
	s.Equal(403, w.Code)
}

func (s *testSuite) TestBadRequest() {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	ctx := New(req, w, l)

	ctx.BadRequest(7)
	s.Equal(400, w.Code)
}
