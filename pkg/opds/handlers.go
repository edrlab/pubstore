// Copyright 2023 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package opds

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/edrlab/pubstore/pkg/lcp"
	"github.com/edrlab/pubstore/pkg/stor"
)

// GetAuthenticationDoc returns an OPDS authentication document
func GetAuthenticationDoc(w http.ResponseWriter, _ *http.Request) {
	jsonData := `
	{
		"id": "org.edrlab.pubstore",
		"title": "LOGIN",
		"description": "PUBSTORE LOGIN",
		"links": [
		  {"rel": "logo", "href": "` + publicBaseUrl + `/static/images/edrlab-logo.jpeg", "type": "image/jpeg", "width": 90, "height": 90}
		],
		"authentication": [
		  {
			"type": "http://opds-spec.org/auth/oauth/password",
			"links": [
			  {"rel": "authenticate", "href": "` + publicBaseUrl + `/opds/token", "type": "application/json"}
			]
		  }
		]
	  }`

	w.Header().Set("Content-Type", "application/opds-authentication+json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(jsonData))
}

// GetCatalog returns the full OPDS Catalog
func (opds *Opds) GetCatalog(w http.ResponseWriter, r *http.Request) {

	// TODO add pagination in opds feed
	opdsFeed, err := opds.GenerateOpdsFeed(1, 50)
	if err != nil {
		fmt.Fprintf(w, "opds feed : %v!", err)
	}
	// Encode the publication as JSON and write it to the response
	w.Header().Set("Content-Type", "application/opds+json")
	err = json.NewEncoder(w).Encode(opdsFeed)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetPublication returns an OPDS Publication
func (o *Opds) GetPublication(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	credential, ok := ctx.Value(CredentialContext).(string)
	authentified := false
	if ok {
		authentified = true
	}

	storPublication, ok := ctx.Value("publication").(*stor.Publication)
	if !ok {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	pub, err := convertToOpdsPublication(storPublication)
	if err != nil {
		http.Error(w, "Failed to convert publication", http.StatusInternalServerError)
		return
	}

	defer func() {
		// Encode the publication as JSON and write it to the response
		w.Header().Set("Content-Type", "application/opds+json")
		err = json.NewEncoder(w).Encode(pub)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}()

	if authentified {
		// add authentified acquisition link

		user, err := o.Store.GetUserByEmail(credential)
		if err != nil {
			// http.Error(w, "Failed to get user", http.StatusInternalServerError)
			fmt.Println("Failed to get user : " + credential)
			pub.Links = append(pub.Links, publicationAcquisitionLinkChoice("notAuthentified", storPublication.UUID, "", "", time.Time{}, time.Time{}))
			return
		}

		transactions, err := o.getTransactionFromUserAndPubUUID(user, storPublication.UUID)
		if err != nil {
			pub.Links = append(pub.Links, publicationAcquisitionLinkChoice("authentified", storPublication.UUID, "", "", time.Time{}, time.Time{}))
			return
		}

		lsdStatus, err := lcp.GetStatusDocument(o.Config.LCPServer, transactions.LicenceId, user.Email, user.TextHint, user.HPassphrase)
		if err != nil {
			lsdStatus = &lcp.LsdStatus{}
		}
		pub.Links = append(pub.Links, publicationAcquisitionLinkChoice("authentifiedAndBorrowed", storPublication.UUID, lsdStatus.StatusCode, user.HPassphrase, lsdStatus.StartDate, lsdStatus.EndDate))

	} else {
		// add borrow link
		pub.Links = append(pub.Links, publicationAcquisitionLinkChoice("notAuthentified", storPublication.UUID, "", "", time.Time{}, time.Time{}))
	}
}

// GetPublicationBorrow redirects to GetPublication
func (opds *Opds) GetPublicationBorrow(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(CredentialContext).(string)
	if !ok {
		GetAuthenticationDoc(w, r)
		return
	}
	storPublication, ok := ctx.Value("publication").(*stor.Publication)
	if !ok {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	http.Redirect(w, r, "/opds/publication/"+storPublication.UUID, http.StatusFound)
}

// GetPublicationLoan is currently empty
func (opds *Opds) GetPublicationLoan(w http.ResponseWriter, r *http.Request) {

	/*
		ctx := r.Context()
		credential, ok := ctx.Value(CredentialContext).(string)
		if !ok {
			opds.GetAuthenticationDoc(w, r)
			return
		}
		user, err := opds.stor.GetUserByEmail(credential)
		if err != nil {
			http.Error(w, "Failed to get transaction", http.StatusInternalServerError)
			return
		}

		storPublication, ok := ctx.Value("publication").(*stor.Publication)
		if !ok {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		licenceBytes, err := lcp.LicenceLoan(storPublication.UUID, user.UUID, user.Email, user.TextHint, user.HPassphrase, 100, 2000, time.Now(), time.Now().AddDate(0, 0, 7))
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		w.Header().Set("Content-Type", "application/vnd.readium.lcp.license.v1.0+json")
		io.Copy(w, bytes.NewReader(licenceBytes))
	*/
}

// GetPublicationLicense returns the LCP license attached to a user / publication tuple
func (o *Opds) GetPublicationLicense(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	credential, ok := ctx.Value(CredentialContext).(string)
	if !ok {
		GetAuthenticationDoc(w, r)
		return
	}
	user, err := o.Store.GetUserByEmail(credential)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	storPublication, ok := ctx.Value("publication").(*stor.Publication)
	if !ok {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	transaction, err := o.Store.GetTransactionByUserAndPublication(user.ID, storPublication.ID)
	if err != nil {
		http.Error(w, "Failed to get transaction", http.StatusInternalServerError)
		return
	}

	licenceBytes, err := lcp.GetFreshLicense(o.Config.LCPServer, transaction.LicenceId, user.Email, user.TextHint, user.HPassphrase)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.readium.lcp.license.v1.0+json")
	io.Copy(w, bytes.NewReader(licenceBytes))
}

// GetBookshelf returns a personal bookshelf feed
func (opds *Opds) GetBookshelf(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	credential, ok := ctx.Value(CredentialContext).(string)
	if !ok {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	opdsFeed, err := opds.GenerateBookshelfFeed(credential)
	if err != nil {
		fmt.Println("Bookshelf : " + err.Error())
	}

	// Encode the publication as JSON and write it to the response
	w.Header().Set("Content-Type", "application/opds+json")
	err = json.NewEncoder(w).Encode(opdsFeed)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
