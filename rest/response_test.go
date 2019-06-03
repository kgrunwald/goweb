package rest

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResponse(t *testing.T) {
	r := NewResponse(1, "s")
	assert.Equal(t, r.Body, "s", "Response body not correct")
	assert.Equal(t, r.StatusCode, 1, "Response status code not correct")
}

func TestOK(t *testing.T) {
	r := OK("s")
	assert.Equal(t, r.Body, "s", "Response body not correct")
	assert.Equal(t, r.StatusCode, http.StatusOK, "Response status code not correct")
}

func TestNotFound(t *testing.T) {
	r := NotFound("s")
	assert.Equal(t, r.Body, "s", "Response body not correct")
	assert.Equal(t, r.StatusCode, http.StatusNotFound, "Response status code not correct")
}

func TestUnauthorized(t *testing.T) {
	r := Unauthorized("s")
	assert.Equal(t, r.Body, "s", "Response body not correct")
	assert.Equal(t, r.StatusCode, http.StatusUnauthorized, "Response status code not correct")
}

func TestForbidden(t *testing.T) {
	r := Forbidden("s")
	assert.Equal(t, r.Body, "s", "Response body not correct")
	assert.Equal(t, r.StatusCode, http.StatusForbidden, "Response status code not correct")
}

func TestBadRequest(t *testing.T) {
	r := BadRequest("s")
	assert.Equal(t, r.Body, "s", "Response body not correct")
	assert.Equal(t, r.StatusCode, http.StatusBadRequest, "Response status code not correct")
}
