// Copyright 2023 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package stor

import (
	"errors"

	"github.com/edrlab/pubstore/pkg/lcp"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User defines the user entity
type User struct {
	gorm.Model
	UUID        string `json:"uuid" validate:"omitempty,uuid4_rfc4122" gorm:"uniqueIndex"`
	Name        string `json:"name"`
	Email       string `json:"email" gorm:"index"`
	Password    string `json:"password" gorm:"-"`
	HPassword   string `json:"hpassword"`
	TextHint    string `json:"text_hint"`
	Passphrase  string `json:"passphrase" gorm:"-"`
	HPassphrase string `json:"hpassphrase"`
	SessionId   string `json:"-" gorm:"index"`
	// does not work : `gorm:"uniqueIndex:idx_name_not_empty,where:name IS NOT NULL"`
	// sessionId is empty at first and then filed with a unique UUID v4 when the user is connecting
}

// Validate checks required fields and values
func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

// BeforeSave creates a hash of the user password and lcp passphrase.
// This applies only if the password and/or passphrase are set.
// Note: the clear password and passphrase are not saved.
func (u *User) BeforeSave(tx *gorm.DB) error {

	// generate a hash of the user password
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("failed to hash the user password; " + err.Error())
		}
		u.HPassword = string(hashedPassword)
	}
	// generate a hash of the lcp passphrase
	if u.Passphrase != "" {
		u.HPassphrase = lcp.HashPassphrase(u.Passphrase)
	}
	return nil
}

// BeforeCreate creates user uuid if missing
func (u *User) BeforeCreate(tx *gorm.DB) error {

	if u.Password == "" {
		return errors.New("missing user authentication password")
	}
	if u.Passphrase == "" {
		return errors.New("missing user LCP passphrase")
	}
	// generate a user UUID if empty
	if u.UUID == "" {
		u.UUID = uuid.New().String()
	}
	return nil
}

// BeforeUpdate checks the uuid
func (u *User) BeforeUpdate(tx *gorm.DB) error {

	// generate a user UUID if empty
	if u.UUID == "" {
		return errors.New("missing user UUID")
	}
	return nil
}

// CreateUser creates a new user
func (s *Store) CreateUser(user *User) error {
	return s.db.Create(user).Error
}

// UpdateUser updates a user
func (s *Store) UpdateUser(user *User) error {
	return s.db.Save(user).Error
}

// GetUser returns a user, found by uuid
func (s *Store) GetUser(uuid string) (*User, error) {
	var user User
	return &user, s.db.Where("uuid = ?", uuid).First(&user).Error
}

// GetUserBySession returns a user, found by session id
func (s *Store) GetUserBySession(sessionId string) (*User, error) {
	var user User
	return &user, s.db.Where("session_id = ?", sessionId).First(&user).Error
}

// GetUserByEmail returns a user, found by email
func (s *Store) GetUserByEmail(email string) (*User, error) {
	var user User
	err := s.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

// DeleteUser deletes a user
func (s *Store) DeleteUser(user *User) error {
	return s.db.Delete(user).Error
}

// ListUsers lists users, with pagination
func (s *Store) ListUsers(page, pageSize int) ([]User, error) {
	users := []User{}
	// page starts at 1, pageSize >= 1
	offset := (page - 1) * pageSize
	if offset < 0 {
		return users, errors.New("invalid page or pageSize")
	}
	// result sorted to assure the same order for each request
	return users, s.db.Offset(offset).Limit(pageSize).Order("id ASC").Find(&users).Error
}

// CountUsers returns the user count
func (s *Store) CountUsers() (int64, error) {
	var count int64
	return count, s.db.Model(User{}).Count(&count).Error
}
