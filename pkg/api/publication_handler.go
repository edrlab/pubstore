// Copyright 2023 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package api

import (
	"errors"
	"net/http"

	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/go-chi/render"
)

// @Summary Create a new publication
// @Description Create a new publication with the provided payload
// @Tags publications
// @Accept json
// @Produce json
// @Param publication body Publication true "Publication object"
// @Success 201 {object} Publication "Publication created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request payload or validation errors"
// @Failure 500 {object} ErrorResponse "Failed to create publication"
// @Router /publication [post]

// createPublication adds a new Publication to the database.
func (a *Api) createPublication(w http.ResponseWriter, r *http.Request) {

	// get the payload
	data := &PublicationRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	publication := data.Publication

	// db create
	err := a.Store.CreatePublication(publication)
	if err != nil {
		render.Render(w, r, ErrServer(err))
		return
	}

	render.Status(r, http.StatusCreated)
	if err := render.Render(w, r, NewPublicationResponse(publication)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// @Summary Get a publication by ID
// @Description Retrieve a publication by its ID
// @Tags publications
// @Accept json
// @Produce json
// @Param id path string true "Publication ID"
// @Success 200 {object} Publication "OK"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /publication/{id} [get]

// getPublication returns a specific publication
func (a *Api) getPublication(w http.ResponseWriter, r *http.Request) {

	publication := fromPubContext(r.Context())

	if err := render.Render(w, r, NewPublicationResponse(publication)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// @Summary Update a publication by ID
// @Description Update a publication with the provided payload
// @Tags publications
// @Accept json
// @Produce json
// @Param id path string true "Publication ID"
// @Param publication body Publication true "Publication object"
// @Success 200 {object} Publication "Publication updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request payload or validation errors"
// @Failure 500 {object} ErrorResponse "Failed to update publication"
// @Router /publication/{id} [put]

// updatePublication updates an existing Publication in the database.
func (a *Api) updatePublication(w http.ResponseWriter, r *http.Request) {

	// get the payload
	data := &PublicationRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	publication := data.Publication

	// get the existing publication
	currentPub := fromPubContext(r.Context())

	// force the ID field
	publication.ID = currentPub.ID

	// update
	err := a.Store.UpdatePublication(publication)
	if err != nil {
		render.Render(w, r, ErrServer(err))
		return
	}

	if err := render.Render(w, r, NewPublicationResponse(publication)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// listPublications lists all publications present in the database.
func (a *Api) listPublications(w http.ResponseWriter, r *http.Request) {

	pg := fromPaginateContext(r.Context())

	publications, err := a.Store.ListPublications(pg.Page, pg.PageSize)
	if err != nil {
		render.Render(w, r, ErrServer(err))
		return
	}
	if err := render.RenderList(w, r, NewPublicationListResponse(publications)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// searchPublications searches publications corresponding to a specific criteria.
func (a *Api) searchPublications(w http.ResponseWriter, r *http.Request) {
	var publications []stor.Publication
	var err error

	pg := fromPaginateContext(r.Context())

	// by format
	if format := r.URL.Query().Get("format"); format != "" {
		var contentType string
		switch format {
		case "epub":
			contentType = "application/epub+zip"
		case "pdf":
			contentType = "application/pdf+lcp"
		case "audiobook":
			contentType = "application/audiobook+lcp"
		case "divina":
			contentType = "application/divina+lcp"
		default:
			err = errors.New("invalid content type query string parameter")
		}
		if contentType != "" {
			publications, err = a.Store.FindPublicationsByType(contentType, pg.Page, pg.PageSize)
		}
	} else {
		render.Render(w, r, ErrNotFound)
		return
	}
	if err != nil {
		render.Render(w, r, ErrNotFound)
		return
	}
	if err := render.RenderList(w, r, NewPublicationListResponse(publications)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// @Summary Delete a publication by ID
// @Description Delete a publication by its ID
// @Tags publications
// @Accept json
// @Produce json
// @Param id path string true "Publication ID"
// @Success 200 "Publication deleted successfully"
// @Failure 500 {object} ErrorResponse "Failed to delete publication"
// @Router /publication/{id} [delete]

// deletePublication removes an existing Publication from the database.
func (a *Api) deletePublication(w http.ResponseWriter, r *http.Request) {

	// get the existing publication
	publication := fromPubContext(r.Context())

	// db delete
	err := a.Store.DeletePublication(publication)
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

//type omit *struct{}

// PublicationRequest is the request publication payload.
type PublicationRequest struct {
	*stor.Publication
}

// PublicationResponse is the response publication payload.
type PublicationResponse struct {
	*stor.Publication
	// do not serialize the following properties
	ID        omit `json:"ID,omitempty"`
	CreatedAt omit `json:"CreatedAt,omitempty"`
	UpdatedAt omit `json:"UpdatedAt,omitempty"`
	DeletedAt omit `json:"DeletedAt,omitempty"`
}

// NewPublicationListResponse creates a rendered list of publications
func NewPublicationListResponse(publications []stor.Publication) []render.Renderer {
	list := []render.Renderer{}
	for i := 0; i < len(publications); i++ {
		list = append(list, NewPublicationResponse(&publications[i]))
	}
	return list
}

// NewPublicationResponse creates a rendered publication.
func NewPublicationResponse(pub *stor.Publication) *PublicationResponse {
	return &PublicationResponse{Publication: pub}
}

// Bind post-processes requests after unmarshalling.
func (p *PublicationRequest) Bind(r *http.Request) error {
	return p.Publication.Validate()
}

// Render processes responses before marshalling.
func (pub *PublicationResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
