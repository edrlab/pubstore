// Copyright 2023 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package opds

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/edrlab/pubstore/pkg/lcp"
	"github.com/edrlab/pubstore/pkg/stor"
)

// publicationAcquisitionLinkChoice
// choice is "authentified" || "notAuthentified" || "authentifiedAndBorrowed"
func publicationAcquisitionLinkChoice(choice string, pubUUID, statusCode, lcpHashedPassphrase string, startDate, endDate time.Time) Link {

	if choice == "authentified" {
		return Link{
			Type: "application/vnd.readium.lcp.license.v1.0+json",
			Rel:  "http://opds-spec.org/acquisition/borrow",
			Href: publicBaseUrl + "/opds/publication/" + pubUUID + "/loan",
			Properties: &Properties{
				Availability: &Availability{
					Status: "available",
				},
				IndirectAcquisition: []Link{
					{
						Type: "application/vnd.readium.lcp.license.v1.0+json",
						Child: []Link{
							{
								Type: "application/epub+zip",
							},
						},
					},
				},
			},
		}

	} else if choice == "notAuthentified" {
		return Link{
			Type: "application/opds-publication+json",
			Rel:  "http://opds-spec.org/acquisition/borrow",
			Href: publicBaseUrl + "/opds/publication/" + pubUUID + "/borrow",
			Properties: &Properties{
				Availability: &Availability{
					Status: "available",
				},
				IndirectAcquisition: []Link{
					{
						Type: "application/vnd.readium.lcp.license.v1.0+json",
						Child: []Link{
							{
								Type: "application/epub+zip",
							},
						},
					},
				},
			},
		}

	}
	return Link{
		Type: "application/vnd.readium.lcp.license.v1.0+json",
		Rel:  "http://opds-spec.org/acquisition",
		Href: publicBaseUrl + "/opds/publication/" + pubUUID + "/license",
		Properties: &Properties{
			Availability: &Availability{
				Status:    statusCode,
				StartDate: &startDate,
				EndDate:   &endDate,
			},
			LcpHashedPassphrase: lcpHashedPassphrase,
			IndirectAcquisition: []Link{
				{
					Type: "application/vnd.readium.lcp.license.v1.0+json",
					Child: []Link{
						{
							Type: "application/epub+zip",
						},
					},
				},
			},
		},
	}

}

// convertToOpdsPublication converts a stored Publication to an OPDS Publication
func convertToOpdsPublication(storPublication *stor.Publication) (Publication, error) {
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
			Published:  storPublication.DatePublished,
		},
		Links: []Link{
			{
				Rel:  "self",
				Href: publicBaseUrl + "/opds/publication/" + storPublication.UUID,
				Type: "application/opds+json",
			},
		},
		Images: getImages(storPublication.CoverUrl),
	}

	return publication, nil
}

// getAuthorNames creates a comma separated string out of an array of authors
func getAuthorNames(authors []stor.Author) string {
	var names []string
	for _, author := range authors {
		names = append(names, author.Name)
	}
	return strings.Join(names, ", ")
}

// getLanguageCode returns the first language
// TODO: an OPDS Publication supports x languages as an array -> update the model and mapping
func getLanguageCode(languages []stor.Language) string {
	var codes []string
	for _, language := range languages {
		codes = append(codes, language.Code)
	}
	if len(codes) > 0 {
		return codes[0]
	}
	return ""
}

// getImages returns an Image object out of cover info
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

// GenerateOpdsFeed create an OPDS feed from all existing publications
func (opds *Opds) GenerateOpdsFeed(page, pageSize int) (Root, error) {

	publications, err := opds.Store.ListPublications(page, pageSize)
	if err != nil {
		return Root{}, errors.New("Error fetching publications:" + err.Error())
	}

	root := Root{
		Metadata: MetadataFeed{
			Title: "Pubstore OPDS Feed",
		},
		Links: []Link{
			{
				Rel:  "self",
				Href: publicBaseUrl + "/opds/catalog",
				Type: "application/opds+json",
			},
			{
				Rel:  "http://opds-spec.org/shelf",
				Href: publicBaseUrl + "/opds/bookshelf",
				Type: "application/opds+json",
			},
		},
		Publications: make([]Publication, len(publications)),
	}

	for i, storPub := range publications {
		root.Publications[i], err = convertToOpdsPublication(&storPub)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println(root.Links)

	return root, nil
}

// getTransactionFromUserAndPubUUID
func (opds *Opds) getTransactionFromUserAndPubUUID(user *stor.User, pubUUID string) (*stor.Transaction, error) {
	if user == nil {
		return nil, errors.New("no user")
	}

	pub, err := opds.Store.GetPublication(pubUUID)
	if err != nil {
		return nil, err
	}

	transaction, err := opds.Store.GetTransactionByUserAndPublication(user.ID, pub.ID)
	if err != nil {
		return nil, err

	}
	return transaction, nil
}

// GenerateBookshelfFeed create a personal OPDS feed
func (opds *Opds) GenerateBookshelfFeed(credential string) (Root, error) {

	user, err := opds.Store.GetUserByEmail(credential)
	if err != nil {
		return Root{}, err
	}
	transactions, err := opds.Store.FindTransactionsByUser(user.ID)
	if err != nil {
		return Root{}, err
	}

	var publicationOpdsView []Publication = make([]Publication, len(*transactions))
	var lsdStatus []*lcp.LsdStatus = make([]*lcp.LsdStatus, len(*transactions))
	for i, transactionStor := range *transactions {
		// TODO: try to find a better solution than a loop on http requests
		lsdStatus[i], err = lcp.GetStatusDocument(opds.Config.LCPServer, transactionStor.LicenceId, transactionStor.User.Email, transactionStor.User.TextHint, transactionStor.User.HPassphrase)
		if err != nil {
			lsdStatus[i] = &lcp.LsdStatus{}
		}
		publicationOpdsView[i], err = convertToOpdsPublication(&transactionStor.Publication)
		if err != nil {
			publicationOpdsView[i] = Publication{}
		}
	}

	root := Root{
		Metadata: MetadataFeed{
			Title: "bookshelf",
		},
		Links: []Link{
			{
				Rel:  "self",
				Href: publicBaseUrl + "/opds/bookshelf",
				Type: "application/opds+json",
			},
			{
				Rel:  "http://opds-spec.org/shelf",
				Href: publicBaseUrl + "/opds/bookshelf",
				Type: "application/opds+json",
			},
		},
		Publications: publicationOpdsView,
	}

	for i, status := range lsdStatus {
		root.Publications[i].Links = append(root.Publications[i].Links, publicationAcquisitionLinkChoice("authentified", root.Publications[i].Metadata.Identifier, status.StatusCode, user.HPassphrase, status.StartDate, status.EndDate))
	}

	return root, nil
}
