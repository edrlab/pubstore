package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

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
	var publication stor.Publication
	err := json.NewDecoder(r.Body).Decode(&publication)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fmt.Println("######")
	fmt.Println(publication)
	fmt.Println("######")

	// Validate the publication struct using the validator
	err = validate.Struct(publication)
	if err != nil {
		// If validation fails, return the validation errors
		validationErrors := err.(validator.ValidationErrors)
		http.Error(w, validationErrors.Error(), http.StatusBadRequest)
		return
	}

	// Generate UUID for the publication
	publication.UUID = uuid.New().String()

	err = api.stor.CreatePublication(&publication)
	if err != nil {
		http.Error(w, "Failed to create publication", http.StatusInternalServerError)
		return
	}
	// Return success response
	w.WriteHeader(http.StatusCreated)
	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("######")
	fmt.Println(publication)
	fmt.Println("######")

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
	publication, err := api.stor.GetPublicationByUUID(publicationID)
	if err != nil {
		http.Error(w, "Publication not found", http.StatusNotFound)
		return
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the publication as JSON and write it to the response
	err = json.NewEncoder(w).Encode(publication)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
