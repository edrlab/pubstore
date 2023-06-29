package opds

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/edrlab/pubstore/pkg/stor"
)

type MetadataFeed struct {
	Title string `json:"title"`
}

type Metadata struct {
	Type       string `json:"@type"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	Identifier string `json:"identifier"`
	Language   string `json:"language"`
	Modified   string `json:"modified"`
}

type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
	Type string `json:"type"`
}

type Publication struct {
	Metadata Metadata `json:"metadata"`
	Links    []Link   `json:"links"`
	Images   []Image  `json:"images"`
}

type Image struct {
	Href   string `json:"href"`
	Type   string `json:"type"`
	Height int    `json:"height,omitempty"`
	Width  int    `json:"width,omitempty"`
}

type Root struct {
	Metadata     MetadataFeed  `json:"metadata"`
	Links        []Link        `json:"links"`
	Publications []Publication `json:"publications"`
}

func ConvertToPublication(storPublication *stor.Publication) (Publication, error) {
	if storPublication == nil {
		return Publication{}, errors.New("invalid stor.Publication")
	}

	publication := Publication{
		Metadata: Metadata{
			Type:       "http://schema.org/Book",
			Title:      storPublication.Title,
			Author:     getAuthorNames(storPublication.Author),
			Identifier: storPublication.UUID,
			Language:   getLanguageCode(storPublication.Language),
			Modified:   storPublication.DatePublication.Format(time.RFC3339),
		},
		Links: []Link{
			{
				Rel:  "self",
				Href: "http://localhost:8080/opds/v1/publication/" + storPublication.UUID,
				Type: "application/opds+json",
			},
		},
		Images: getImages(storPublication.CoverUrl),
	}

	return publication, nil
}

func getAuthorNames(authors []stor.Author) string {
	var names []string
	for _, author := range authors {
		names = append(names, author.Name)
	}
	return strings.Join(names, ", ")
}

func getLanguageCode(languages []stor.Language) string {
	var codes []string
	for _, language := range languages {
		codes = append(codes, language.Code)
	}
	return codes[0]
}

func getImages(coverURL string) []Image {
	if coverURL == "" {
		return nil
	}

	images := []Image{
		{
			Href: coverURL,
			Type: "image/jpeg",
		},
	}
	return images
}

func (opds *Opds) GenerateOpdsFeed(page, pageSize int) (Root, error) {

	publications, _, err := opds.stor.GetAllPublications(page, pageSize)
	if err != nil {
		return Root{}, errors.New("Error fetching publications:" + err.Error())
	}

	root := Root{
		Metadata: MetadataFeed{
			Title: "Example listing publications",
		},
		Links: []Link{
			{
				Rel:  "self",
				Href: "http://localhost:8080/opds/v1/catalog",
				Type: "application/opds+json",
			},
			{
				Rel:  "http://opds-spec.org/shelf",
				Href: "http://localhost:8080/opds/v1/bookshelf",
				Type: "application/opds+json",
			},
		},
		Publications: make([]Publication, len(publications)),
	}

	for i, storPub := range publications {
		root.Publications[i], err = ConvertToPublication(&storPub)
		if err != nil {
			fmt.Println(err)
		}
	}

	return root, nil
}

func (opds *Opds) GenerateBookshelfFeed() (Root, error) {

	publications, _, err := opds.stor.GetAllPublications(1, 50)
	if err != nil {
		return Root{}, errors.New("Error fetching publications:" + err.Error())
	}

	// user := web.getUserByCookie(r)
	// if user == nil {
	// 	fmt.Fprintf(w, "bookshelf error")
	// 	return
	// }
	// transactions, err := web.stor.GetTransactionsByUserID(user.ID)
	// if err != nil {
	// 	fmt.Fprintf(w, "bookshelf error")
	// 	return
	// }

	// var transactionsView []*view.TransactionView = make([]*view.TransactionView, len(*transactions))
	// for i, transactionStor := range *transactions {
	// 	transactionsView[i] = web.view.GetTransactionViewFromTransactionStor(&transactionStor)
	// }

	root := Root{
		Metadata: MetadataFeed{
			Title: "bookshelf",
		},
		Links: []Link{
			{
				Rel:  "self",
				Href: "http://localhost:8080/opds/v1/bookshelf",
				Type: "application/opds+json",
			},
			{
				Rel:  "http://opds-spec.org/shelf",
				Href: "http://localhost:8080/opds/v1/bookshelf",
				Type: "application/opds+json",
			},
		},
		Publications: make([]Publication, len(publications)),
	}

	for i, storPub := range publications {
		root.Publications[i], err = ConvertToPublication(&storPub)
		if err != nil {
			fmt.Println(err)
		}
	}

	return root, nil
}
