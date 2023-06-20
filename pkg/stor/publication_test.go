package stor

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

var stor Stor

func checkLanguageEquality(a, b []Language) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i].Code != b[i].Code {
			return false
		}
	}

	return true
}

func checkPublisherEquality(a, b []Publisher) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i].Name != b[i].Name {
			return false
		}
	}

	return true
}

func checkAuthorEquality(a, b []Author) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i].Name != b[i].Name {
			return false
		}
	}

	return true
}

func checkCategoryEquality(a, b []Category) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i].Name != b[i].Name {
			return false
		}
	}

	return true
}

func TestPublicationCRUD(t *testing.T) {
	// Test CreatePublication

	pubUUID := uuid.New().String()
	publication := &Publication{
		Title:           "Test Publication",
		UUID:            pubUUID,
		DatePublication: time.Now(),
		Description:     "Test description",
		CoverUrl:        "http://example.com/cover.jpg",
		Language: []Language{
			{Code: "en"},
			{Code: "fr"},
		},
		Publisher: []Publisher{
			{Name: "Test Publisher A"},
			{Name: "Test Publisher B"},
		},
		Author: []Author{
			{Name: "Test Author A"},
			{Name: "Test Author B"},
		},
		Category: []Category{
			{Name: "Test Category A"},
			{Name: "Test Category B"},
		},
	}

	err := stor.CreatePublication(publication)
	if err != nil {
		t.Errorf("Error creating publication: %s", err.Error())
	}

	// Test GetPublicationByID
	fetchedPublication, err := stor.GetPublicationByUUID(pubUUID)
	if err != nil {
		t.Errorf("Error getting publication by ID: %s", err.Error())
	} else {
		// Ensure fetched publication matches the created publication
		if fetchedPublication.ID != publication.ID ||
			fetchedPublication.UUID != publication.UUID ||
			!fetchedPublication.DatePublication.Equal(publication.DatePublication) ||
			fetchedPublication.Description != publication.Description ||
			!checkLanguageEquality(fetchedPublication.Language, publication.Language) ||
			!checkPublisherEquality(fetchedPublication.Publisher, publication.Publisher) ||
			!checkAuthorEquality(fetchedPublication.Author, publication.Author) ||
			!checkCategoryEquality(fetchedPublication.Category, publication.Category) {
			t.Error("Fetched publication does not match the created publication")
		}

	}

	// Test GetPublicationByTitle
	fetchedPublication, err = stor.GetPublicationByTitle(publication.Title)
	if err != nil {
		t.Errorf("Error getting publication by title: %s", err.Error())
	} else {
		// Ensure fetched publication matches the created publication
		if fetchedPublication.ID != publication.ID ||
			fetchedPublication.UUID != publication.UUID ||
			!fetchedPublication.DatePublication.Equal(publication.DatePublication) ||
			fetchedPublication.Description != publication.Description ||
			!checkLanguageEquality(fetchedPublication.Language, publication.Language) ||
			!checkPublisherEquality(fetchedPublication.Publisher, publication.Publisher) ||
			!checkAuthorEquality(fetchedPublication.Author, publication.Author) ||
			!checkCategoryEquality(fetchedPublication.Category, publication.Category) {
			t.Error("Fetched publication does not match the created publication")
		}
	}

	// Test UpdatePublication
	fetchedPublication.Title = "Updated Test Publication"
	err = stor.UpdatePublication(fetchedPublication)
	if err != nil {
		t.Errorf("Error updating publication: %s", err.Error())
	}

	// Fetch the updated publication again to ensure the changes were saved
	updatedPublication, err := stor.GetPublicationByUUID(pubUUID)
	if err != nil {
		t.Errorf("Error getting publication by ID: %s", err.Error())
	} else {
		if updatedPublication.Title != fetchedPublication.Title {
			t.Error("Publication title was not updated")
		}
	}

	// Test DeletePublication
	err = stor.DeletePublication(updatedPublication)
	if err != nil {
		t.Errorf("Error deleting publication: %s", err.Error())
	}

	// Ensure the publication is no longer present in the database
	_, err = stor.GetPublicationByUUID(pubUUID)
	if err == nil {
		t.Error("Publication was not deleted")
	}
}

