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
// DatePublished is a string: we do not process its value as a dateTime (or a simpler date, which is more complex to validate)
type Publication struct {
	gorm.Model
	UUID          string      `json:"uuid" validate:"omitempty,uuid4_rfc4122" gorm:"uniqueIndex"`
	Title         string      `json:"title" gorm:"index"`
	ContentType   string      `json:"content_type" gorm:"index"`
	DatePublished string      `json:"date_published"`
	Description   string      `json:"description"`
	CoverUrl      string      `json:"cover_url"`
	Language      []Language  `json:"language" gorm:"many2many:publication_language;"`
	Publisher     []Publisher `json:"publisher" gorm:"many2many:publication_publisher;"`
	Author        []Author    `json:"author" gorm:"many2many:publication_author;"`
	Category      []Category  `json:"category" gorm:"many2many:publication_category;"`
}

// TODO : remove gorm.Model from these tables

type Language struct {
	gorm.Model
	Code string `json:"code" gorm:"size:2;uniqueIndex"`
}

// TODO : verify if the BeforeSave clauses are good

func (l *Language) BeforeSave(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(clause.OnConflict{
		DoNothing: false,
		Columns:   []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"code"}), //map[string]interface{}{"code": "EXCLUDED.code"}),
	})
	return
}

type Publisher struct {
	gorm.Model
	Name string `json:"name" gorm:"uniqueIndex"`
}

func (l *Publisher) BeforeSave(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(clause.OnConflict{
		DoNothing: false,
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}), //map[string]interface{}{"name": "EXCLUDED.name"}),
	})
	return
}

type Author struct {
	gorm.Model
	Name string `json:"name" gorm:"uniqueIndex"`
}

func (l *Author) BeforeSave(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(clause.OnConflict{
		DoNothing: false,
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}), //map[string]interface{}{"name": "EXCLUDED.name"}),
	})
	return
}

type Category struct {
	gorm.Model
	Name string `json:"name" gorm:"uniqueIndex"`
}

func (l *Category) BeforeSave(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(clause.OnConflict{
		DoNothing: false,
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}), //map[string]interface{}{"name": "EXCLUDED.name"}),
	})
	return
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
	return s.db.Session(&gorm.Session{FullSaveAssociations: true}).Model(&Publication{}).Preload("Author").Preload("Publisher").Preload("Language").Preload("Category")
	// return s.db.Model(&Publication{}).Preload("Author").Preload("Publisher").Preload("Language").Preload("Category")
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

// FindPublicationsByCategory retrieves publications by category
func (s *Store) FindPublicationsByCategory(category string, page int, pageSize int) ([]Publication, error) {
	var publications []Publication
	offset := (page - 1) * pageSize
	if offset < 0 {
		return publications, errors.New("invalid pagination")
	}
	return publications, s.preloadPublication().Joins("JOIN publication_category ON publication_category.publication_id = publications.id").
		Joins("JOIN categories ON categories.id = publication_category.category_id").
		Where("categories.name = ?", category).Order(clause.OrderByColumn{Column: clause.Column{Table: "publications", Name: "updated_at"}, Desc: true}).Offset(offset).Limit(pageSize).Find(&publications).Error
}

// FindPublicationsByAuthor retrieves publications by author
func (s *Store) FindPublicationsByAuthor(author string, page int, pageSize int) ([]Publication, error) {
	var publications []Publication
	offset := (page - 1) * pageSize
	if offset < 0 {
		return publications, errors.New("invalid pagination")
	}
	return publications, s.preloadPublication().Joins("JOIN publication_author ON publication_author.publication_id = publications.id").
		Joins("JOIN authors ON authors.id = publication_author.author_id").
		Where("authors.name = ?", author).Order(clause.OrderByColumn{Column: clause.Column{Table: "publications", Name: "updated_at"}, Desc: true}).Offset(offset).Limit(pageSize).Find(&publications).Error
}

// FindPublicationsByPublisher retrieves publications by publisher
func (s *Store) FindPublicationsByPublisher(publisher string, page int, pageSize int) ([]Publication, error) {
	var publications []Publication
	offset := (page - 1) * pageSize
	if offset < 0 {
		return publications, errors.New("invalid pagination")
	}
	return publications, s.preloadPublication().Joins("JOIN publication_publisher ON publication_publisher.publication_id = publications.id").
		Joins("JOIN publishers ON publishers.id = publication_publisher.publisher_id").
		Where("publishers.name = ?", publisher).Order(clause.OrderByColumn{Column: clause.Column{Table: "publications", Name: "updated_at"}, Desc: true}).Offset(offset).Limit(pageSize).Find(&publications).Error
}

// FindPublicationsByLanguage retrieves publications by language
func (s *Store) FindPublicationsByLanguage(code string, page int, pageSize int) ([]Publication, error) {
	var publications []Publication
	offset := (page - 1) * pageSize
	if offset < 0 {
		return publications, errors.New("invalid pagination")
	}
	return publications, s.preloadPublication().Joins("JOIN publication_language ON publication_language.publication_id = publications.id").
		Joins("JOIN languages ON languages.id = publication_language.language_id").
		Where("languages.code = ?", code).Order(clause.OrderByColumn{Column: clause.Column{Table: "publications", Name: "updated_at"}, Desc: true}).Offset(offset).Limit(pageSize).Find(&publications).Error
}

// Count returns the publication count
func (s *Store) CountPublications() (int64, error) {
	var count int64
	return count, s.db.Model(Publication{}).Count(&count).Error
}

// GetCategories lists available categories
func (s *Store) GetCategories() ([]Category, error) {
	var categories []Category
	return categories, s.db.Order(clause.OrderByColumn{Column: clause.Column{Name: "name"}, Desc: false}).Find(&categories).Error
}

// GetAuthors lists available authors
func (s *Store) GetAuthors() ([]Author, error) {
	var authors []Author
	return authors, s.db.Order(clause.OrderByColumn{Column: clause.Column{Name: "name"}, Desc: false}).Find(&authors).Error
}

// GetPublishers lists available publishers
func (s *Store) GetPublishers() ([]Publisher, error) {
	var publishers []Publisher
	return publishers, s.db.Order(clause.OrderByColumn{Column: clause.Column{Name: "name"}, Desc: false}).Find(&publishers).Error
}

// GetLanguages lists available languages
func (s *Store) GetLanguages() ([]Language, error) {
	var languages []Language
	return languages, s.db.Order(clause.OrderByColumn{Column: clause.Column{Name: "code"}, Desc: false}).Find(&languages).Error
}
