package stor

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Language struct {
	gorm.Model
	Code string `gorm:"size:2;index:idx_code;unique"`
}

func (l *Language) BeforeSave(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(clause.OnConflict{DoNothing: false})
	return
}

type Publisher struct {
	gorm.Model
	Name string `gorm:"unique"`
}

func (l *Publisher) BeforeSave(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(clause.OnConflict{DoNothing: false})
	return
}

type Author struct {
	gorm.Model
	Name string `gorm:"unique"`
}

func (l *Author) BeforeSave(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(clause.OnConflict{DoNothing: false})
	return
}

type Category struct {
	gorm.Model
	Name string `gorm:"unique"`
}

func (l *Category) BeforeSave(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(clause.OnConflict{DoNothing: false})
	return
}

type Publication struct {
	gorm.Model
	Title           string
	UUID            string `gorm:"uniqueIndex"`
	DatePublication time.Time
	Description     string
	CoverUrl        string
	Language        []Language  `gorm:"many2many:publication_language;save_association:false"`
	Publisher       []Publisher `gorm:"many2many:publication_publisher;"`
	Author          []Author    `gorm:"many2many:publication_author;"`
	Category        []Category  `gorm:"many2many:publication_category;"`
}

// CreatePublication creates a new publicatio
func (stor *Stor) CreatePublication(publication *Publication) error {
	if err := stor.db.Create(publication).Error; err != nil {
		return err
	}

	return nil
}

// UpdatePublication updates a publication
func (stor *Stor) UpdatePublication(publication *Publication) error {
	if err := stor.db.Save(publication).Error; err != nil {
		return err
	}

	return nil
}

// DeletePublication deletes a publication
// TODO: delete many2many link if empty
// category,publisher,... items are not deleted if is only linked with this deleted publication
func (stor *Stor) DeletePublication(publication *Publication) error {
	if err := stor.db.Delete(publication).Error; err != nil {
		return err
	}

	return nil
}

func (stor *Stor) preloadPublication() *gorm.DB {
	return stor.db.Preload("Author").Preload("Publisher").Preload("Language").Preload("Category")
}

// GetPublicationByID retrieves a publication by ID
func (stor *Stor) GetPublicationByUUID(uuid string) (*Publication, error) {
	var publication Publication
	if err := stor.preloadPublication().Where("UUID = ?", uuid).First(&publication).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Publication not found")
		}
		return nil, err
	}

	return &publication, nil
}

func (stor *Stor) GetAllPublications(page int, pageSize int) ([]Publication, error) {
	var publications []Publication
	offset := (page - 1) * pageSize

	if err := stor.preloadPublication().Offset(offset).Limit(pageSize).Find(&publications).Error; err != nil {
		return nil, err
	}

	return publications, nil
}

// GetPublicationByID retrieves a publication by Title
func (stor *Stor) GetPublicationByTitle(title string) (*Publication, error) {
	var publication Publication
	if err := stor.preloadPublication().Where("Title = ?", title).First(&publication).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Publication not found")

		}
		return nil, err
	}

	return &publication, nil
}

// GetPublicationByCategory retrieves publications by category
func (stor *Stor) GetPublicationsByCategory(category string) ([]Publication, error) {
	var publications []Publication
	if err := stor.preloadPublication().Joins("JOIN publication_category ON publication_category.publication_id = publications.id").
		Joins("JOIN categories ON categories.id = publication_category.category_id").
		Where("categories.name = ?", category).Find(&publications).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Publication not found")
		}
		return nil, err
	}

	return publications, nil
}

// GetPublicationByAuthor retrieves publications by author
func (stor *Stor) GetPublicationsByAuthor(author string) ([]Publication, error) {
	var publications []Publication
	if err := stor.preloadPublication().Joins("JOIN publication_author ON publication_author.publication_id = publications.id").
		Joins("JOIN authors ON authors.id = publication_author.author_id").
		Where("authors.name = ?", author).Find(&publications).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Publication not found")
		}
		return nil, err
	}

	return publications, nil
}

// GetPublicationByAuthor retrieves publications by publisher
func (stor *Stor) GetPublicationsByPublisher(publisher string) ([]Publication, error) {
	var publications []Publication
	if err := stor.preloadPublication().Joins("JOIN publication_publisher ON publication_publisher.publication_id = publications.id").
		Joins("JOIN publishers ON publishers.id = publication_publisher.publisher_id").
		Where("publishers.name = ?", publisher).Find(&publications).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Publication not found")
		}
		return nil, err
	}

	return publications, nil
}

// GetPublicationByAuthor retrieves publications by language
func (stor *Stor) GetPublicationsByLanguage(code string) ([]Publication, error) {
	var publications []Publication
	if err := stor.preloadPublication().Joins("JOIN publication_language ON publication_language.publication_id = publications.id").
		Joins("JOIN languages ON languages.id = publication_language.language_id").
		Where("languages.code = ?", code).Find(&publications).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Publication not found")
		}
		return nil, err
	}

	return publications, nil
}

func (stor *Stor) GetCategories() ([]Category, error) {
	var categories []Category
	if err := stor.db.Find(&categories).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no categories found")
		}
		return nil, err
	}

	return categories, nil
}

func (stor *Stor) GetAuthors() ([]Author, error) {
	var authors []Author
	if err := stor.db.Find(&authors).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no authors found")
		}
		return nil, err
	}

	return authors, nil
}

func (stor *Stor) GetPublishers() ([]Publisher, error) {
	var publishers []Publisher
	if err := stor.db.Find(&publishers).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no publishers found")
		}
		return nil, err
	}

	return publishers, nil
}

func (stor *Stor) GetLanguages() ([]Language, error) {
	var languages []Language
	if err := stor.db.Find(&languages).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no languages found")
		}
		return nil, err
	}

	return languages, nil
}