func TestGetPublicationByCategory(t *testing.T) {
	// Create test publications
	publication1 := &Publication{
		Title:           "Test Publication 1",
		UUID:            uuid.New().String(),
		DatePublication: time.Now(),
		Description:     "Test description",
		Category: []Category{
			{Name: "Category A"},
		},
	}
	publication2 := &Publication{
		Title:           "Test Publication 2",
		UUID:            uuid.New().String(),
		DatePublication: time.Now(),
		Description:     "Test description",
		Category: []Category{
			{Name: "Category B"},
		},
	}

	err := stor.CreatePublication(publication1)
	if err != nil {
		t.Errorf("Error creating publication 1: %s", err.Error())
	}
	err = stor.CreatePublication(publication2)
	if err != nil {
		t.Errorf("Error creating publication 2: %s", err.Error())
	}

	// Test GetPublicationByCategory
	publications, err := stor.GetPublicationsByCategory("Category B")
	if err != nil {
		t.Errorf("Error getting publications by category: %s", err.Error())
	} else {
		// Ensure the correct number of publications is retrieved
		if len(publications) != 1 {
			t.Errorf("Expected 1 publication, got %d", len(publications))
		}

		// Ensure the retrieved publication matches the created publication
		if publications[0].Title != publication2.Title ||
			publications[0].UUID != publication2.UUID ||
			!publications[0].DatePublication.Equal(publication2.DatePublication) ||
			publications[0].Description != publication2.Description ||
			!checkLanguageEquality(publications[0].Language, publication2.Language) ||
			!checkPublisherEquality(publications[0].Publisher, publication2.Publisher) ||
			!checkAuthorEquality(publications[0].Author, publication2.Author) ||
			!checkCategoryEquality(publications[0].Category, publication2.Category) {
			t.Error("Fetched publication does not match the created publication")
		}
	}

	// Clean up the test data
	err = stor.DeletePublication(publication1)
	if err != nil {
		t.Errorf("Error deleting publication 1: %s", err.Error())
	}
	err = stor.DeletePublication(publication2)
	if err != nil {
		t.Errorf("Error deleting publication 2: %s", err.Error())
	}
}

func TestGetPublicationByAuthor(t *testing.T) {
	// Create test publications
	publication1 := &Publication{
		Title:           "Test Publication 1",
		UUID:            uuid.New().String(),
		DatePublication: time.Now(),
		Description:     "Test description",
		Author: []Author{
			{Name: "Author A"},
		},
	}
	publication2 := &Publication{
		Title:           "Test Publication 2",
		UUID:            uuid.New().String(),
		DatePublication: time.Now(),
		Description:     "Test description",
		Author: []Author{
			{Name: "Author B"},
		},
	}

	err := stor.CreatePublication(publication1)
	if err != nil {
		t.Errorf("Error creating publication 1: %s", err.Error())
	}
	err = stor.CreatePublication(publication2)
	if err != nil {
		t.Errorf("Error creating publication 2: %s", err.Error())
	}

	// Test GetPublicationByAuthor
	publications, err := stor.GetPublicationsByAuthor("Author B")
	if err != nil {
		t.Errorf("Error getting publications by author: %s", err.Error())
	} else {
		// Ensure the correct number of publications is retrieved
		if len(publications) != 1 {
			t.Errorf("Expected 1 publication, got %d", len(publications))
		}

		// Ensure the retrieved publication matches the created publication
		if publications[0].Title != publication2.Title ||
			publications[0].UUID != publication2.UUID ||
			!publications[0].DatePublication.Equal(publication2.DatePublication) ||
			publications[0].Description != publication2.Description ||
			!checkLanguageEquality(publications[0].Language, publication2.Language) ||
			!checkPublisherEquality(publications[0].Publisher, publication2.Publisher) ||
			!checkAuthorEquality(publications[0].Author, publication2.Author) ||
			!checkCategoryEquality(publications[0].Category, publication2.Category) {
			t.Error("Fetched publication does not match the created publication")
		}
	}

	// Clean up the test data
	err = stor.DeletePublication(publication1)
	if err != nil {
		t.Errorf("Error deleting publication 1: %s", err.Error())
	}
	err = stor.DeletePublication(publication2)
	if err != nil {
		t.Errorf("Error deleting publication 2: %s", err.Error())
	}
}

