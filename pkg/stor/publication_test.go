package stor

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func checkPublicationEquality(a, b Publication) bool {
	if a.ID != b.ID ||
		a.UUID != b.UUID ||
		a.Title != b.Title ||
		a.ContentType != b.ContentType ||
		a.DatePublished != b.DatePublished ||
		a.Description != b.Description ||
		!checkLanguageEquality(a.Language, b.Language) ||
		!checkPublisherEquality(a.Publisher, b.Publisher) ||
		!checkAuthorEquality(a.Author, b.Author) ||
		!checkCategoryEquality(a.Category, b.Category) {
		return false
	}

	return true
}

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

	// start by checking categories where there are none
	categories, err := store.GetCategories()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(categories))

	// create a new publication
	publication := &Publication{
		UUID:          uuid.New().String(),
		Title:         "Test Publication",
		ContentType:   "application/zip",
		DatePublished: "2022-12-31",
		Description:   "Test description",
		CoverUrl:      "http://example.com/cover.jpg",
		Language: []Language{
			{Code: "en"},
			{Code: "fr"},
		},
		Publisher: []Publisher{
			{Name: "Publisher A"},
			{Name: "Publisher B"},
		},
		Author: []Author{
			{Name: "Author A"},
			{Name: "Author B"},
		},
		Category: []Category{
			{Name: "Category A"},
			{Name: "Category B"},
		},
	}

	err = store.CreatePublication(publication)
	assert.NoError(t, err)
	assert.NotEmpty(t, publication.ID)

	// validate the publication
	err = publication.Validate()
	assert.NoError(t, err)

	// retrieve the user by uuid
	storedPublication, err := store.GetPublication(publication.UUID)
	assert.NoError(t, err)

	// check equality
	if !checkPublicationEquality(*publication, *storedPublication) {
		t.Error("Fetched publication does not match the created publication")
	}

	// update the title
	storedPublication.Title = "Updated Test Publication"
	err = store.UpdatePublication(storedPublication)
	assert.NoError(t, err)

	// retrieve the publication by ID and validate the updated title
	updatedPublication, err := store.GetPublication(publication.UUID)
	assert.NoError(t, err)
	assert.Equal(t, updatedPublication.Title, storedPublication.Title)

	// count publications
	pubCount, err := store.CountPublications()
	assert.NoError(t, err)
	assert.Equal(t, 1, int(pubCount))

	// delete the publication
	err = store.DeletePublication(updatedPublication)
	assert.NoError(t, err)

	// retrieve the publication by ID and ensure it's not found
	_, err = store.GetPublication(publication.UUID)
	assert.Error(t, err)

	// check publishers. They were not deleted.
	publishers, err := store.GetPublishers()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(publishers))

}

func TestFindByCategory(t *testing.T) {

	// Create test publications
	publication1 := &Publication{
		Title:         "Test Publication 1",
		UUID:          uuid.New().String(),
		DatePublished: "2022-12-31",
		Description:   "Test description",
		Category: []Category{
			{Name: "Category A"},
		},
	}
	publication2 := &Publication{
		Title:         "Test Publication 2",
		UUID:          uuid.New().String(),
		DatePublished: "2022-12-31",
		Description:   "Test description",
		Category: []Category{
			{Name: "Category B"},
		},
	}

	err := store.CreatePublication(publication1)
	assert.NoError(t, err)

	err = store.CreatePublication(publication2)
	assert.NoError(t, err)

	// Test FindByCategory
	publications, err := store.FindPublicationsByCategory("Category B", 1, 10)
	assert.NoError(t, err)

	// Ensure the correct number of publications is retrieved
	if len(publications) != 1 {
		t.Errorf("Expected 1 publication, got %d", len(publications))
	}

	if !checkPublicationEquality(publications[0], *publication2) {
		t.Error("Fetched publication does not match the created publication")
	}

	// Clean up test data
	err = store.DeletePublication(publication1)
	if err != nil {
		t.Errorf("Error deleting publication 1: %s", err.Error())
	}
	err = store.DeletePublication(publication2)
	if err != nil {
		t.Errorf("Error deleting publication 2: %s", err.Error())
	}
}

func TestGetPublicationByAuthor(t *testing.T) {

	// Create test publications
	publication1 := &Publication{
		Title:         "Test Publication 1",
		UUID:          uuid.New().String(),
		DatePublished: "2022-12-31",
		Description:   "Test description",
		Author: []Author{
			{Name: "Author A"},
		},
	}
	publication2 := &Publication{
		Title:         "Test Publication 2",
		UUID:          uuid.New().String(),
		DatePublished: "2022-12-31",
		Description:   "Test description",
		Author: []Author{
			{Name: "Author B"},
		},
	}

	err := store.CreatePublication(publication1)
	assert.NoError(t, err)

	err = store.CreatePublication(publication2)
	assert.NoError(t, err)

	publications, err := store.FindPublicationsByAuthor("Author B", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(publications))

	// Ensure the retrieved publication matches the created publication
	if !checkPublicationEquality(*publication2, publications[0]) {
		t.Error("Fetched publication does not match the created publication")
	}

	// check authors
	authors, err := store.GetAuthors()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(authors))

	// Clean up test data
	err = store.DeletePublication(publication1)
	assert.NoError(t, err)

	err = store.DeletePublication(publication2)
	assert.NoError(t, err)
}

