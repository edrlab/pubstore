// Copyright 2025 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package web

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/edrlab/pubstore/pkg/lcp"
	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/go-chi/chi/v5"
)

// createLicense sends a request to the License Server and returns the license to the caller
func (web *Web) createLicense(w http.ResponseWriter, r *http.Request) {

	// get request params
	pubUUID := chi.URLParam(r, "id")
	printParam := r.URL.Query().Get("printRights")
	copyParam := r.URL.Query().Get("copyRights")
	startParam := r.URL.Query().Get("startDate")
	endParam := r.URL.Query().Get("endDate")

	licenseReq := lcp.LicenseRequest{}

	// sanitize params
	var print, copy int32
	var start, end time.Time
	var i int
	var err error
	if i, err = strconv.Atoi(printParam); err != nil {
		fmt.Println(err.Error())
		print = int32(web.Config.PrintLimit)
	} else {
		print = int32(i)
	}
	if i, err = strconv.Atoi(copyParam); err != nil {
		fmt.Println(err.Error())
		copy = int32(web.Config.CopyLimit)
	} else {
		copy = int32(i)
	}
	// start & end params may be empty strings. In this case their time representation keep a zero value
	if startParam != "" {
		start, err = time.Parse(time.RFC3339, startParam)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	if endParam != "" {
		end, err = time.Parse(time.RFC3339, endParam)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	// get user information
	user := web.getUserByCookie(r)

	// get publication information
	errMessage := "License acquisition failed: "

	publication, err := web.GetPublication(pubUUID)
	if err != nil {
		acquisitionFailure(w, r, pubUUID, errMessage+err.Error())
		return
	}

	licenseReq.PublicationID = pubUUID

	// negative values for print and copy are considered unconstrained
	if print >= 0 {
		licenseReq.Print = &print
	}
	if copy >= 0 {
		licenseReq.Copy = &copy
	}
	// zero start and end are considered void (therefore unconstrained)
	if !start.IsZero() {
		licenseReq.Start = &start
	}
	if !end.IsZero() {
		licenseReq.End = &end
	}
	licenseReq.UserID = user.UUID
	licenseReq.UserName = user.Name
	licenseReq.UserEmail = user.Email
	licenseReq.UserEncrypted = []string{"email"}
	licenseReq.Profile = web.Config.EncryptionProfile
	licenseReq.TextHint = user.TextHint
	licenseReq.PassHash = user.HPassphrase

	license, _, err := lcp.GenerateLicense(web.Config.LCPServer, licenseReq)
	if err != nil {
		acquisitionFailure(w, r, pubUUID, errMessage+err.Error())
		return
	}

	// Extract the link to the status document
	var statusDocLink lcp.Link
	for _, l := range license.Links {
		if l.Rel == "status" {
			statusDocLink = l
			break
		}
	}

	noLimit := int32(-1) // -1 stored for no print/copy limits
	if license.Rights.Copy == nil {
		license.Rights.Copy = &noLimit
	}
	if license.Rights.Print == nil {
		license.Rights.Print = &noLimit
	}

	// create a transaction
	transaction := &stor.Transaction{
		UserID:        user.ID,
		PublicationID: publication.ID,
		LicenseId:     license.ID,
		Status:        "ready", // license is ready to be used
		StatusDocLink: statusDocLink.Href,
		Print:         *license.Rights.Print,
		Copy:          *license.Rights.Copy,
		Start:         license.Rights.Start,
		End:           license.Rights.End,
	}

	err = web.CreateTransaction(transaction)
	if err != nil {
		acquisitionFailure(w, r, pubUUID, errMessage+err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/bookshelf/publications/%s?license=%s&success=%s", pubUUID, license.ID, url.QueryEscape("License created successfully")), http.StatusFound)
}

// acquisitionFailure is a helper function that redirects to the publication page with an error message
func acquisitionFailure(w http.ResponseWriter, r *http.Request, pubID string, message string) {
	http.Redirect(w, r, fmt.Sprintf("/catalog/publications/%s?err=%s", pubID, url.QueryEscape(message)), http.StatusFound)
}