func TestGetPublicationByPublisher(t *testing.T) {
	// Create test publications
	publication1 := &Publication{
		Title:           "Test Publication 1",
		UUID:            uuid.New().String(),
		DatePublication: time.Now(),
		Description:     "Test description",
		Publisher: []Publisher{
			{Name: "Publisher A"},
		},
	}
	publication2 := &Publication{
		Title:           "Test Publication 2",
		UUID:            uuid.New().String(),
		DatePublication: time.Now(),
		Description:     "Test description",
		Publisher: []Publisher{
			{Name: "Publisher B"},
		},
	}

	err := stor.CreatePublication(publication1)
	if err != nil {
		t.Errorf("Error creating publication 1: %s", err.Error())
	}
	err = stor.CreatePublication(publication2)
	if err != nil {
		t.Errorf("Error creating publication 2: %s", err.Error())
	}

	// Test GetPublicationByPublisher
	publications, err := stor.GetPublicationsByPublisher("Publisher B")
	if err != nil {
		t.Errorf("Error getting publications by publisher: %s", err.Error())
	} else {
		// Ensure the correct number of publications is retrieved
		if len(publications) != 1 {
			t.Errorf("Expected 1 publication, got %d", len(publications))
		}

		// Ensure the retrieved publication matches the created publication
		if publications[0].Title != publication2.Title ||
			publications[0].UUID != publication2.UUID ||
			!publications[0].DatePublication.Equal(publication2.DatePublication) ||
			publications[0].Description != publication2.Description ||
			!checkLanguageEquality(publications[0].Language, publication2.Language) ||
			!checkPublisherEquality(publications[0].Publisher, publication2.Publisher) ||
			!checkAuthorEquality(publications[0].Author, publication2.Author) ||
			!checkCategoryEquality(publications[0].Category, publication2.Category) {
			t.Error("Fetched publication does not match the created publication")
		}
	}

	// Clean up the test data
	err = stor.DeletePublication(publication1)
	if err != nil {
		t.Errorf("Error deleting publication 1: %s", err.Error())
	}
	err = stor.DeletePublication(publication2)
	if err != nil {
		t.Errorf("Error deleting publication 2: %s", err.Error())
	}
}

func TestGetPublicationByLanguage(t *testing.T) {
	// Create test publications
	publication1 := &Publication{
		Title:           "Test Publication 1",
		UUID:            uuid.New().String(),
		DatePublication: time.Now(),
		Description:     "Test description",
		Language: []Language{
			{Code: "aa"},
		},
	}
	publication2 := &Publication{
		Title:           "Test Publication 2",
		UUID:            uuid.New().String(),
		DatePublication: time.Now(),
		Description:     "Test description",
		Language: []Language{
			{Code: "bb"},
		},
	}

	err := stor.CreatePublication(publication1)
	if err != nil {
		t.Errorf("Error creating publication 1: %s", err.Error())
	}
	err = stor.CreatePublication(publication2)
	if err != nil {
		t.Errorf("Error creating publication 2: %s", err.Error())
	}

	// Test GetPublicationByLanguage
	publications, err := stor.GetPublicationsByLanguage("bb")
	if err != nil {
		t.Errorf("Error getting publications by language: %s", err.Error())
	} else {
		// Ensure the correct number of publications is retrieved
		if len(publications) != 1 {
			t.Errorf("Expected 1 publication, got %d", len(publications))
		}

		// Ensure the retrieved publication matches the created publication
		if publications[0].Title != publication2.Title ||
			publications[0].UUID != publication2.UUID ||
			!publications[0].DatePublication.Equal(publication2.DatePublication) ||
			publications[0].Description != publication2.Description ||
			!checkLanguageEquality(publications[0].Language, publication2.Language) ||
			!checkPublisherEquality(publications[0].Publisher, publication2.Publisher) ||
			!checkAuthorEquality(publications[0].Author, publication2.Author) ||
			!checkCategoryEquality(publications[0].Category, publication2.Category) {
			t.Error("Fetched publication does not match the created publication")
		}
	}

	// Clean up the test data
	err = stor.DeletePublication(publication1)
	if err != nil {
		t.Errorf("Error deleting publication 1: %s", err.Error())
	}
	err = stor.DeletePublication(publication2)
	if err != nil {
		t.Errorf("Error deleting publication 2: %s", err.Error())
	}
}