func TestFindByPublisher(t *testing.T) {

	// Create test publications
	publication1 := &Publication{
		Title:         "Test Publication 1",
		UUID:          uuid.New().String(),
		DatePublished: "2022-12-31",
		Description:   "Test description",
		Publisher: []Publisher{
			{Name: "Publisher A"},
		},
	}
	publication2 := &Publication{
		Title:         "Test Publication 2",
		UUID:          uuid.New().String(),
		DatePublished: "2022-12-31",
		Description:   "Test description",
		Publisher: []Publisher{
			{Name: "Publisher B"},
		},
	}

	err := store.CreatePublication(publication1)
	assert.NoError(t, err)

	err = store.CreatePublication(publication2)
	assert.NoError(t, err)

	publications, err := store.FindPublicationsByPublisher("Publisher B", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(publications))

	// check publishers
	publishers, err := store.GetPublishers()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(publishers))

	// Ensure the retrieved publication matches the created publication
	if !checkPublicationEquality(*publication2, publications[0]) {
		t.Error("Fetched publication does not match the created publication")
	}

	// Clean up test data
	err = store.DeletePublication(publication1)
	assert.NoError(t, err)

	err = store.DeletePublication(publication2)
	assert.NoError(t, err)
}

func TestFindByLanguage(t *testing.T) {

	// Create test publications
	publication1 := &Publication{
		Title:         "Test Publication 1",
		UUID:          uuid.New().String(),
		DatePublished: "2022-12-31",
		Description:   "Test description",
		Language: []Language{
			{Code: "en"},
		},
	}
	publication2 := &Publication{
		Title:         "Test Publication 2",
		UUID:          uuid.New().String(),
		DatePublished: "2022-12-31",
		Description:   "Test description",
		Language: []Language{
			{Code: "fr"},
		},
	}

	err := store.CreatePublication(publication1)
	assert.NoError(t, err)

	err = store.CreatePublication(publication2)
	assert.NoError(t, err)

	publications, err := store.FindPublicationsByLanguage("en", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(publications))

	// check languages
	languages, err := store.GetLanguages()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(languages))

	// Ensure the retrieved publication matches the created publication
	if !checkPublicationEquality(*publication1, publications[0]) {
		t.Error("Fetched publication does not match the created publication")
	}

	// Clean up test data
	err = store.DeletePublication(publication1)
	assert.NoError(t, err)

	err = store.DeletePublication(publication2)
	assert.NoError(t, err)
}

func TestCreate2PublicationsWithSameCategory(t *testing.T) {

	publication := &Publication{
		Title:         "Test Publication 1",
		UUID:          uuid.New().String(),
		ContentType:   "application/epub+zip",
		DatePublished: "2022-12-31",
		Description:   "Test description",
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
		Title:         "Test Publication 2",
		UUID:          uuid.New().String(),
		ContentType:   "text/plain",
		DatePublished: "2022-12-31",
		Description:   "Test description",
		Language: []Language{
			{Code: "en"},
			{Code: "fr"},
		},
		Publisher: []Publisher{
			{Name: "Publisher A"},
			{Name: "Publisher B"},
		},
		Author: []Author{
			{Name: "Author A"},
			{Name: "Author B"},
		},
		Category: []Category{
			{Name: "Category A"},
			{Name: "Category B"},
		},
	}

	// create both publications
	err := store.CreatePublication(publication)
	assert.NoError(t, err)

	err = store.CreatePublication(publication2)
	assert.NoError(t, err)

	// list publications
	publications, err := store.ListPublications(1, 5)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(publications))

	// get categories
	categories, err := store.GetCategories()
	assert.NoError(t, err)

	// Ensure the correct number of categories is retrieved
	// if len(categories) != 2 {
	// 	t.Errorf("Expected 2 categories, got %d", len(categories))
	// }

	category1 := Category{Name: "Category A"}
	category2 := Category{Name: "Category B"}

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

	// find by type
	publications, err = store.FindPublicationsByType("text/plain", 1, 5)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(publications))

	// find by title
	publications, err = store.FindPublicationsByTitle("Test Publication 1", 1, 5)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(publications))

	// Clean up test data
	err = store.DeletePublication(publication)
	assert.NoError(t, err)

	err = store.DeletePublication(publication2)
	assert.NoError(t, err)
}
