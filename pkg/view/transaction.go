package view

import (
	"fmt"

	"github.com/edrlab/pubstore/pkg/lcp"
	"github.com/edrlab/pubstore/pkg/stor"
)

type TransactionView struct {
	// TransactionID             string
	// TransactionDate           time.Time
	PublicationUUID           string
	PublicationTitle          string
	PublicationAuthor         string
	PublicationCoverUrl       string
	PublicationPrintRights    string
	PublicationCopyRights     string
	PublicationStartDate      string
	PublicationEndDate        string
	LicenseStatusMessage      string
	LicenseStatusCode         string
	LicenseEndPotentialRights string
}

func (view *View) GetTransactionViewFromTransactionStor(transaction *stor.Transaction) *TransactionView {

	var publicationAuthor string
	publication, err := view.Store.GetPublication(transaction.Publication.UUID)
	if err == nil && len(publication.Author) > 0 {
		publicationAuthor = publication.Author[0].Name
	}

	// TODO: avoid fetching the Status Document in this function
	lsdStatus, err := lcp.GetStatusDocument(view.Config.LCPServer, transaction)
	if err != nil {
		fmt.Println("LSD STATUS Error from (" + transaction.LicenceId + ")")
		lsdStatus = &lcp.LsdStatus{}
	}

	return &TransactionView{
		PublicationUUID:           transaction.Publication.UUID,
		PublicationTitle:          transaction.Publication.Title,
		PublicationAuthor:         publicationAuthor,
		PublicationCoverUrl:       publication.CoverUrl,
		PublicationPrintRights:    fmt.Sprintf("%d", lsdStatus.PrintLimit),
		PublicationCopyRights:     fmt.Sprintf("%d", lsdStatus.CopyLimit),
		PublicationStartDate:      lsdStatus.StartDate.Format("2006-01-02 15:04:05"),
		PublicationEndDate:        lsdStatus.EndDate.Format("2006-01-02 15:04:05"),
		LicenseStatusMessage:      lsdStatus.StatusMessage,
		LicenseStatusCode:         lsdStatus.StatusCode,
		LicenseEndPotentialRights: lsdStatus.EndPotentialRights.Format("2006-01-02 15:04:05"),
	}
}
