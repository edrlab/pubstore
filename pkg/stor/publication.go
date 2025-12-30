// Copyright 2023 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package stor

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// A Publication
type Publication struct {
	gorm.Model
	UUID          string `json:"uuid" validate:"omitempty,uuid4_rfc4122" gorm:"uniqueIndex"`
	AltId         string `json:"alt_id"`
	ContentType   string `json:"content_type" gorm:"index"`
	Title         string `json:"title" gorm:"index"`
	Description   string `json:"description"`
	Authors       string `json:"authors"`
	Publishers    string `json:"publishers"`
	CoverUrl      string `json:"cover_url"`
}

// Validate checks required fields and values
func (p *Publication) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// CreatePublication creates a new publication
func (s *Store) CreatePublication(publication *Publication) error {
	return s.db.Create(publication).Error
}

// preloadPublication preloads a publication
func (s *Store) preloadPublication() *gorm.DB {
	return s.db.Session(&gorm.Session{FullSaveAssociations: true}).Model(&Publication{})
}

// GetPublication returns a publication, found by uuid
func (s *Store) GetPublication(uuid string) (*Publication, error) {
	var publication Publication
	return &publication, s.preloadPublication().Where("uuid = ?", uuid).First(&publication).Error
}

// UpdatePublication updates a publication
func (s *Store) UpdatePublication(publication *Publication) error {
	return s.db.Save(publication).Error
}

// DeletePublication deletes a publication
func (s *Store) DeletePublication(publication *Publication) error {
	return s.db.Delete(publication).Error
}

// ListPublications retrieves all publications
func (s *Store) ListPublications(page int, pageSize int) ([]Publication, error) {
	var publications []Publication

	// page starts at 1, pageSize >= 1
	offset := (page - 1) * pageSize
	if offset < 0 {
		return publications, errors.New("invalid pagination")
	}
	// result sorted to assure the same order for each request
	return publications, s.preloadPublication().Order(clause.OrderByColumn{Column: clause.Column{Name: "updated_at"}, Desc: true}).Offset(offset).Limit(pageSize).Find(&publications).Error
}

// FindPublicationsByType retrieves publications by content type
func (s *Store) FindPublicationsByType(contentType string, page int, pageSize int) ([]Publication, error) {
	var publications []Publication
	offset := (page - 1) * pageSize
	if offset < 0 {
		return publications, errors.New("invalid pagination")
	}
	return publications, s.db.Offset(offset).Limit(pageSize).Find(&publications, "content_type= ?", contentType).Error
}

// FindPublicationsByTitle retrieves publications by Title
func (s *Store) FindPublicationsByTitle(title string, page int, pageSize int) ([]Publication, error) {
	var publications []Publication
	offset := (page - 1) * pageSize
	if offset < 0 {
		return publications, errors.New("invalid pagination")
	}
	return publications, s.preloadPublication().Where("Title LIKE ?", "%"+title+"%").Order(clause.OrderByColumn{Column: clause.Column{Table: "publications", Name: "updated_at"}, Desc: true}).Offset(offset).Limit(pageSize).Find(&publications).Error
}

// FindPublicationsByAuthor retrieves publications by author
func (s *Store) FindPublicationsByAuthor(author string, page int, pageSize int) ([]Publication, error) {
	var publications []Publication
	offset := (page - 1) * pageSize
	if offset < 0 {
		return publications, errors.New("invalid pagination")
	}
	return publications, s.preloadPublication().Where("authors LIKE ?", "%"+author+"%").Order(clause.OrderByColumn{Column: clause.Column{Table: "publications", Name: "updated_at"}, Desc: true}).Offset(offset).Limit(pageSize).Find(&publications).Error
}

// FindPublicationsByPublisher retrieves publications by publisher
func (s *Store) FindPublicationsByPublisher(publisher string, page int, pageSize int) ([]Publication, error) {
	var publications []Publication
	offset := (page - 1) * pageSize
	if offset < 0 {
		return publications, errors.New("invalid pagination")
	}
	return publications, s.preloadPublication().Where("publishers LIKE ?", "%"+publisher+"%").Order(clause.OrderByColumn{Column: clause.Column{Table: "publications", Name: "updated_at"}, Desc: true}).Offset(offset).Limit(pageSize).Find(&publications).Error
}

// Count returns the publication count
func (s *Store) CountPublications() (int64, error) {
	var count int64
	return count, s.db.Model(Publication{}).Count(&count).Error
}

// GetContentTypes lists available content types
func (s *Store) GetContentTypes() ([]string, error) {
	var contentTypes []string

	// Find distinct content types in publications
	err := s.db.Model(&Publication{}).Distinct("content_type").Pluck("content_type", &contentTypes).Error
	if err != nil {
		return contentTypes, err
	}
	return contentTypes, nil
}
