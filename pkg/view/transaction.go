package view

import (
	"fmt"
	"time"

	"github.com/edrlab/pubstore/pkg/lcp"
	"github.com/edrlab/pubstore/pkg/stor"
)

type TransactionView struct {
	TransactionID             string
	TransactionDate           time.Time
	PublicationUUID           string
	PublicationTitle          string
	PublicationAuthor         string
	PublicationCoverUrl       string
	PublicationPrintRights    string
	PublicationCopyRights     string
	PublicationStartDate      time.Time
	PublicationEndDate        time.Time
	LicenseStatusMessage      string
	LicenseStatusCode         string
	LicenseEndPotentialRights time.Time
}

func (view *View) GetTransactionViewFromTransactionStor(transaction *stor.Transaction) *TransactionView {

	var publicationAuthor string
	publication, err := view.stor.GetPublicationByUUID(transaction.Publication.UUID)
	if err == nil && len(publication.Author) > 0 {
		publicationAuthor = publication.Author[0].Name
	}

	statusMessage, statusCode, endPotentialRights, printRights, copyRights, startDate, endDate, err := lcp.GetLsdStatus(transaction.LicenceId, transaction.User.Email, transaction.User.LcpHintMsg, transaction.User.LcpPassHash)

	return &TransactionView{
		TransactionID:             fmt.Sprintf("%d", transaction.ID),
		TransactionDate:           transaction.CreatedAt,
		PublicationUUID:           transaction.Publication.UUID,
		PublicationTitle:          transaction.Publication.Title,
		PublicationAuthor:         publicationAuthor,
		PublicationCoverUrl:       publication.CoverUrl,
		PublicationPrintRights:    fmt.Sprintf("%d", printRights),
		PublicationCopyRights:     fmt.Sprintf("%d", copyRights),
		PublicationStartDate:      startDate,
		PublicationEndDate:        endDate,
		LicenseStatusMessage:      statusMessage,
		LicenseStatusCode:         statusCode,
		LicenseEndPotentialRights: endPotentialRights,
	}
}
