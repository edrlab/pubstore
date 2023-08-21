package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Language struct {
	Code string `json:"code"`
}

type Publisher struct {
	Name string `json:"name"`
}

type Author struct {
	Name string `json:"name"`
}

type Category struct {
	Name string `json:"name"`
}

type Publication struct {
	UUID            string      `json:"uuid" validate:"omitempty,uuid4_rfc4122"`
	Title           string      `json:"title"`
	DatePublication time.Time   `json:"datePublication"`
	Description     string      `json:"description"`
	CoverUrl        string      `json:"coverUrl"`
	Language        []Language  `json:"language"`
	Publisher       []Publisher `json:"publisher"`
	Author          []Author    `json:"author"`
	Category        []Category  `json:"category"`
}

func convertPublicationFromStor(originalPublication stor.Publication) Publication {
	convertedPublication := Publication{
		Title:           originalPublication.Title,
		UUID:            originalPublication.UUID,
		DatePublication: originalPublication.DatePublication,
		Description:     originalPublication.Description,
		CoverUrl:        originalPublication.CoverUrl,
	}

	// Convert Language slice
	for _, language := range originalPublication.Language {
		convertedPublication.Language = append(convertedPublication.Language, Language{Code: language.Code})
	}

	// Convert Publisher slice
	for _, publisher := range originalPublication.Publisher {
		convertedPublication.Publisher = append(convertedPublication.Publisher, Publisher{Name: publisher.Name})
	}

	// Convert Author slice
	for _, author := range originalPublication.Author {
		convertedPublication.Author = append(convertedPublication.Author, Author{Name: author.Name})
	}

	// Convert Category slice
	for _, category := range originalPublication.Category {
		convertedPublication.Category = append(convertedPublication.Category, Category{Name: category.Name})
	}

	return convertedPublication
}

func convertPublicationToStor(convertedPublication Publication, originalPublication *stor.Publication) *stor.Publication {
	if originalPublication == nil {
		originalPublication = &stor.Publication{}
	}

	if convertedPublication.Title != "" {
		originalPublication.Title = convertedPublication.Title
	}
	if convertedPublication.UUID != "" {
		originalPublication.UUID = convertedPublication.UUID
	}
	if convertedPublication.DatePublication.IsZero() {
		originalPublication.DatePublication = convertedPublication.DatePublication
	}
	if convertedPublication.Description != "" {
		originalPublication.Description = convertedPublication.Description
	}
	if convertedPublication.CoverUrl != "" {
		originalPublication.CoverUrl = convertedPublication.CoverUrl
	}

	// Convert Language slice
	for _, language := range convertedPublication.Language {
		if language.Code != "" {
			originalPublication.Language = append(originalPublication.Language, stor.Language{Code: language.Code})
		}
	}

	// Convert Publisher slice
	for _, publisher := range convertedPublication.Publisher {
		if publisher.Name != "" {
			originalPublication.Publisher = append(originalPublication.Publisher, stor.Publisher{Name: publisher.Name})
		}
	}

	// Convert Author slice
	for _, author := range convertedPublication.Author {
		if author.Name != "" {
			originalPublication.Author = append(originalPublication.Author, stor.Author{Name: author.Name})
		}
	}

	// Convert Category slice
	for _, category := range convertedPublication.Category {
		if category.Name != "" {
			originalPublication.Category = append(originalPublication.Category, stor.Category{Name: category.Name})
		}
	}

	return originalPublication
}

/*
*

	curl -X POST \
	  -H "Content-Type: application/json" \
	  -d '{
	    "title": "Test Publication",
	    "datePublication": "2023-06-16T12:00:00Z",
	    "description": "Test description",
	    "coverUrl": "http://example.com/cover.jpg",
	    "language": [
	      {"code": "en"},
	      {"code": "fr"}
	    ],
	    "publisher": [
	      {"name": "Test Publisher A"},
	      {"name": "Test Publisher B"}
	    ],
	    "author": [
	      {"name": "Test Author A"},
	      {"name": "Test Author B"}
	    ],
	    "category": [
	      {"name": "Test Category A"},
	      {"name": "Test Category B"}
	    ]
	  }' \
	  http://localhost:8080/api/v1/publications
*/

type ErrorResponse string

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
func (api *Api) createPublicationHandler(w http.ResponseWriter, r *http.Request) {
	// Parse and validate the request body
	var publication Publication
	err := json.NewDecoder(r.Body).Decode(&publication)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the publication struct using the validator
	err = validate.Struct(publication)
	if err != nil {
		// If validation fails, return the validation errors
		validationErrors := err.(validator.ValidationErrors)
		http.Error(w, validationErrors.Error(), http.StatusBadRequest)
		return
	}

	// Generate UUID for the publication
	if publication.UUID == "" {
		publication.UUID = uuid.New().String()
	}

	publicationStor := convertPublicationToStor(publication, &stor.Publication{})

	err = api.stor.CreatePublication(publicationStor)
	if err != nil {
		http.Error(w, "Failed to create publication", http.StatusInternalServerError)
		return
	}
	// Return success response
	w.WriteHeader(http.StatusCreated)
	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the publication as JSON and write it to the response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(publication)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	// fmt.Fprint(w, "Publication created successfully")
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
func (api *Api) getPublicationHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	storPublication, ok := ctx.Value("publication").(*stor.Publication)
	if !ok {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	publication := convertPublicationFromStor(*storPublication)

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the publication as JSON and write it to the response
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(publication)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
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
func (api *Api) updatePublicationHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	storPublication, ok := ctx.Value("publication").(*stor.Publication)
	if !ok {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Parse and validate the request body
	var publication Publication
	err := json.NewDecoder(r.Body).Decode(&publication)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the publication struct using the validator
	err = validate.Struct(publication)
	if err != nil {
		// If validation fails, return the validation errors
		validationErrors := err.(validator.ValidationErrors)
		http.Error(w, validationErrors.Error(), http.StatusBadRequest)
		return
	}

	storPublicationConverted := convertPublicationToStor(publication, storPublication)

	err = api.stor.UpdatePublication(storPublicationConverted)
	if err != nil {
		http.Error(w, "Failed to update publication", http.StatusInternalServerError)
		return
	}

	publication = convertPublicationFromStor(*storPublicationConverted)

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the publication as JSON and write it to the response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(publication)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
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
func (api *Api) deletePublicationHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	storPublication, ok := ctx.Value("publication").(*stor.Publication)
	if !ok {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err := api.stor.DeletePublication(storPublication)
	if err != nil {
		http.Error(w, "Failed to delete publication", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
