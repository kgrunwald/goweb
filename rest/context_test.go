package rest

import (
	"bytes"
	"testing"

	"net/http/httptest"

	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
}

func TestContext(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (s *testSuite) TestNewContext() {
	req := httptest.NewRequest("GET", "/", nil)
	ctx := NewContext(req)
	s.EqualValues(ctx.Request(), req)
}

func (s *testSuite) TestBindJSON() {
	type T struct {
		X int `json:"x"`
	}

	req := httptest.NewRequest("GET", "/", bytes.NewBufferString(`{"x": 5}`))
	ctx := NewContext(req)
	t := T{}
	ctx.Bind(&t)
	s.Equal(5, t.X, "Bind did not parse JSON correctly")
}

func (s *testSuite) TestBindXML() {
	type T struct {
		X int `xml:"x"`
	}

	req := httptest.NewRequest("GET", "/", bytes.NewBufferString(`<T><x>5</x></T>`))
	req.Header.Add(HeaderAccept, ContentTypeXML)
	ctx := NewContext(req)
	t := T{}
	ctx.Bind(&t)
	s.Equal(5, t.X, "Bind did not parse XML correctly")
}

func (s *testSuite) TestMarshalJSON() {
	type T struct {
		X int `json:"x"`
	}

	req := httptest.NewRequest("GET", "/", nil)
	ctx := NewContext(req)
	t := T{X: 7}
	str, _ := ctx.Marshal(OK(t))
	s.Equal(`{"x":7}`, string(str), "Bind did not marshal JSON correctly")
}

func (s *testSuite) TestMarshalXML() {
	type T struct {
		X int `xml:"x"`
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add(HeaderAccept, ContentTypeXML)
	ctx := NewContext(req)
	t := T{X: 7}
	str, _ := ctx.Marshal(OK(t))
	s.Equal(`<T><x>7</x></T>`, string(str), "Bind did not marshal XML correctly")
}
