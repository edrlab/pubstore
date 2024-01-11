// Copyright 2023 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package api

import (
	"net/http"

	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/go-chi/render"
)

// @Summary Create a new user
// @Description Create a new user with the provided payload
// @Tags users
// @Accept json
// @Produce json
// @Param user body User true "User object"
// @Success 201 {object} User "User created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request payload or validation error"
// @Failure 500 {object} ErrorResponse "Failed to create user"
// @Router /users [post]
func (a *Api) createUser(w http.ResponseWriter, r *http.Request) {

	// get the payload
	data := &UserRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	user := data.User

	// force an empty sessionId
	user.SessionId = ""

	// db create
	err := a.Store.CreateUser(user)
	if err != nil {
		render.Render(w, r, ErrServer(err))
		return
	}

	render.Status(r, http.StatusCreated)
	if err := render.Render(w, r, NewUserResponse(user)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// @Summary Get a user by ID
// @Description Retrieve a user by its ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} User "OK"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /users/{id} [get]
func (a *Api) getUser(w http.ResponseWriter, r *http.Request) {

	user := fromUserContext(r.Context())

	if err := render.Render(w, r, NewUserResponse(user)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// @Summary Update a user by ID
// @Description Update a user with the provided payload
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body User true "User object"
// @Success 200 {object} User "User updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request payload or validation errors"
// @Failure 500 {object} ErrorResponse "Failed to update user"
// @Router /users/{id} [put]
func (a *Api) updateUser(w http.ResponseWriter, r *http.Request) {

	// get the payload
	data := &UserRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	user := data.User

	// get the existing user
	currentUser := fromUserContext(r.Context())

	// force the ID field
	user.ID = currentUser.ID

	// update
	if err := a.Store.UpdateUser(user); err != nil {
		render.Render(w, r, ErrServer(err))
		return
	}

	if err := render.Render(w, r, NewUserResponse(user)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// @Summary List users
// @Description List users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} stor.User
// @Failure 422 {object} ErrorResponse "Error rendering response"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /users [get]
func (a *Api) listUsers(w http.ResponseWriter, r *http.Request) {

	pg := fromPaginateContext(r.Context())

	users, err := a.Store.ListUsers(pg.Page, pg.PageSize)
	if err != nil {
		render.Render(w, r, ErrServer(err))
		return
	}
	if err := render.RenderList(w, r, NewUserListResponse(users)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// @Summary Delete a user by ID
// @Description Delete a user by its ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 "User deleted successfully"
// @Failure 500 {object} ErrorResponse "Failed to delete user"
// @Router /users/{id} [delete]
func (a *Api) deleteUser(w http.ResponseWriter, r *http.Request) {

	// get the existing user
	user := fromUserContext(r.Context())

	// delete
	err := a.Store.DeleteUser(user)
	if err != nil {
		render.Render(w, r, ErrServer(err))
		return
	}

	// return a simple ok status
	w.WriteHeader(http.StatusOK)
}

// --
// Request and Response payloads for the REST api.
// --

type omit *struct{}

// UserRequest is the request user payload.
type UserRequest struct {
	*stor.User
}

// UserResponse is the response user payload.
type UserResponse struct {
	*stor.User
	// do not serialize the following properties
	ID          omit `json:"ID,omitempty"`
	CreatedAt   omit `json:"CreatedAt,omitempty"`
	UpdatedAt   omit `json:"UpdatedAt,omitempty"`
	DeletedAt   omit `json:"DeletedAt,omitempty"`
	Password    omit `json:"password,omitempty"`
	Passphrase  omit `json:"passphrase,omitempty"`
	HPassword   omit `json:"hpassword,omitempty"`
	HPassphrase omit `json:"hpassphrase,omitempty"`
}

// NewUserListResponse creates a rendered list of users
func NewUserListResponse(users []stor.User) []render.Renderer {
	list := []render.Renderer{}
	for i := 0; i < len(users); i++ {
		list = append(list, NewUserResponse(&users[i]))
	}
	return list
}

// NewUserResponse creates a rendered user.
func NewUserResponse(user *stor.User) *UserResponse {
	return &UserResponse{User: user}
}

// Bind post-processes requests after unmarshalling.
func (p *UserRequest) Bind(r *http.Request) error {
	return p.User.Validate()
}

// Render processes responses before marshalling.
func (pub *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
