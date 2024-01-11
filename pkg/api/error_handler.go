// Copyright 2023 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package api

import (
	"net/http"

	"github.com/go-chi/render"
)

// See Problem Details for HTTP APIs, rfc 7807 : https://tools.ietf.org/html/rfc7807
// Error response payloads & renderers

// Error types
const ERROR_BASE_URL = "http://edrlab.org/pubstore/error/"
const SERVER_ERROR = ERROR_BASE_URL + "server"

type ErrResponse struct {
	//not serialized
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code
	//mandatory
	Type  string `json:"type"`
	Title string `json:"title"`
	//optional
	Detail   string `json:"detail,omitempty"` // application-level error message
	Instance string `json:"instance,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		Type:           "about:blank",
		Title:          "Invalid request",
		Detail:         err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		Type:           "about:blank",
		Title:          "Error rendering response",
		Detail:         err.Error(),
	}
}

func ErrServer(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		Type:           "about:blank",
		Title:          "An unexpected server error has occurred.",
		Detail:         err.Error(),
	}
}

var ErrNotFound = &ErrResponse{
	HTTPStatusCode: 404,
	Type:           "about:blank",
	Title:          "Resource not found.",
}
