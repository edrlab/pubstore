// Copyright 2023 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package api

import (
	"bytes"
	"io"
	"net/http"

	"github.com/edrlab/pubstore/pkg/lcp"
)

// @Summary Get a fresh license
// @Description Get a fresh licen from a license id
// @Tags licences
// @Accept -
// @Produce json
// @Param user body User true "User object"
// @Success 200 {object} A fresh license is returned
// @Router /licenses [get]
func (a *Api) getFreshLicense(w http.ResponseWriter, r *http.Request) {

	transaction := fromTransContext(r.Context())

	licenceBytes, err := lcp.GetFreshLicense(a.Config.LCPServer, transaction)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.readium.lcp.license.v1.0+json")
	io.Copy(w, bytes.NewReader(licenceBytes))

}