func TestCreate2PublicationsWithSameCategory(t *testing.T) {
	publication := &Publication{
		Title:           "Test Publication",
		UUID:            uuid.New().String(),
		DatePublication: time.Now(),
		Description:     "Test description",
		Language: []Language{
			{Code: "en"},
			{Code: "fr"},
		},
		Publisher: []Publisher{
			{Name: "Test Publisher A"},
			{Name: "Test Publisher B"},
		},
		Author: []Author{
			{Name: "Test Author A"},
			{Name: "Test Author B"},
		},
		Category: []Category{
			{Name: "Test Category A"},
			{Name: "Test Category B"},
		},
	}

	publication2 := &Publication{
		Title:           "Test Publication",
		UUID:            uuid.New().String(),
		DatePublication: time.Now(),
		Description:     "Test description",
		Language: []Language{
			{Code: "en"},
			{Code: "fr"},
		},
		Publisher: []Publisher{
			{Name: "Test Publisher A"},
			{Name: "Test Publisher B"},
		},
		Author: []Author{
			{Name: "Test Author A"},
			{Name: "Test Author B"},
		},
		Category: []Category{
			{Name: "Test Category A"},
			{Name: "Test Category B"},
		},
	}

	err := stor.CreatePublication(publication)
	if err != nil {
		t.Errorf("Error creating publication: %s", err.Error())
	}

	err = stor.CreatePublication(publication2)
	if err != nil {
		t.Errorf("Error creating publication: %s", err.Error())
	}

	categories, err2 := stor.GetCategories()
	if err2 != nil {
		t.Errorf("Error getting categories: %s", err.Error())
	}

	// Ensure the correct number of categories is retrieved
	// if len(categories) != 2 {
	// 	t.Errorf("Expected 2 categories, got %d", len(categories))
	// }

	category1 := Category{Name: "Test Category A"}
	category2 := Category{Name: "Test Category B"}

	// Ensure the retrieved categories match the created categories
	foundCategory1 := false
	foundCategory2 := false
	for _, category := range categories {
		if category.Name == category1.Name {
			foundCategory1 = true
		} else if category.Name == category2.Name {
			foundCategory2 = true
		}
	}
	if !foundCategory1 || !foundCategory2 {
		t.Error("Not all created categories were retrieved")
	}

	// Clean up the test data
	err = stor.DeletePublication(publication)
	if err != nil {
		t.Errorf("Error deleting publication 1: %s", err.Error())
	}
	err = stor.DeletePublication(publication2)
	if err != nil {
		t.Errorf("Error deleting publication 2: %s", err.Error())
	}

}

func TestGetPublicationByUUID(t *testing.T) {
	// Create test publications

	pubUUID := uuid.New().String()
	publication := &Publication{
		Title:           "Test Publication",
		UUID:            pubUUID,
		DatePublication: time.Now(),
		Description:     "Test description",
		Language: []Language{
			{Code: "en"},
			{Code: "fr"},
		},
		Publisher: []Publisher{
			{Name: "Test Publisher A"},
			{Name: "Test Publisher B"},
		},
		Author: []Author{
			{Name: "Test Author A"},
			{Name: "Test Author B"},
		},
		Category: []Category{
			{Name: "Test Category A"},
			{Name: "Test Category B"},
		},
	}

	err := stor.CreatePublication(publication)
	if err != nil {
		t.Errorf("Error creating publication 1: %s", err.Error())
	}

	fetchedPublication, _ := stor.GetPublicationByUUID(pubUUID)
	fmt.Println(fetchedPublication)
	// Ensure fetched publication matches the created publication
	if fetchedPublication.ID != publication.ID ||
		fetchedPublication.UUID != publication.UUID ||
		!fetchedPublication.DatePublication.Equal(publication.DatePublication) ||
		fetchedPublication.Description != publication.Description ||
		!checkLanguageEquality(fetchedPublication.Language, publication.Language) ||
		!checkPublisherEquality(fetchedPublication.Publisher, publication.Publisher) ||
		!checkAuthorEquality(fetchedPublication.Author, publication.Author) ||
		!checkCategoryEquality(fetchedPublication.Category, publication.Category) {
		t.Error("Fetched publication does not match the created publication")
	}

	// Clean up the test data
	err = stor.DeletePublication(publication)
	if err != nil {
		t.Errorf("Error deleting publication 1: %s", err.Error())
	}
}
