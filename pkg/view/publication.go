package view

import (
	"github.com/edrlab/pubstore/pkg/stor"
)

type PublicationView struct {
	Title         string
	UUID          string
	DatePublished string
	Description   string
	CoverUrl      string
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
