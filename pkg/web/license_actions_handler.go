// Copyright 2025 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package web

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/edrlab/pubstore/pkg/lcp"
	"github.com/go-chi/chi/v5"
)

// licenseRegisterHandler sends a register request to the License Server
func (web *Web) licenseRegisterHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Registering license...")
	web.executeAction(w, r, "register")
}

// licenseRenewHandler sends a renew request to the License Server
func (web *Web) licenseRenewHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Renewing license...")
	web.executeAction(w, r, "renew")
}

// licenseReturnHandler sends a return request to the License Server
func (web *Web) licenseReturnHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Returning license...")
	web.executeAction(w, r, "return")
}

// licenseRevokeHandler sends a return request to the License Server
func (web *Web) licenseRevokeHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Revoking license...")
	web.executeAction(w, r, "revoke")
}

// executeAction is a helper function to execute register, renew, return and revoke actions
func (web *Web) executeAction(w http.ResponseWriter, r *http.Request, action string) {
	licenseID := chi.URLParam(r, "id")

	// get the transaction associated with this license
	transaction, err := web.Store.GetTransactionByLicense(licenseID)
	if err != nil {
		message := url.QueryEscape("Failed to get transaction by license: " + err.Error())
		http.Error(w, message, 500)
		return
	}

	publicationUUID := transaction.Publication.UUID
	currentLicenseUpdated := transaction.LicenseUpdated
	currentStatus := transaction.Status
	log.Println("Current license status: " + currentStatus)

	err = lcp.ExecuteAction(web.Config.LCPServer, transaction, action)
	if err != nil {
		log.Printf("%s: %v", action, err)
		message := url.QueryEscape(err.Error())
		http.Redirect(w, r, fmt.Sprintf("/bookshelf/publications/%s?license=%s&err=%s", publicationUUID, licenseID, message), http.StatusFound)
		return
	}

	// the transaction has been updated during the lcp action
	if transaction.Status != currentStatus || transaction.LicenseUpdated.After(currentLicenseUpdated) {
		// store the updated transaction
		if transaction.Status != currentStatus {
			log.Println("License status has been updated: " + transaction.Status)
		}
		if transaction.LicenseUpdated.After(currentLicenseUpdated) {
			log.Println("The license has been updated: " + transaction.LicenseUpdated.String())
		}
		err = web.Store.UpdateTransaction(transaction)
		if err != nil {
			log.Printf("%v", err)
			message := url.QueryEscape("Failed to update transaction: " + err.Error())
			http.Redirect(w, r, fmt.Sprintf("/bookshelf/publications/%s?license=%s&err=%s", publicationUUID, licenseID, message), http.StatusFound)
			return
		}
	}

	// re-direct to the bookshelf publication page with success message
	var successMessage string
	switch action {
	case "register":
		successMessage = "License registered successfully"
	case "renew":
		successMessage = "License renewed successfully"
	case "return":
		successMessage = "License returned successfully"
	case "revoke":
		successMessage = "License revoked successfully"
	}
	message := url.QueryEscape(successMessage)
	http.Redirect(w, r, fmt.Sprintf("/bookshelf/publications/%s?license=%s&success=%s", publicationUUID, licenseID, message), http.StatusFound)
}
