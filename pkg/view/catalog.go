package view

import (
	"fmt"
)

type PublicationCatalogView struct {
	CoverHref string
	Title     string
	Author    string
	UUID      string
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

	if authorArray, err := view.stor.GetAuthors(); err != nil {
		fmt.Println(err)
		facets.Authors = make([]string, 0)
	} else {
		facets.Authors = make([]string, len(authorArray))
		for i, element := range authorArray {
			facets.Authors[i] = element.Name
		}
	}

	if publisherArray, err := view.stor.GetPublishers(); err != nil {
		fmt.Println(err)
		facets.Publishers = make([]string, 0)
	} else {
		facets.Publishers = make([]string, len(publisherArray))
		for i, element := range publisherArray {
			facets.Publishers[i] = element.Name
		}
	}

	if languageArray, err := view.stor.GetLanguages(); err != nil {
		fmt.Println(err)
		facets.Languages = make([]string, 0)
	} else {
		facets.Languages = make([]string, len(languageArray))
		for i, element := range languageArray {
			facets.Languages[i] = element.Code
		}
	}

	if categoryArray, err := view.stor.GetCategories(); err != nil {
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

func (view *View) GetCatalogPublicationsView(facet string, value string) *[]PublicationCatalogView {

	var publications []PublicationCatalogView
	switch facet {

	case "author":
		if pubs, err := view.stor.GetPublicationsByAuthor(value); err != nil {
			fmt.Println(err)
			publications = make([]PublicationCatalogView, 0)
		} else {
			publications = make([]PublicationCatalogView, len(pubs))
			for i, element := range pubs {
				publications[i] = PublicationCatalogView{CoverHref: element.CoverUrl, Title: element.Title, Author: element.Author[0].Name, UUID: element.UUID}
			}
		}

	case "publisher":
		if pubs, err := view.stor.GetPublicationsByPublisher(value); err != nil {
			fmt.Println(err)
			publications = make([]PublicationCatalogView, 0)
		} else {
			publications = make([]PublicationCatalogView, len(pubs))
			for i, element := range pubs {
				var author = ""
				if len(element.Author) > 0 {
					author = element.Author[0].Name
				}
				publications[i] = PublicationCatalogView{CoverHref: element.CoverUrl, Title: element.Title, Author: author, UUID: element.UUID}
			}
		}

	case "language":
		if pubs, err := view.stor.GetPublicationsByLanguage(value); err != nil {
			fmt.Println(err)
			publications = make([]PublicationCatalogView, 0)
		} else {
			publications = make([]PublicationCatalogView, len(pubs))
			for i, element := range pubs {
				var author = ""
				if len(element.Author) > 0 {
					author = element.Author[0].Name
				}
				publications[i] = PublicationCatalogView{CoverHref: element.CoverUrl, Title: element.Title, Author: author, UUID: element.UUID}
			}
		}

	case "category":
		if pubs, err := view.stor.GetPublicationsByCategory(value); err != nil {
			fmt.Println(err)
			publications = make([]PublicationCatalogView, 0)
		} else {
			publications = make([]PublicationCatalogView, len(pubs))
			for i, element := range pubs {
				var author = ""
				if len(element.Author) > 0 {
					author = element.Author[0].Name
				}
				publications[i] = PublicationCatalogView{CoverHref: element.CoverUrl, Title: element.Title, Author: author, UUID: element.UUID}
			}
		}

	default:
		if pubs, err := view.stor.GetAllPublications(1, 50); err != nil {
			fmt.Println(err)
			publications = make([]PublicationCatalogView, 0)
		} else {
			publications = make([]PublicationCatalogView, len(pubs))
			for i, element := range pubs {
				var author = ""
				if len(element.Author) > 0 {
					author = element.Author[0].Name
				}
				publications[i] = PublicationCatalogView{CoverHref: element.CoverUrl, Title: element.Title, Author: author, UUID: element.UUID}
			}
		}
	}

	return &publications
}

func GetCatalogView(pubs *[]PublicationCatalogView, facets *FacetsView) *CatalogView {

	var catalogView CatalogView

	catalogView.Authors = facets.Authors
	catalogView.Categories = facets.Categories
	catalogView.Languages = facets.Languages
	catalogView.Publishers = facets.Publishers
	catalogView.Publications = make([]PublicationCatalogView, len(*pubs))
	for i, element := range *pubs {
		catalogView.Publications[i] = PublicationCatalogView{CoverHref: element.CoverHref, Title: element.Title, Author: element.Author, UUID: element.UUID}
	}

	return &catalogView
}
