package stor

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Language struct {
	gorm.Model
	Code string `gorm:"size:2;index:idx_code;unique"`
}

type Publisher struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type Author struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type Category struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type Publication struct {
	gorm.Model
	Title           string
	UUID            string `gorm:"uniqueIndex"`
	DatePublication time.Time
	Description     string
	Language        []Language  `gorm:"many2many:publication_language;"`
	Publisher       []Publisher `gorm:"many2many:publication_publisher;"`
	Author          []Author    `gorm:"many2many:publication_author;"`
	Category        []Category  `gorm:"many2many:publication_category;"`
}

// CreatePublication creates a new publication
func CreatePublication(publication *Publication) error {
	if err := db.Create(publication).Error; err != nil {
		return err
	}

	return nil
}

// GetPublicationByID retrieves a publication by ID
func GetPublicationByID(id uint) (*Publication, error) {
	var publication Publication
	if err := db.First(&publication, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Publication not found")
		}
		return nil, err
	}

	return &publication, nil
}

func GetAllPublications(page int, pageSize int) ([]Publication, error) {
	var publications []Publication
	offset := (page - 1) * pageSize

	if err := db.Offset(offset).Limit(pageSize).Find(&publications).Error != nil {
		return nil, result.Error
	}

	return publications, nil
}

// GetPublicationByID retrieves a publication by Title
func GetPublicationByTitle(title string) (*Publication, error) {
	var publication Publication
	if err := db.Where("Title = ?", title).First(&publication).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Publication not found")

		}
		return nil, err
	}

	return &publication, nil
}

// GetPublicationByCategory retrieves publications by category
func GetPublicationByCategory(category string) ([]Publication, error) {
	var publications []Publication
	if err := db.Preload("Category").Joins("JOIN publication_category ON publication_category.publication_id = publications.id").
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
func GetPublicationByAuthor(author string) ([]Publication, error) {
	var publications []Publication
	if err := db.Preload("Author").Joins("JOIN publication_author ON publication_author.publication_id = publications.id").
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
func GetPublicationByPublisher(publisher string) ([]Publication, error) {
	var publications []Publication
	if err := db.Preload("Publisher").Joins("JOIN publication_publisher ON publication_publisher.publication_id = publications.id").
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
func GetPublicationByLanguage(code string) ([]Publication, error) {
	var publications []Publication
	if err := db.Preload("Language").Joins("JOIN publication_language ON publication_language.publication_id = publications.id").
		Joins("JOIN languages ON languages.id = publication_language.language_id").
		Where("languages.code = ?", code).Find(&publications).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Publication not found")
		}
		return nil, err
	}

	return publications, nil
}

// UpdatePublication updates a publication
func UpdatePublication(publication *Publication) error {
	if err := db.Save(publication).Error; err != nil {
		return err
	}

	return nil
}

// DeletePublication deletes a publication
// TODO: delete many2many link if empty
// category,publisher,... items are not deleted if is only linked with this deleted publication
func DeletePublication(publication *Publication) error {
	if err := db.Delete(publication).Error; err != nil {
		return err
	}

	return nil
}

func GetCategories() ([]Category, error) {
	var categories []Category
	if err := db.Find(&categories).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("No categories found")
		}
		return nil, err
	}

	return categories, nil
}

func GetAuthors() ([]Author, error) {
    var authors []Author
    if err := db.Find(&authors).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("No authors found")
        }
        return nil, err
    }

    return authors, nil
}

func GetPublishers() ([]Publisher, error) {
    var publishers []Publisher
    if err := db.Find(&publishers).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("No publishers found")
        }
        return nil, err
    }

    return publishers, nil
}

func GetLanguages() ([]Language, error) {
    var languages []Language
    if err := db.Find(&languages).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("No languages found")
        }
        return nil, err
    }

    return languages, nil
}
