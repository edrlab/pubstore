package stor

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestUserCRUD(t *testing.T) {

	// create a new user
	user := &User{
		UUID:       gofakeit.UUID(),
		Name:       "Pierre ler",
		Email:      gofakeit.Email(),
		Password:   "password",
		TextHint:   "hint",
		Passphrase: "passphrase",
		SessionId:  gofakeit.UUID(),
	}

	err := store.CreateUser(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)

	// validate the user
	err = user.Validate()
	assert.NoError(t, err)

	// retrieve user by email and validate
	readUser, err := store.GetUserByEmail(user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, readUser.ID)
	assert.Equal(t, user.Name, readUser.Name)
	assert.Equal(t, user.Email, readUser.Email)
	assert.Equal(t, user.HPassword, readUser.HPassword)
	assert.Equal(t, user.TextHint, readUser.TextHint)
	assert.Equal(t, user.HPassphrase, readUser.HPassphrase)
	assert.Equal(t, user.SessionId, readUser.SessionId)

	// update the user name
	user.Name = "Jane Doe"
	err = store.UpdateUser(user)
	assert.NoError(t, err)

	// retrieve user by ID and validate updated name
	readUser, err = store.GetUser(user.UUID)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, readUser.Name)

	// retrieve user by sessionId and validate the updated name
	readUser, err = store.GetUserBySession(user.SessionId)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, readUser.Name)

	// create a second user, with no passphrase
	user2 := &User{
		UUID:      gofakeit.UUID(),
		Name:      "Pierre ler",
		Email:     gofakeit.Email(),
		Password:  "password",
		TextHint:  "hint",
		SessionId: gofakeit.UUID(),
	}

	// check that it is not created
	err = store.CreateUser(user2)
	assert.Error(t, err)

	// add the passphrase
	user2.Passphrase = "passphrase"

	err = store.CreateUser(user2)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)

	// list users
	users, err := store.ListUsers(1, 5)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(users))

	// list users with an error
	_, err = store.ListUsers(0, 5)
	assert.Error(t, err)

	// count users
	userCount, err := store.CountUsers()
	assert.NoError(t, err)
	assert.Equal(t, 2, int(userCount))

	// delete the first user
	err = store.DeleteUser(user)
	assert.NoError(t, err)

	// retrieve user by ID and ensure it's not found
	_, err = store.GetUser(user.UUID)
	assert.Error(t, err)

	// delete the second user
	err = store.DeleteUser(user2)
	assert.NoError(t, err)
}
