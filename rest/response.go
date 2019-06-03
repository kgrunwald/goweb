package rest

import "net/http"

// The Response struct encodes an HTTP response. Controllers must return a Response object so that the REST package can
// encode it according to the Accept header of the request.
type Response struct {
	StatusCode int
	Body       interface{}
}

// NewResponse constructs a new Reponse with the provided status code and body
func NewResponse(status int, body interface{}) *Response {
	return &Response{
		Body:       body,
		StatusCode: status,
	}
}

// OK is a helper method that returns a response with a 200 status code
func OK(body interface{}) *Response {
	return NewResponse(http.StatusOK, body)
}

// NotFound is a helper method that returns a response with a 404 status code
func NotFound(body interface{}) *Response {
	return NewResponse(http.StatusNotFound, body)
}

// Forbidden is a helper method that returns a response with a 403 status code
func Forbidden(body interface{}) *Response {
	return NewResponse(http.StatusForbidden, body)
}

// Unauthorized is a helper method that returns a response with a 401 status code
func Unauthorized(body interface{}) *Response {
	return NewResponse(http.StatusUnauthorized, body)
}

// BadRequest is a helper method that returns a response with a 400 status code
func BadRequest(body interface{}) *Response {
	return NewResponse(http.StatusBadRequest, body)
}
