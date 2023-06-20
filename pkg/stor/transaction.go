package stor

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID        int
	User          User
	PublicationID int
	Publication   Publication
	LicenceId     string
}

// CreateTransaction creates a new transaction
func (stor *Stor) CreateTransaction(transaction *Transaction) error {
	if err := stor.db.Create(transaction).Error; err != nil {
		return err
	}

	return nil
}

// UpdateTransaction updates a transaction
func (stor *Stor) UpdateTransaction(transaction *Transaction) error {
	if err := stor.db.Save(transaction).Error; err != nil {
		return err
	}

	return nil
}

// DeleteTransaction deletes a transaction
func (stor *Stor) DeleteTransaction(transaction *Transaction) error {
	if err := stor.db.Delete(transaction).Error; err != nil {
		return err
	}

	return nil
}

func (stor *Stor) GetTransactionByLicenceId(transactionID string) (*Transaction, error) {
	var transaction Transaction
	if err := stor.db.Preload("User").Preload("Publication").Where("licence_id = ?", transactionID).First(&transaction).Error; err != nil {
		return nil, err
	}

	return &transaction, nil
}
