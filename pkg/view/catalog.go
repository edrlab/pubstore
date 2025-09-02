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
	Authors    []string
	Publishers []string
	Languages  []string
	Categories []string
}

type CatalogView struct {
	FacetsView
	Publications   []PublicationCatalogView
	NbPages        string
	NbPublications string
}

func (view *View) GetCatalogFacetsView() *FacetsView {
	var facets FacetsView

	if authorArray, err := view.Store.GetAuthors(); err != nil {
		fmt.Println(err)
		facets.Authors = make([]string, 0)
	} else {
		facets.Authors = make([]string, len(authorArray))
		for i, element := range authorArray {
			facets.Authors[i] = element.Name
		}
	}

	if publisherArray, err := view.Store.GetPublishers(); err != nil {
		fmt.Println(err)
		facets.Publishers = make([]string, 0)
	} else {
		facets.Publishers = make([]string, len(publisherArray))
		for i, element := range publisherArray {
			facets.Publishers[i] = element.Name
		}
	}

	if languageArray, err := view.Store.GetLanguages(); err != nil {
		fmt.Println(err)
		facets.Languages = make([]string, 0)
	} else {
		facets.Languages = make([]string, len(languageArray))
		for i, element := range languageArray {
			facets.Languages[i] = element.Code
		}
	}

	if categoryArray, err := view.Store.GetCategories(); err != nil {
		fmt.Println(err)
		facets.Categories = make([]string, 0)
	} else {
		facets.Categories = make([]string, len(categoryArray))
		for i, element := range categoryArray {
			facets.Categories[i] = element.Name
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
				var author string
				if len(element.Author) > 0 {
					author = element.Author[0].Name
				}
				publications[i] = PublicationCatalogView{CoverHref: element.CoverUrl, Title: element.Title, Author: author, UUID: element.UUID, Format: value}
			}
		}
	case "author":
		if pubs, err = view.Store.FindPublicationsByAuthor(value, page, pageSize); err != nil {
			publications = make([]PublicationCatalogView, 0)
		} else {
			publications = make([]PublicationCatalogView, len(pubs))
			for i, element := range pubs {
				publications[i] = PublicationCatalogView{CoverHref: element.CoverUrl, Title: element.Title, Author: element.Author[0].Name, UUID: element.UUID, Format: contentTypeToFormat(element.ContentType)}
			}
		}

	case "publisher":
		if pubs, err = view.Store.FindPublicationsByPublisher(value, page, pageSize); err != nil {
			publications = make([]PublicationCatalogView, 0)
		} else {
			publications = make([]PublicationCatalogView, len(pubs))
			for i, element := range pubs {
				var author = ""
				if len(element.Author) > 0 {
					author = element.Author[0].Name
				}
				publications[i] = PublicationCatalogView{CoverHref: element.CoverUrl, Title: element.Title, Author: author, UUID: element.UUID, Format: contentTypeToFormat(element.ContentType)}
			}
		}

	case "language":
		if pubs, err = view.Store.FindPublicationsByLanguage(value, page, pageSize); err != nil {
			publications = make([]PublicationCatalogView, 0)
		} else {
			publications = make([]PublicationCatalogView, len(pubs))
			for i, element := range pubs {
				var author = ""
				if len(element.Author) > 0 {
					author = element.Author[0].Name
				}
				publications[i] = PublicationCatalogView{CoverHref: element.CoverUrl, Title: element.Title, Author: author, UUID: element.UUID, Format: contentTypeToFormat(element.ContentType)}
			}
		}

	case "category":
		if pubs, err = view.Store.FindPublicationsByCategory(value, page, pageSize); err != nil {
			publications = make([]PublicationCatalogView, 0)
		} else {
			publications = make([]PublicationCatalogView, len(pubs))
			for i, element := range pubs {
				var author = ""
				if len(element.Author) > 0 {
					author = element.Author[0].Name
				}
				publications[i] = PublicationCatalogView{CoverHref: element.CoverUrl, Title: element.Title, Author: author, UUID: element.UUID, Format: contentTypeToFormat(element.ContentType)}
			}
		}

	case "search":
		if pubs, err = view.Store.FindPublicationsByTitle(value, page, pageSize); err != nil {
			publications = make([]PublicationCatalogView, 0)
		} else {
			publications = make([]PublicationCatalogView, len(pubs))
			for i, element := range pubs {
				var author = ""
				if len(element.Author) > 0 {
					author = element.Author[0].Name
				}
				publications[i] = PublicationCatalogView{CoverHref: element.CoverUrl, Title: element.Title, Author: author, UUID: element.UUID, Format: contentTypeToFormat(element.ContentType)}
			}
		}

	default:
		if pubs, err = view.Store.ListPublications(page, pageSize); err != nil {
			publications = make([]PublicationCatalogView, 0)
		} else {
			publications = make([]PublicationCatalogView, len(pubs))
			for i, element := range pubs {
				var author = ""
				if len(element.Author) > 0 {
					author = element.Author[0].Name
				}
				publications[i] = PublicationCatalogView{CoverHref: element.CoverUrl, Title: element.Title, Author: author, UUID: element.UUID, Format: contentTypeToFormat(element.ContentType)}
			}
		}
	}

	return &publications, int64(len(publications))
}

func GetCatalogView(pubs *[]PublicationCatalogView, facets *FacetsView) *CatalogView {

	var catalogView CatalogView

	catalogView.Authors = facets.Authors
	catalogView.Categories = facets.Categories
	catalogView.Languages = facets.Languages
	catalogView.Publishers = facets.Publishers
	catalogView.Publications = make([]PublicationCatalogView, len(*pubs))
	for i, element := range *pubs {
		catalogView.Publications[i] = PublicationCatalogView{CoverHref: element.CoverHref, Title: element.Title, Author: element.Author, UUID: element.UUID, Format: element.Format}
	}

	return &catalogView
}

func contentTypeToFormat(contentType string) string {
	switch contentType {
	case "application/epub+zip":
		return "epub"
	case "application/pdf+lcp":
		return "pdf"
	case "application/audiobook+lcp":
		return "audiobook"
	case "application/divina+lcp":
		return "divina"
	default:
		return "unknown"
	}
}

func formatToContentType(format string) string {
	switch format {
	case "epub":
		return "application/epub+zip"
	case "pdf":
		return "application/pdf+lcp"
	case "audiobook":
		return "application/audiobook+lcp"
	case "divina":
		return "application/divina+lcp"
	default:
		return "unknown"
	}
}
