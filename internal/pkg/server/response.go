package server

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

// Response is basic response instance.
type Response struct {
	StatusCode int    `json:"-"`
	StatusText string `json:"-"`
	Status     bool   `json:"status"`
}

// ErrorResponse is basic error response instance.
type ErrorResponse struct {
	*Response
	Error     error  `json:"-"`
	ErrorText string `json:"error"`
}

// NewResponse returns new basic response instance.
func NewResponse(status int) *Response {
	return &Response{
		StatusCode: status,
		StatusText: http.StatusText(status),
		Status:     true,
	}
}

// NewErrorResponse returns new basic error response instance.
func NewErrorResponse(status int, err error) *ErrorResponse {
	r := NewResponse(status)
	r.Status = false

	return &ErrorResponse{
		Response:  r,
		Error:     err,
		ErrorText: err.Error(),
	}
}

// Render renders status code from response instance to response writer.
func (r Response) Render(_ http.ResponseWriter, req *http.Request) error {
	render.Status(req, r.StatusCode)

	return nil
}

// RenderResponse is supposed to be the only method to return any response to the client.
func RenderResponse(w http.ResponseWriter, r *http.Request, response render.Renderer) {
	if err := render.Render(w, r, response); err != nil {
		logrus.Error(err)

		render.Status(r, http.StatusInternalServerError)

		return
	}
}
