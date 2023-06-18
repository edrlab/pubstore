package view

import (
	"fmt"

	"github.com/edrlab/pubstore/pkg/stor"
)

type PublicationView struct {
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
	Publications   []PublicationView
	NbPages        string
	NbPublications string
}

type View struct {
	stor *stor.Stor
}

func Init(s *stor.Stor) *View {
	return &View{stor: s}
}

func (view *View) GetFacetsView() *FacetsView {
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

func (view *View) GetPublicationsView(facet string, value string) *[]PublicationView {

	var publications []PublicationView
	switch facet {

	case "author":
		if pubs, err := view.stor.GetPublicationByAuthor(value); err != nil {
			fmt.Println(err)
			publications = make([]PublicationView, 0)
		} else {
			publications = make([]PublicationView, len(pubs))
			for i, element := range pubs {
				publications[i] = PublicationView{CoverHref: element.CoverUrl, Title: element.Title, Author: "", UUID: element.UUID}
			}
		}

	case "publisher":
		if pubs, err := view.stor.GetPublicationByPublisher(value); err != nil {
			fmt.Println(err)
			publications = make([]PublicationView, 0)
		} else {
			publications = make([]PublicationView, len(pubs))
			for i, element := range pubs {
				publications[i] = PublicationView{CoverHref: element.CoverUrl, Title: element.Title, Author: "", UUID: element.UUID}
			}
		}

	case "language":
		if pubs, err := view.stor.GetPublicationByLanguage(value); err != nil {
			fmt.Println(err)
			publications = make([]PublicationView, 0)
		} else {
			publications = make([]PublicationView, len(pubs))
			for i, element := range pubs {
				publications[i] = PublicationView{CoverHref: element.CoverUrl, Title: element.Title, Author: "", UUID: element.UUID}
			}
		}

	case "category":
		if pubs, err := view.stor.GetPublicationByCategory(value); err != nil {
			fmt.Println(err)
			publications = make([]PublicationView, 0)
		} else {
			publications = make([]PublicationView, len(pubs))
			for i, element := range pubs {
				publications[i] = PublicationView{CoverHref: element.CoverUrl, Title: element.Title, Author: "", UUID: element.UUID}
			}
		}

	default:
		if pubs, err := view.stor.GetAllPublications(1, 50); err != nil {
			fmt.Println(err)
			publications = make([]PublicationView, 0)
		} else {
			publications = make([]PublicationView, len(pubs))
			for i, element := range pubs {
				publications[i] = PublicationView{CoverHref: element.CoverUrl, Title: element.Title, Author: "", UUID: element.UUID}
			}
		}
	}

	return &publications
}

func GetCatalogView(pubs *[]PublicationView, facets *FacetsView) *CatalogView {

	var catalogView CatalogView

	catalogView.Authors = facets.Authors
	catalogView.Categories = facets.Categories
	catalogView.Languages = facets.Languages
	catalogView.Publishers = facets.Publishers
	catalogView.Publications = make([]PublicationView, len(*pubs))
	for i, element := range *pubs {
		catalogView.Publications[i] = PublicationView{CoverHref: element.CoverHref, Title: element.Title, Author: element.Author}
	}

	return &catalogView
}
