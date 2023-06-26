package view

import (
	"fmt"

	"github.com/edrlab/pubstore/pkg/stor"
)

type TransactionView struct {
	TransactionID       string
	PublicationUUID     string
	PublicationTitle    string
	PublicationAuthor   string
	PublicationCoverUrl string
	// StartDate time.Time
	// EndDate   time.Time
	// Status    string
}

func (view *View) GetTransactionViewFromTransactionStor(transaction *stor.Transaction) *TransactionView {

	var publicationAuthor string
	publication, err := view.stor.GetPublicationByUUID(transaction.Publication.UUID)
	if err == nil && len(publication.Author) > 0 {
		publicationAuthor = publication.Author[0].Name
	}

	return &TransactionView{
		TransactionID:       fmt.Sprintf("%d", transaction.ID),
		PublicationUUID:     transaction.Publication.UUID,
		PublicationTitle:    transaction.Publication.Title,
		PublicationAuthor:   publicationAuthor,
		PublicationCoverUrl: publication.CoverUrl,
	}
}
