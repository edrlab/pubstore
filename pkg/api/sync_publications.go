// Copyright 2023 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package api

import (
	"log"
	"net/http"
	"time"

	"github.com/edrlab/pubstore/pkg/lcp"
	"github.com/go-chi/render"
)

// @Summary Synchronize publications from an LCP Server
// @Description Update the publications table by fetching data from an LCP Server source
// @Tags publications
// @Success 200 {object} UpdateResponse "Publications updated successfully"
// @Failure 500 {object} ErrorResponse "Failed to update publications"
// @Router /synchronize [get]
func (a *Api) syncPublications(w http.ResponseWriter, r *http.Request) {

	if a.LCPServer.Version != "v2" {
		render.Render(w, r, ErrServer(
			http.ErrNotSupported,
		))
		return
	}

	// get the last synchronization time
	lastSync, err := a.Store.GetLastSyncTime()
	if err != nil {
		log.Printf("Error getting last sync time: %v", err)
		render.Render(w, r, ErrServer(err))
		return
	}
	
	// Log the sync status
	if lastSync.IsZero() {
		log.Println("No previous synchronization found, fetching all publications")
	} else {
		log.Printf("Last synchronization time: %v", lastSync)
	}

	// get publications from the LCP Server v2
	publications, err := lcp.GetPublications(a.LCPServer, lastSync)
	if err != nil {
		log.Printf("Error fetching publications from LCP Server: %v", err)
		render.Render(w, r, ErrServer(err))
		return
	}

	log.Printf("Fetched %d publications from LCP Server", len(publications))

	// save a new synchronization time
	currentTime := time.Now()
	err = a.Store.SaveSyncTime(currentTime)
	if err != nil {
		log.Printf("Error saving sync time: %v", err)
		render.Render(w, r, ErrServer(err))
		return
	}

	// create publications accordingly
	for _, pub := range publications {
		err := a.Store.CreatePublication(&pub)
		if err != nil {
			log.Printf("Error creating publication ID %d: %v", pub.ID, err)
			render.Render(w, r, ErrServer(err))
			return
		}
		log.Printf("Created publication ID %d successfully", pub.ID)
	}

	// respond	
	if err := render.Render(w, r, NewSyncResponse("Publications synchronized successfully")); err != nil {
	render.Render(w, r, ErrRender(err))
	return
	}
}

// SyncResponse represents the response after synchronization
type SyncResponse struct {
	Message string `json:"message"`
}

// NewSyncResponse creates a new SyncResponse
func NewSyncResponse(message string) *SyncResponse {
	return &SyncResponse{
		Message: message,
	}
}

// Render processes responses before marshalling.
func (pub *SyncResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
