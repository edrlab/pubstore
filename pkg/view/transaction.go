package view

import (
	"strconv"

	"github.com/edrlab/pubstore/pkg/stor"
)

type TransactionView struct {
	// TransactionID             string
	// TransactionDate           time.Time
	PublicationUUID        string
	PublicationTitle       string
	PublicationAuthor      string
	PublicationFormat      string
	PublicationCoverUrl    string
	PublicationPrintRights string
	PublicationCopyRights  string
	PublicationStartDate   string
	PublicationEndDate     string
	LicenseStatusMessage   string
	LicenseStatus          string
	LicenseMaxEnd          string
	LicenseId              string
}

func (view *View) GetTransactionViewFromTransactionStor(transaction *stor.Transaction) *TransactionView {

	var publicationAuthor string
	publication, err := view.Store.GetPublication(transaction.Publication.UUID)
	if err == nil && len(publication.Author) > 0 {
		publicationAuthor = publication.Author[0].Name
	}

	var start, end, copy, print string
	unknown := "unknown"
	if transaction.Start != nil {
		start = transaction.Start.Format("2006-01-02 15:04:05")
	} else {
		start = unknown
	}
	if transaction.End != nil {
		end = transaction.End.Format("2006-01-02 15:04:05")
	} else {
		end = unknown
	}
	if transaction.Copy >= 0 {
		copy = strconv.Itoa(int(transaction.Copy))
	} else {
		copy = unknown
	}
	if transaction.Print >= 0 {
		print = strconv.Itoa(int(transaction.Print))
	} else {
		print = unknown
	}

	return &TransactionView{
		PublicationUUID:        transaction.Publication.UUID,
		PublicationTitle:       transaction.Publication.Title,
		PublicationFormat:      contentTypeToFormat(transaction.Publication.ContentType),
		LicenseId:              transaction.LicenseId,
		PublicationAuthor:      publicationAuthor,
		PublicationCoverUrl:    publication.CoverUrl,
		PublicationPrintRights: print,
		PublicationCopyRights:  copy,
		PublicationStartDate:   start,
		PublicationEndDate:     end,
		LicenseStatus:          transaction.Status,
	}
}
