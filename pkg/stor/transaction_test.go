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
		UUID:       gofakeit.UUID(),
		Name:       "Pierre 1er",
		Email:      gofakeit.Email(),
		Password:   "password",
		TextHint:   "hint",
		Passphrase: "passphrase",
	}

	err := store.CreateUser(user)
	assert.NoError(t, err)

	// create a new publication
	publication := &Publication{
		Title:         "Test Publication",
		UUID:          gofakeit.UUID(),
		DatePublished: "2022-12-31",
		Description:   "Test description",
		CoverUrl:      "http://example.com/cover.jpg",
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

	err = store.CreatePublication(publication)
	if err != nil {
		t.Errorf("Error creating publication: %s", err.Error())
	}

	// create a new transaction
	transaction := &Transaction{
		UserID:        user.ID,
		PublicationID: publication.ID,
		LicenceId:     gofakeit.UUID(),
	}

	err = store.CreateTransaction(transaction)
	assert.NoError(t, err)
	assert.NotZero(t, transaction.ID)

	// read the transaction by licence ID
	readTransaction, err := store.GetTransactionByLicence(transaction.LicenceId)
	assert.NoError(t, err)
	assert.Equal(t, transaction.ID, readTransaction.ID)
	assert.Equal(t, transaction.UserID, readTransaction.UserID)
	assert.Equal(t, transaction.PublicationID, readTransaction.PublicationID)
	assert.Equal(t, transaction.LicenceId, readTransaction.LicenceId)

	// update the transaction
	transaction.LicenceId = gofakeit.UUID()
	err = store.UpdateTransaction(transaction)
	assert.NoError(t, err)

	// verify the updated transaction
	updatedTransaction, err := store.GetTransactionByLicence(transaction.LicenceId)
	assert.NoError(t, err)
	assert.Equal(t, transaction.LicenceId, updatedTransaction.LicenceId)

	// retrieve the transaction by userID and publicationID
	readTransaction2, err := store.GetTransactionByUserAndPublication(transaction.UserID, transaction.PublicationID)
	assert.NoError(t, err)
	assert.Equal(t, readTransaction2.UserID, updatedTransaction.UserID)
	assert.Equal(t, readTransaction2.PublicationID, updatedTransaction.PublicationID)
	assert.Equal(t, readTransaction2.LicenceId, updatedTransaction.LicenceId)

	// retrieves the array to transactions made by the user
	transactions, err := store.FindTransactionsByUser(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, readTransaction2.UserID, (*transactions)[0].UserID)
	assert.Equal(t, readTransaction2.PublicationID, (*transactions)[0].PublicationID)
	assert.Equal(t, readTransaction2.LicenceId, (*transactions)[0].LicenceId)

	// delete the transaction
	err = store.DeleteTransaction(transaction)
	assert.NoError(t, err)

	// verify that the transaction is deleted
	_, err = store.GetTransactionByLicence(transaction.LicenceId)
	assert.Error(t, err)

	// delete the publication
	err = store.DeletePublication(publication)
	assert.NoError(t, err)

	// delete the user
	err = store.DeleteUser(user)
	assert.NoError(t, err)
}
