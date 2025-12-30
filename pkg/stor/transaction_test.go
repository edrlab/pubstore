// Copyright 2023 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package stor

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestTransactionCRUD(t *testing.T) {

	// create a new user
	user := &User{
		UUID:        gofakeit.UUID(),
		Name:        "Pierre 1er",
		Email:       gofakeit.Email(),
		HPassword:   "hashed-password",
		HPassphrase: "hashed-passphrase",
		TextHint:    "hint",
	}

	err := store.CreateUser(user)
	assert.NoError(t, err)

	// create a new publication
	publication := &Publication{
		Title:         "Test Publication",
		UUID:          gofakeit.UUID(),
		AltId:         "test-alt-id-6",
		Description:   "Test description",
		CoverUrl:      "http://example.com/cover.jpg",
		Publishers:    "Test Publisher A, Test Publisher B",
		Authors:       "Test Author A, Test Author B",
	}

	err = store.CreatePublication(publication)
	if err != nil {
		t.Errorf("Error creating publication: %s", err.Error())
	}

	// create a new transaction
	transaction := &Transaction{
		UserID:        user.ID,
		PublicationID: publication.ID,
		LicenseId:     gofakeit.UUID(),
		Status:        "ready",
	}

	err = store.CreateTransaction(transaction)
	assert.NoError(t, err)
	assert.NotZero(t, transaction.ID)

	// read the transaction by license ID
	readTransaction, err := store.GetTransactionByLicense(transaction.LicenseId)
	assert.NoError(t, err)
	assert.Equal(t, transaction.ID, readTransaction.ID)
	assert.Equal(t, transaction.UserID, readTransaction.UserID)
	assert.Equal(t, transaction.PublicationID, readTransaction.PublicationID)
	assert.Equal(t, transaction.LicenseId, readTransaction.LicenseId)

	// update the transaction
	transaction.Status = "expired"
	transaction.LicenseId = gofakeit.UUID()
	err = store.UpdateTransaction(transaction)
	assert.NoError(t, err)

	// verify the updated transaction
	updatedTransaction, err := store.GetTransactionByLicense(transaction.LicenseId)
	assert.NoError(t, err)
	assert.Equal(t, transaction.LicenseId, updatedTransaction.LicenseId)

	// retrieve the transaction by userID and publicationID
	readTransaction2, err := store.GetTransactionByUserAndPublication(transaction.UserID, transaction.PublicationID)
	assert.NoError(t, err)
	assert.Equal(t, readTransaction2.UserID, updatedTransaction.UserID)
	assert.Equal(t, readTransaction2.PublicationID, updatedTransaction.PublicationID)
	assert.Equal(t, readTransaction2.LicenseId, updatedTransaction.LicenseId)

	// retrieves the array to transactions made by the user
	transactions, err := store.FindTransactionsByUser(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, readTransaction2.UserID, (*transactions)[0].UserID)
	assert.Equal(t, readTransaction2.PublicationID, (*transactions)[0].PublicationID)
	assert.Equal(t, readTransaction2.LicenseId, (*transactions)[0].LicenseId)

	// delete the transaction
	err = store.DeleteTransaction(transaction)
	assert.NoError(t, err)

	// verify that the transaction is deleted
	_, err = store.GetTransactionByLicense(transaction.LicenseId)
	assert.Error(t, err)

	// delete the publication
	err = store.DeletePublication(publication)
	assert.NoError(t, err)

	// delete the user
	err = store.DeleteUser(user)
	assert.NoError(t, err)
}
