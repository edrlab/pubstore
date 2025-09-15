// Copyright 2025 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package web

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/edrlab/pubstore/pkg/lcp"
	"github.com/edrlab/pubstore/pkg/view"
	"github.com/foolin/goview"
	"github.com/go-chi/chi/v5"
)

// publicationHandler handles a request for a publication page
func (web *Web) publicationHandler(w http.ResponseWriter, r *http.Request) {

	pubUUID := chi.URLParam(r, "id")
	errLcp := r.URL.Query().Get("err")
	licenseOK := false
	IsBookshelfPage := false

	publicationStor, err := web.Store.GetPublication(pubUUID)
	if err != nil {
		http.ServeFile(w, r, "static/404.html")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	viewTransaction := view.TransactionView{}
	var userName string
	// detect if the request comes from the bookshelf page
	if strings.Contains(r.URL.Path, "bookshelf") {
		IsBookshelfPage = true

		// get the user (must be authenticated)
		user := web.getUserByCookie(r)
		userName = user.Name

		// get the transaction for this user and publication (if any)
		transaction, err := web.Store.GetTransactionByUserAndPublication(user.ID, publicationStor.ID)
		if err != nil {
			message := url.QueryEscape("Failed to get transaction by user and publication: " + err.Error())
			http.Error(w, message, 500)
			return
		}

		currentLicenseUpdated := transaction.LicenseUpdated
		currentStatus := transaction.Status
		log.Println("The status of the license is " + currentStatus)

		// request the status document
		lsdStatus, err := lcp.GetStatusDocument(web.Config.LCPServer, transaction)
		if err != nil {
			message := url.QueryEscape("Failed to get the status document: " + err.Error())
			http.Error(w, message, 500)
			return
		}

		viewTransaction = *web.View.GetTransactionViewFromTransactionStor(transaction)
		viewTransaction.LicenseStatusMessage = lsdStatus.StatusMessage
		var maxEnd string
		if lsdStatus.MaxEnd != nil {
			maxEnd = lsdStatus.MaxEnd.Format("2006-01-02 15:04:05")
		} else {
			maxEnd = "unknown"
		}
		viewTransaction.LicenseMaxEnd = maxEnd

		if transaction.Status == "ready" || transaction.Status == "active" {
			licenseOK = true
		}
		if currentStatus != transaction.Status || transaction.LicenseUpdated.After(currentLicenseUpdated) {
			// store the updated transaction
			log.Println("License status has been recently updated: " + transaction.Status)
			err = web.Store.UpdateTransaction(transaction)
			if err != nil {
				log.Printf("%v", err) // let's continue despite the error
			}
		}
	}

	publicationView := web.View.GetPublicationViewFromPublicationStor(publicationStor)

	goviewModel := goview.M{
		"pageTitle":             fmt.Sprintf("pubstore - %s", publicationView.Title),
		"host":                  strings.Split(web.Config.PublicBaseUrl, "://")[1],
		"userIsAuthenticated":   web.userIsAuthenticated(r),
		"userName":              userName,
		"errLcp":                errLcp,
		"title":                 publicationView.Title,
		"uuid":                  publicationView.UUID,
		"format":                publicationView.Format,
		"datePublished":         publicationView.DatePublished,
		"description":           publicationView.Description,
		"coverUrl":              publicationView.CoverUrl,
		"authors":               publicationView.Author,
		"publishers":            publicationView.Publisher,
		"languages":             publicationView.Language,
		"categories":            publicationView.Category,
		"licenseFound":          bool(viewTransaction.PublicationUUID != ""),
		"licenseFoundAndActive": licenseOK,
		"transaction":           viewTransaction,
		"IsBookshelfPage":       IsBookshelfPage,
	}
	err = goview.Render(w, http.StatusOK, "publication", goviewModel)
	if err != nil {
		fmt.Fprintf(w, "Render index error: %v!", err)
	}

}
