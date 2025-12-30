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
		a.AltId != b.AltId ||
		a.Description != b.Description ||
		a.Publishers != b.Publishers ||
		a.Authors != b.Authors {
		return false
	}

	return true
}

func TestPublicationCRUD(t *testing.T) {

	// create a new publication
	publication := &Publication{
		UUID:          uuid.New().String(),
		Title:         "Test Publication",
		ContentType:   "application/zip",
		AltId:         "test-alt-id-1",
		Description:   "Test description",
		CoverUrl:      "http://example.com/cover.jpg",
		Publishers:    "Publisher A, Publisher B",
		Authors:       "Author A, Author B",
	}

	err := store.CreatePublication(publication)
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

}

func TestGetPublicationByAuthor(t *testing.T) {

	// Create test publications
	publication1 := &Publication{
		Title:         "Test Publication 1",
		UUID:          uuid.New().String(),
		AltId:         "test-alt-id-2",
		Description:   "Test description",
		Authors:       "Author A",
	}
	publication2 := &Publication{
		Title:         "Test Publication 2",
		UUID:          uuid.New().String(),
		AltId:         "test-alt-id-3",
		Description:   "Test description",
		Authors:       "Author B",
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
		AltId:         "test-alt-id-4",
		Description:   "Test description",
		Publishers:    "Publisher A",
	}
	publication2 := &Publication{
		Title:         "Test Publication 2",
		UUID:          uuid.New().String(),
		AltId:         "test-alt-id-5",
		Description:   "Test description",
		Publishers:    "Publisher B",
	}

	err := store.CreatePublication(publication1)
	assert.NoError(t, err)

	err = store.CreatePublication(publication2)
	assert.NoError(t, err)

	publications, err := store.FindPublicationsByPublisher("Publisher B", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(publications))

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
