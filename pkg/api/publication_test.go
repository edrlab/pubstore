package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/stretchr/testify/assert"
)

var api *Api

func TestGetPublicationHandler(t *testing.T) {
	// Initialize router
	r := chi.NewRouter()
	r.Get("/api/v1/publication/{id}", api.getPublicationHandler)
	r.Post("/api/v1/publication", api.createPublicationHandler)

	jsonData := `
	{
	    "title": "Test Publication",
	    "datePublication": "2023-06-16T12:00:00Z",
	    "description": "Test description",
	    "coverUrl": "http://example.com/cover.jpg",
	    "language": [
	        {
	            "code": "en"
	        },
	        {
	            "code": "fr"
	        }
	    ],
	    "publisher": [
	        {
	            "name": "Test Publisher A"
	        },
	        {
	            "name": "Test Publisher B"
	        }
	    ],
	    "author": [
	        {
	            "name": "Test Author A"
	        },
	        {
	            "name": "Test Author B"
	        }
	    ],
	    "category": [
	        {
	            "name": "Test Category A"
	        },
	        {
	            "name": "Test Category B"
	        }
	    ]
	}
	`

	var publication Publication
	err := json.Unmarshal([]byte(jsonData), &publication)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Encode the request body
	reqBody, err := json.Marshal(publication)
	assert.NoError(t, err)

	// Perform a POST request
	req := httptest.NewRequest("POST", "/api/v1/publication", bytes.NewBuffer(reqBody))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	// Check the response status code
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Decode the response body
	var result stor.Publication
	err = json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)

	// Check the response body
	pubUUID := result.UUID

	// Perform a GET request
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/publication/%s", pubUUID), nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rec.Code)

	// Decode the response body
	err = json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)

	// Check the response body
	assert.Equal(t, publication.Title, result.Title)
	assert.Equal(t, pubUUID, result.UUID)
	assert.Equal(t, publication.Description, result.Description)
}

func TestCreatePublicationHandler(t *testing.T) {
	// Initialize router
	r := chi.NewRouter()
	r.Post("/api/v1/publication", api.createPublicationHandler)

	jsonData := `
	{
	    "title": "Test Publication",
	    "datePublication": "2023-06-16T12:00:00Z",
	    "description": "Test description",
	    "coverUrl": "http://example.com/cover.jpg",
	    "language": [
	        {
	            "code": "en"
	        },
	        {
	            "code": "fr"
	        }
	    ],
	    "publisher": [
	        {
	            "name": "Test Publisher A"
	        },
	        {
	            "name": "Test Publisher B"
	        }
	    ],
	    "author": [
	        {
	            "name": "Test Author A"
	        },
	        {
	            "name": "Test Author B"
	        }
	    ],
	    "category": [
	        {
	            "name": "Test Category A"
	        },
	        {
	            "name": "Test Category B"
	        }
	    ]
	}
	`

	var publication Publication
	err := json.Unmarshal([]byte(jsonData), &publication)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Encode the request body
	reqBody, err := json.Marshal(publication)
	assert.NoError(t, err)

	// Perform a POST request
	req := httptest.NewRequest("POST", "/api/v1/publication", bytes.NewBuffer(reqBody))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	// Check the response status code
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Decode the response body
	var result stor.Publication
	err = json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)

	// Check the response body
	assert.Equal(t, publication.Title, result.Title)
	// assert.Equal(t, publication.UUID, result.UUID)
	assert.Equal(t, publication.Description, result.Description)
}

func TestMain(m *testing.M) {

	validate = validator.New()

	s := stor.Init("file::memory:?cache=shared")

	api = &Api{stor: s}

	// Run the tests
	exitCode := m.Run()

	s.Stop()

	fmt.Println("ExitCode", exitCode)
	// Exit with the appropriate exit code
	os.Exit(exitCode)
}

func TestSuite(t *testing.T) {
	t.Run("TestCreatePublicationHandler", TestCreatePublicationHandler)
	t.Run("TestGetPublicationHandler", TestGetPublicationHandler)
}
