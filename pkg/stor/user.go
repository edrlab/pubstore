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

// CreateUser creates a new user
func (stor *Stor) CreateUser(user *User) error {
	if err := stor.db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

// UpdateUser updates a user
func (stor *Stor) UpdateUser(user *User) error {
	if err := stor.db.Save(user).Error; err != nil {
		return err
	}

	return nil
}

// DeleteUser deletes a user
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
