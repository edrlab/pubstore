package stor

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestUserCRUD(t *testing.T) {
	// Create a new user
	user := &User{
		Name:        "John Doe",
		Email:       gofakeit.Email(),
		Pass:        "password123",
		LcpHintMsg:  "Hint",
		LcpPassHash: "Hash",
		SessionId:   gofakeit.UUID(),
	}

	err := stor.CreateUser(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	// Read the user by session ID
	readUser, err := stor.GetUserBySessionId(user.SessionId)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, readUser.ID)
	assert.Equal(t, user.Name, readUser.Name)
	assert.Equal(t, user.Email, readUser.Email)
	assert.Equal(t, user.Pass, readUser.Pass)
	assert.Equal(t, user.LcpHintMsg, readUser.LcpHintMsg)
	assert.Equal(t, user.LcpPassHash, readUser.LcpPassHash)
	assert.Equal(t, user.SessionId, readUser.SessionId)

	// Update the user
	user.Name = "John Smith"
	err = stor.UpdateUser(user)
	assert.NoError(t, err)

	// Verify the updated user
	updatedUser, err := stor.GetUserBySessionId(user.SessionId)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, updatedUser.Name)

	// Delete the user
	err = stor.DeleteUser(user)
	assert.NoError(t, err)

	// Verify that the user is deleted
	deletedUser, err := stor.GetUserBySessionId(user.SessionId)
	assert.Error(t, err)
	assert.Nil(t, deletedUser)
}
