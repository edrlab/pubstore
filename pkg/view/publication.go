package view

import (
	"github.com/edrlab/pubstore/pkg/stor"
)

type PublicationView struct {
	Title         string
	UUID          string
	AltId         string
	Description   string
	CoverUrl      string
	Format        string
	Authors       string
	Publishers    string
}

func (view *View) GetPublicationViewFromPublicationStor(originalPublication *stor.Publication) *PublicationView {
	convertedPublication := PublicationView{
		Title:         originalPublication.Title,
		UUID:          originalPublication.UUID,
		AltId:         originalPublication.AltId,
		Description:   originalPublication.Description,
		CoverUrl:      originalPublication.CoverUrl,
		Authors:       originalPublication.Authors,
		Publishers:    originalPublication.Publishers,
	}

	// Convert content type to format label
	convertedPublication.Format = contentTypeToFormat(originalPublication.ContentType)

	return &convertedPublication
}
