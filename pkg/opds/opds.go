package opds

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/edrlab/pubstore/pkg/lcp"
	"github.com/edrlab/pubstore/pkg/stor"
)

type MetadataFeed struct {
	Title string `json:"title"`
}

type Metadata struct {
	Type       string `json:"@type"`
	Title      string `json:"title"`
	Author     string `json:"author,omitempty"`
	Identifier string `json:"identifier,omitempty"`
	Language   string `json:"language,omitempty"`
	Modified   string `json:"modified,omitempty"`
}

type Link struct {
	Rel        string      `json:"rel,omitempty"`
	Href       string      `json:"href,omitempty"`
	Type       string      `json:"type,omitempty"`
	Child      []Link      `json:"child,omitempty"`
	Properties *Properties `json:"properties,omitempty"`
}

type Publication struct {
	Metadata Metadata `json:"metadata"`
	Links    []Link   `json:"links,omitempty"`
	Images   []Image  `json:"images,omitempty"`
}

type Image struct {
	Href   string `json:"href"`
	Type   string `json:"type,omitempty"`
	Height int    `json:"height,omitempty"`
	Width  int    `json:"width,omitempty"`
}

type Availability struct {
	Status    string     `json:"state,omitempty"`
	StartDate *time.Time `json:"since,omitempty"`
	EndDate   *time.Time `json:"until,omitempty"`
}

type Properties struct {
	Availability        *Availability `json:"availability,omitempty"`
	IndirectAcquisition []Link        `json:"indirectAcquisition,omitempty"`
	LcpHashedPassphrase string        `json:"lcp_hashed_passphrase,omitempty"`
}

type Root struct {
	Metadata     MetadataFeed  `json:"metadata"`
	Links        []Link        `json:"links"`
	Publications []Publication `json:"publications"`
}

func publicationAcquisitionLinkChoice(borrowLink bool, pubUUID, statusCode, lcpHashedPassphrase string, startDate, endDate time.Time) Link {

	if borrowLink {
		return Link{
			Type: "application/opds-publication+json",
			Rel:  "http://opds-spec.org/acquisition/borrow",
			Href: "http://localhost:8080/opds/v1/publication/" + pubUUID + "/loan",
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
		Type: "application/opds-publication+json",
		Rel:  "http://opds-spec.org/acquisition",
		Href: "http://localhost:8080/opds/v1/publication/" + pubUUID + "/license",
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

func convertToPublication(storPublication *stor.Publication) (Publication, error) {
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
	if len(codes) > 0 {
		return codes[0]
	}
	return ""
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
		root.Publications[i], err = convertToPublication(&storPub)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println(root.Links)

	return root, nil
}

func (opds *Opds) getTransactionFromUserAndPubUUID(user *stor.User, pubUUID string) (*stor.Transaction, error) {
	if user == nil {
		return nil, errors.New("no user")
	}

	pub, err := opds.stor.GetPublicationByUUID(pubUUID)
	if err != nil {
		return nil, err
	}

	transactions, err := opds.stor.GetTransactionByUserAndPublication(user.ID, pub.ID)
	if err != nil {
		return nil, err

	}
	return transactions, nil
}

func (opds *Opds) GenerateBookshelfFeed(credential string) (Root, error) {

	user, err := opds.stor.GetUserByEmail(credential)
	if err != nil {
		return Root{}, err
	}
	transactions, err := opds.stor.GetTransactionsByUserID(user.ID)
	if err != nil {
		return Root{}, err
	}

	var publicationOpdsView []Publication = make([]Publication, len(*transactions))
	var lsdStatus []*lcp.LsdStatus = make([]*lcp.LsdStatus, len(*transactions))
	for i, transactionStor := range *transactions {
		lsdStatus[i], err = lcp.GetLsdStatus(transactionStor.LicenceId, transactionStor.User.Email, transactionStor.User.LcpHintMsg, transactionStor.User.LcpPassHash)
		if err != nil {
			lsdStatus[i] = &lcp.LsdStatus{}
		}
		publicationOpdsView[i], err = convertToPublication(&transactionStor.Publication)
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
				Href: "http://localhost:8080/opds/v1/bookshelf",
				Type: "application/opds+json",
			},
			{
				Rel:  "http://opds-spec.org/shelf",
				Href: "http://localhost:8080/opds/v1/bookshelf",
				Type: "application/opds+json",
			},
		},
		Publications: publicationOpdsView,
	}

	for i, status := range lsdStatus {
		root.Publications[i].Links = append(root.Publications[i].Links, publicationAcquisitionLinkChoice(false, root.Publications[i].Metadata.Identifier, status.StatusCode, user.LcpPassHash, status.StartDate, status.EndDate))
	}

	return root, nil
}
