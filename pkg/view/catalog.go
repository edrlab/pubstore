// Copyright 2023 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package view

import (
	"fmt"

	"github.com/edrlab/pubstore/pkg/stor"
)

type PublicationCatalogView struct {
	CoverHref string
	Title     string
	Author    string
	UUID      string
	Format    string
}

type FacetsView struct {
	Formats    []string
}

type CatalogView struct {
	FacetsView
	Publications   []PublicationCatalogView
	NbPages        string
	NbPublications string
}

func (view *View) GetCatalogFacetsView() *FacetsView {
	var facets FacetsView

	if contentTypeArray, err := view.Store.GetContentTypes(); err != nil {
		fmt.Println(err)
		facets.Formats = make([]string, 0)
	} else {
		facets.Formats = make([]string, len(contentTypeArray))
		for i, element := range contentTypeArray {
			facets.Formats[i] = contentTypeToFormat(element)
		}
	}

	return &facets
}

func (view *View) GetCatalogPublicationsView(facet string, value string, page int, pageSize int) (*[]PublicationCatalogView, int64) {

	var publications []PublicationCatalogView
	var pubs []stor.Publication
	var err error

	switch facet {
	case "format":
		contentType := formatToContentType(value)
		if pubs, err = view.Store.FindPublicationsByType(contentType, page, pageSize); err != nil {
			publications = make([]PublicationCatalogView, 0)
		} else {
			publications = make([]PublicationCatalogView, len(pubs))
			for i, element := range pubs {
				publications[i] = PublicationCatalogView{CoverHref: element.CoverUrl, Title: element.Title, Author: element.Authors, UUID: element.UUID, Format: value}
			}
		}
	case "author":
		if pubs, err = view.Store.FindPublicationsByAuthor(value, page, pageSize); err != nil {
			publications = make([]PublicationCatalogView, 0)
		} else {
			publications = make([]PublicationCatalogView, len(pubs))
			for i, element := range pubs {
				publications[i] = PublicationCatalogView{CoverHref: element.CoverUrl, Title: element.Title, Author: element.Authors, UUID: element.UUID, Format: contentTypeToFormat(element.ContentType)}
			}
		}

	case "publisher":
		if pubs, err = view.Store.FindPublicationsByPublisher(value, page, pageSize); err != nil {
			publications = make([]PublicationCatalogView, 0)
		} else {
			publications = make([]PublicationCatalogView, len(pubs))
			for i, element := range pubs {
				publications[i] = PublicationCatalogView{CoverHref: element.CoverUrl, Title: element.Title, Author: element.Authors, UUID: element.UUID, Format: contentTypeToFormat(element.ContentType)}
			}
		}


	case "search":
		if pubs, err = view.Store.FindPublicationsByTitle(value, page, pageSize); err != nil {
			publications = make([]PublicationCatalogView, 0)
		} else {
			publications = make([]PublicationCatalogView, len(pubs))
			for i, element := range pubs {
				publications[i] = PublicationCatalogView{CoverHref: element.CoverUrl, Title: element.Title, Author: element.Authors, UUID: element.UUID, Format: contentTypeToFormat(element.ContentType)}
			}
		}

	default:
		if pubs, err = view.Store.ListPublications(page, pageSize); err != nil {
			publications = make([]PublicationCatalogView, 0)
		} else {
			publications = make([]PublicationCatalogView, len(pubs))
			for i, element := range pubs {
				publications[i] = PublicationCatalogView{CoverHref: element.CoverUrl, Title: element.Title, Author: element.Authors, UUID: element.UUID, Format: contentTypeToFormat(element.ContentType)}
			}
		}
	}

	return &publications, int64(len(publications))
}

func GetCatalogView(pubs *[]PublicationCatalogView, facets *FacetsView) *CatalogView {

	var catalogView CatalogView

	catalogView.Formats = facets.Formats
	catalogView.Publications = make([]PublicationCatalogView, len(*pubs))
	for i, element := range *pubs {
		catalogView.Publications[i] = PublicationCatalogView{CoverHref: element.CoverHref, Title: element.Title, Author: element.Author, UUID: element.UUID, Format: element.Format}
	}

	return &catalogView
}
