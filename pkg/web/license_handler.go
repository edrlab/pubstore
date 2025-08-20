// Copyright 2022 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package web

import (
	"bytes"
	"fmt"
	"io"
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
	var print, copy int
	var start, end time.Time
	var err error
	if print, err = strconv.Atoi(printParam); err != nil {
		fmt.Println(err.Error())
		print = web.Config.PrintLimit
	}
	if copy, err = strconv.Atoi(copyParam); err != nil {
		fmt.Println(err.Error())
		copy = web.Config.CopyLimit
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

	licenseReq.PublicationID = pubUUID

	// negative values for print and copy are considered void (therefore unconstrained)
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
	licenseReq.TextHint = user.TextHint
	licenseReq.PassHash = user.HPassphrase

	errMessage := "License acquisition failed: "

	licence, err := lcp.GenerateLicense(web.Config.LCPServer, licenseReq)
	if err != nil {
		acquisitionFailure(w, r, pubUUID, errMessage+err.Error())
		return
	}

	licenseId, pubTitle, _, _, _, _, _, err := lcp.ParseLicense(licence)
	if err != nil {
		acquisitionFailure(w, r, pubUUID, errMessage+err.Error())
		return
	}

	publication, err := web.GetPublication(pubUUID)
	if err != nil {
		acquisitionFailure(w, r, pubUUID, errMessage+err.Error())
		return
	}

	transaction := &stor.Transaction{
		UserID:        user.ID,
		PublicationID: publication.ID,
		LicenceId:     licenseId,
	}

	err = web.CreateTransaction(transaction)
	if err != nil {
		acquisitionFailure(w, r, pubUUID, errMessage+err.Error())
		return
	}

	// return the license to the caller
	w.Header().Set("Content-Disposition", "attachment; filename="+pubTitle+".lcpl")
	w.Header().Set("Content-Type", "application/vnd.readium.lcp.license.v1.0+json")
	w.Header().Set("Content-Length", strconv.Itoa(len(licence)))

	io.Copy(w, bytes.NewReader(licence))
}

func acquisitionFailure(w http.ResponseWriter, r *http.Request, pubID string, message string) {
	http.Redirect(w, r, fmt.Sprintf("/catalog/publication/%s?err=%s", pubID, url.QueryEscape(message)), http.StatusFound)
}
