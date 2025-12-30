// Copyright 2025 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.
package web

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/edrlab/pubstore/pkg/lcp"
	"github.com/go-chi/chi/v5"
)

// freshLicenseHandler handles a request for a fresh license
// This is a public license gateway
func (web *Web) freshLicenseHandler(w http.ResponseWriter, r *http.Request) {

	licenseID := chi.URLParam(r, "id")

	// get the transaction
	transaction, err := web.Store.GetTransactionByLicense(licenseID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(404), 404)
		return
	}

	_, data, err := lcp.GetFreshLicense(web.Config.LCPServer, transaction)
	if err != nil {
		message := url.QueryEscape("Failed to generate a fresh license: " + err.Error())
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+transaction.Publication.Title+".lcpl")
	w.Header().Set("Content-Type", "application/vnd.readium.lcp.license.v1.0+json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))

	io.Copy(w, bytes.NewReader(data))
}
