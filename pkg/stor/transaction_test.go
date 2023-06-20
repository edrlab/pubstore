package stor

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestTransactionCRUD(t *testing.T) {
	// Create a new transaction
	transaction := &Transaction{
		UserID:        1,
		PublicationID: 1,
		LicenceId:     gofakeit.UUID(),
	}

	err := stor.CreateTransaction(transaction)
	assert.NoError(t, err)
	assert.NotZero(t, transaction.ID)

	// Read the transaction by licence ID
	readTransaction, err := stor.GetTransactionByLicenceId(transaction.LicenceId)
	assert.NoError(t, err)
	assert.Equal(t, transaction.ID, readTransaction.ID)
	assert.Equal(t, transaction.UserID, readTransaction.UserID)
	assert.Equal(t, transaction.PublicationID, readTransaction.PublicationID)
	assert.Equal(t, transaction.LicenceId, readTransaction.LicenceId)

	// Update the transaction
	transaction.LicenceId = gofakeit.UUID()
	err = stor.UpdateTransaction(transaction)
	assert.NoError(t, err)

	// Verify the updated transaction
	updatedTransaction, err := stor.GetTransactionByLicenceId(transaction.LicenceId)
	assert.NoError(t, err)
	assert.Equal(t, transaction.LicenceId, updatedTransaction.LicenceId)

	// Delete the transaction
	err = stor.DeleteTransaction(transaction)
	assert.NoError(t, err)

	// Verify that the transaction is deleted
	deletedTransaction, err := stor.GetTransactionByLicenceId(transaction.LicenceId)
	assert.Error(t, err)
	assert.Nil(t, deletedTransaction)
}
