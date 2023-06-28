package stor

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestUserCRUD(t *testing.T) {

	// Create a new user
	user := &User{
		UUID:        gofakeit.UUID(),
		Name:        "Pierre ler",
		Email:       gofakeit.Email(),
		Pass:        "password123",
		LcpHintMsg:  "Hint",
		LcpPassHash: "Hash",
		SessionId:   gofakeit.UUID(),
	}

	err := stor.CreateUser(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)

	// Retrieve user by email and validate
	readUser, err := stor.GetUserByEmailAndPass(user.Email, user.Pass)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, readUser.ID)
	assert.Equal(t, user.Name, readUser.Name)
	assert.Equal(t, user.Email, readUser.Email)
	assert.Equal(t, user.Pass, readUser.Pass)
	assert.Equal(t, user.LcpHintMsg, readUser.LcpHintMsg)
	assert.Equal(t, user.LcpPassHash, readUser.LcpPassHash)
	assert.Equal(t, user.SessionId, readUser.SessionId)

	// Update user name
	user.Name = "Jane Doe"
	err = stor.UpdateUser(user)
	assert.NoError(t, err)

	// Retrieve user by ID and validate updated name
	readUser, err = stor.GetUserByUUID(user.UUID)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, readUser.Name)

	// Delete user
	err = stor.DeleteUser(user)
	assert.NoError(t, err)

	// Retrieve user by ID and ensure it's not found
	deletedUser, err := stor.GetUserByUUID(user.UUID)
	assert.Error(t, err)
	assert.Nil(t, deletedUser)
}
