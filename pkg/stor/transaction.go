package stor

import (
	"errors"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID        uint
	User          User
	PublicationID uint
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

// CreateTransaction creates a new transaction
func (stor *Stor) CreateTransactionWithUUID(pubUUID, userUUID, licenceUUID string) error {

	publication, err := stor.GetPublicationByUUID(pubUUID)
	if err != nil {
		return errors.New("can't get publication")
	}

	user, err := stor.GetUserByUUID(userUUID)
	if err != nil {
		return errors.New("can't get user")
	}

	transaction := &Transaction{
		UserID:        user.ID,
		PublicationID: publication.ID,
		LicenceId:     licenceUUID,
	}

	return stor.CreateTransaction(transaction)
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Transaction not found")
		}
		return nil, err
	}

	return &transaction, nil
}

func (stor *Stor) GetTransactionByUserAndPublication(userID, publicationID uint) (*Transaction, error) {
	var transaction Transaction
	if err := stor.db.Preload("User").Preload("Publication").Where("user_id = ?", userID).Where("publication_id = ?", publicationID).Order("created_at DESC").First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Transaction not found")
		}
		return nil, err
	}

	return &transaction, nil
}

func (stor *Stor) GetTransactionsByUserID(userID uint) (*[]Transaction, error) {
	var transaction []Transaction
	if err := stor.db.Preload("User").Preload("Publication").Where("user_id = ?", userID).Order("created_at DESC").Find(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Transaction not found")
		}
		return nil, err
	}

	return &transaction, nil
}
