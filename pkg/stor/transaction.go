// Copyright 2023 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package stor

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID        uint // implicit foreign key to the related user
	User          User
	PublicationID uint // implicit foreign key to the related publication
	Publication   Publication
	LicenceId     string
}

// CreateTransaction creates a new transaction
func (s *Store) CreateTransaction(transaction *Transaction) error {
	return s.db.Create(transaction).Error
}

// UpdateTransaction updates a transaction
func (s *Store) UpdateTransaction(transaction *Transaction) error {
	return s.db.Save(transaction).Error
}

// GetTransactionByLicense retrieves a transaction using its licenseID
func (s *Store) GetTransactionByLicence(licenseID string) (*Transaction, error) {
	var transaction Transaction
	return &transaction, s.db.Preload("User").Preload("Publication").Where("licence_id = ?", licenseID).First(&transaction).Error
}

// GetTransactionByUserAndPublication retrieves a transaction using its userID and publicationID
func (s *Store) GetTransactionByUserAndPublication(userID, publicationID uint) (*Transaction, error) {
	var transaction Transaction
	return &transaction, s.db.Preload("User").Preload("Publication").Where("user_id = ?", userID).Where("publication_id = ?", publicationID).Order("created_at DESC").First(&transaction).Error
}

// FindTransactionsByUser retrieves the array to transactions made by a specific user
func (s *Store) FindTransactionsByUser(userID uint) (*[]Transaction, error) {
	var transaction []Transaction
	return &transaction, s.db.Preload("User").Preload("Publication").Where("user_id = ?", userID).Order("created_at DESC").Find(&transaction).Error
}

// DeleteTransaction deletes a transaction
func (s *Store) DeleteTransaction(transaction *Transaction) error {
	return s.db.Delete(transaction).Error
}
