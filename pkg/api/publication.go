package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/go-chi/chi/v5"
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
	Title           string      `json:"title"`
	UUID            string      `json:"uuid" validate:"omitempty,uuid4_rfc4122"`
	DatePublication time.Time   `json:"datePublication"`
	Description     string      `json:"description"`
	CoverUrl        string      `json:"coverUrl"`
	Language        []Language  `json:"language"`
	Publisher       []Publisher `json:"publisher"`
	Author          []Author    `json:"author"`
	Category        []Category  `json:"category"`
}

func ConvertPublicationFromStor(originalPublication stor.Publication) Publication {
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

func ConvertPublicationToStor(convertedPublication Publication) stor.Publication {
	originalPublication := stor.Publication{
		Title:           convertedPublication.Title,
		UUID:            convertedPublication.UUID,
		DatePublication: convertedPublication.DatePublication,
		Description:     convertedPublication.Description,
		CoverUrl:        convertedPublication.CoverUrl,
	}

	// Convert Language slice
	for _, language := range convertedPublication.Language {
		originalPublication.Language = append(originalPublication.Language, stor.Language{Code: language.Code})
	}

	// Convert Publisher slice
	for _, publisher := range convertedPublication.Publisher {
		originalPublication.Publisher = append(originalPublication.Publisher, stor.Publisher{Name: publisher.Name})
	}

	// Convert Author slice
	for _, author := range convertedPublication.Author {
		originalPublication.Author = append(originalPublication.Author, stor.Author{Name: author.Name})
	}

	// Convert Category slice
	for _, category := range convertedPublication.Category {
		originalPublication.Category = append(originalPublication.Category, stor.Category{Name: category.Name})
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
	  http://localhost:8080/api/v1/publication
*/
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
	if len(publication.UUID) == 0 {
		publication.UUID = uuid.New().String()
	}

	publicationStor := ConvertPublicationToStor(publication)

	err = api.stor.CreatePublication(&publicationStor)
	if err != nil {
		http.Error(w, "Failed to create publication", http.StatusInternalServerError)
		return
	}
	// Return success response
	w.WriteHeader(http.StatusCreated)
	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the publication as JSON and write it to the response
	err = json.NewEncoder(w).Encode(publication)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	// fmt.Fprint(w, "Publication created successfully")
}

func (api *Api) getPublicationHandler(w http.ResponseWriter, r *http.Request) {
	// Get the publication ID from the URL parameters
	publicationID := chi.URLParam(r, "id")

	// Retrieve the publication from the database
	publicationStor, err := api.stor.GetPublicationByUUID(publicationID)
	if err != nil {
		http.Error(w, "Publication not found", http.StatusNotFound)
		return
	}

	fmt.Println(publicationStor)

	publication := ConvertPublicationFromStor(*publicationStor)

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the publication as JSON and write it to the response
	err = json.NewEncoder(w).Encode(publication)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
