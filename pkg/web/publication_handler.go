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
	successMsg := r.URL.Query().Get("success")
	licenseID := r.URL.Query().Get("license")
	licenseReady := false
	licenseUsable := false
	IsBookshelfPage := false

	publicationStor, err := web.Store.GetPublication(pubUUID)
	if err != nil {
		http.ServeFile(w, r, "static/404.html")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	viewTransaction := view.TransactionView{}

	// get the user (if authenticated)
	var userName string
	user := web.getUserByCookie(r)
	if user != nil {
		userName = user.Name
	}

	// detect if the request comes from the bookshelf page
	if strings.Contains(r.URL.Path, "bookshelf") {
		IsBookshelfPage = true

		// get the transaction ID from the query parameter
		if licenseID == "" {
			message := url.QueryEscape("No license ID provided")
			http.Error(w, message, 400)
			return
		}

		// get the transaction for this license (if any)
		transaction, err := web.Store.GetTransactionByLicense(licenseID)
		if err != nil {
			message := url.QueryEscape("Failed to get transaction by license id: " + err.Error())
			http.Error(w, message, 500)
			return
		}

		currentLicenseUpdated := transaction.LicenseUpdated
		currentStatus := transaction.Status
		log.Println("Transaction", transaction.ID, "with license", licenseID, "has current status", currentStatus)

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
			//maxEnd = lsdStatus.MaxEnd.Format("2006-01-02 15:04:05")
			maxEnd = lsdStatus.MaxEnd.Format("2 January 2006, 15:04:05 MST")
		} else {
			maxEnd = "unknown"
		}
		viewTransaction.LicenseMaxEnd = maxEnd

		log.Println("maxEnd:", maxEnd)

		if transaction.Status == "ready" {
			licenseReady = true
		}
		if transaction.Status == "ready" || transaction.Status == "active" {
			licenseUsable = true
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
		"successMsg":            successMsg,
		"title":                 publicationView.Title,
		"uuid":                  publicationView.UUID,
		"altId":                 publicationView.AltId,
		"format":                publicationView.Format,
		"description":           publicationView.Description,
		"coverUrl":              publicationView.CoverUrl,
		"authors":               publicationView.Authors,
		"publishers":            publicationView.Publishers,
		"licenseFound":          bool(viewTransaction.PublicationUUID != ""),
		"licenseReady":          licenseReady,
		"licenseFoundAndUsable": licenseUsable,
		"transaction":           viewTransaction,
		"IsBookshelfPage":       IsBookshelfPage,
	}
	err = goview.Render(w, http.StatusOK, "publication", goviewModel)
	if err != nil {
		fmt.Fprintf(w, "Render index error: %v!", err)
	}

}
