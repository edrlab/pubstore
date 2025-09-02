package view

import (
	"regexp"
	"time"

	"github.com/edrlab/pubstore/pkg/stor"
)

type PublicationView struct {
	Title         string
	UUID          string
	DatePublished string
	Description   string
	CoverUrl      string
	Format        string
	Author        []string
	Publisher     []string
	Category      []string
	Language      []string
}

func (view *View) GetPublicationViewFromPublicationStor(originalPublication *stor.Publication) *PublicationView {
	convertedPublication := PublicationView{
		Title:         originalPublication.Title,
		UUID:          originalPublication.UUID,
		DatePublished: originalPublication.DatePublished,
		Description:   originalPublication.Description,
		CoverUrl:      originalPublication.CoverUrl,
	}

	// Convert content type to format label
	convertedPublication.Format = contentTypeToFormat(originalPublication.ContentType)

	// Override a yyyy-mm-dd date string as a human readable formatted string
	matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, originalPublication.DatePublished)
	if matched {
		date, err := time.Parse("2006-01-02", originalPublication.DatePublished)
		if err == nil {
			convertedPublication.DatePublished = date.Format("02 Jan 2006")
		}
	}

	// Convert Language slice
	for _, language := range originalPublication.Language {
		convertedPublication.Language = append(convertedPublication.Language, language.Code)
	}

	// Convert Language slice
	for _, language := range originalPublication.Language {
		convertedPublication.Language = append(convertedPublication.Language, language.Code)
	}

	// Convert Publisher slice
	for _, publisher := range originalPublication.Publisher {
		convertedPublication.Publisher = append(convertedPublication.Publisher, publisher.Name)
	}

	// Convert Author slice
	for _, author := range originalPublication.Author {
		convertedPublication.Author = append(convertedPublication.Author, author.Name)
	}

	// Convert Category slice
	for _, category := range originalPublication.Category {
		convertedPublication.Category = append(convertedPublication.Category, category.Name)
	}

	return &convertedPublication
}
