// Copyright 2025 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package web

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/edrlab/pubstore/pkg/lcp"
	"github.com/go-chi/render"
)

// syncpubHandler handles a request to synchronize publications from the LCP Server v2
func (web *Web) syncpubHandler(w http.ResponseWriter, r *http.Request) {

	if web.LCPServer.Version != "v2" {
		http.Error(w, url.QueryEscape("LCP Server V1 synchronization not supported"), 500)
		return
	}

	// get the last synchronization time
	lastSync, err := web.Store.GetLastSyncTime()
	if err != nil {
		log.Printf("Error getting last sync time: %v", err)
		http.Error(w, url.QueryEscape(err.Error()), 500)
		return
	}
	
	// Log the sync status
	if lastSync.IsZero() {
		log.Println("No previous synchronization found, fetching all publications")
	} else {
		log.Printf("Last synchronization time: %v", lastSync)
	}

	// get publications from the LCP Server v2
	publications, err := lcp.GetPublications(web.LCPServer, lastSync)
	if err != nil {
		log.Printf("Error fetching publications from LCP Server: %v", err)
		http.Error(w, url.QueryEscape(err.Error()), 500)
		return
	}

	log.Printf("Fetched %d publications from LCP Server", len(publications))

	// save a new synchronization time
	currentTime := time.Now()
	err = web.Store.SaveSyncTime(currentTime)
	if err != nil {
		log.Printf("Error saving sync time: %v", err)
		http.Error(w, url.QueryEscape(err.Error()), 500)
		return
	}

	// create publications accordingly
	for _, pub := range publications {
		err := web.Store.CreatePublication(&pub)
		if err != nil {
			log.Printf("Error creating publication ID %d: %v", pub.ID, err)
			http.Error(w, url.QueryEscape(err.Error()), 500)
			return
		}
		log.Printf("Created publication ID %d successfully", pub.ID)
	}

	// send back a success response
	render.JSON(w, r, map[string]string{"status": "synchronization completed successfully"})
}
