package stor

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID        string `gorm:"uniqueIndex"`
	Name        string
	Email       string `gorm:"uniqueIndex"`
	Pass        string
	LcpHintMsg  string
	LcpPassHash string
	SessionId   string // does not works : `gorm:"uniqueIndex:idx_name_not_empty,where:name IS NOT NULL"`
	// sessionId is empty at first and then filed with a unique UUID v4 when the user is connecting
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.UUID == "" {
		u.UUID = uuid.New().String()
	}

	HashedPass, err := bcrypt.GenerateFromPassword([]byte(u.Pass), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("user before create Password Hashing error " + err.Error())
	}
	u.Pass = string(HashedPass)

	return
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("User not found")
		}
		return nil, err
	}

	return &user, nil
}

func (stor *Stor) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := stor.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("User not found")
		}
		return nil, err
	}

	return &user, nil
}

func (stor *Stor) GetUserByUUID(uuid string) (*User, error) {
	var user User
	if err := stor.db.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("User not found")
		}
		return nil, err
	}

	return &user, nil
}
