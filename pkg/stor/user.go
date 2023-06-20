package stor

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name        string
	Email       string `gorm:"uniqueIndex"`
	Pass        string
	LcpHintMsg  string
	LcpPassHash string
	SessionId   string `gorm:"uniqueIndex"`
}

// CreatePublication creates a new publication
func (stor *Stor) CreateUser(user *User) error {
	if err := stor.db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

// UpdatePublication updates a publication
func (stor *Stor) UpdateUser(user *User) error {
	if err := stor.db.Save(user).Error; err != nil {
		return err
	}

	return nil
}

// DeletePublication deletes a publication
// TODO: delete many2many link if empty
// category,publisher,... items are not deleted if is only linked with this deleted publication
func (stor *Stor) DeleteUser(user *User) error {
	if err := stor.db.Delete(user).Error; err != nil {
		return err
	}

	return nil
}

func (stor *Stor) GetUserBySessionId(sessionId string) (*User, error) {
	var user User
	if err := stor.db.Where("session_id = ?", sessionId).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
