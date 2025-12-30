// Copyright 2025 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package lcp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/edrlab/pubstore/pkg/conf"
	"github.com/edrlab/pubstore/pkg/stor"
)

// GetPublications fetches publications from the LCP Server v2
// Publications are fetched page by page (20 per page), ordered by id DESC (newest first)
// If lastSync is zero, all publications are fetched
// Otherwise, only publications with CreatedAt >= lastSync are returned
func GetPublications(lcpsv conf.LCPServer, lastSync time.Time) ([]stor.Publication, error) {

	if lcpsv.Version != "v2" {
		return nil, errors.New("GetPublications is only supported for LCP Server v2")
	}

	log.Println("Requesting publications from LCP Server v2...")

	if lastSync.IsZero() {
		log.Println("Fetching all publications")
	} else {
		log.Printf("Fetching publications created since: %v", lastSync)
	}

	const perPage = 20
	var accumulatedPubs []stor.Publication
	page := 1
	reachedLastSync := false

	client := &http.Client{}

	for !reachedLastSync {
		// Build the URL with pagination parameters
		url := fmt.Sprintf("%s/publications?page=%d&per_page=%d", lcpsv.Url, page, perPage)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/json")
		req.SetBasicAuth(lcpsv.UserName, lcpsv.Password)

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send a request to the LCP Server: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			if resp.StatusCode >= http.StatusBadRequest && resp.StatusCode < http.StatusInternalServerError {
				return nil, fmt.Errorf("client error occurred. Status code: %d", resp.StatusCode)
			} else if resp.StatusCode == http.StatusInternalServerError {
				return nil, fmt.Errorf("server error occurred. Status code: %d", resp.StatusCode)
			} else {
				return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			}
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		var publications []stor.Publication
		err = json.Unmarshal(body, &publications)
		if err != nil {
			return nil, err
		}

		log.Printf("Fetched page %d with %d publications", page, len(publications))

		// If no publications in this page, we're done
		if len(publications) == 0 {
			break
		}

		// Filter publications by CreatedAt
		for _, pub := range publications {
			// If no last sync, include all
			if lastSync.IsZero() {
				accumulatedPubs = append(accumulatedPubs, pub)
			} else {
				// Compare CreatedAt with lastSync
				if pub.CreatedAt.Before(lastSync) {
					// We've reached publications older than lastSync, stop
					reachedLastSync = true
					break
				}
				// Publication is newer than or equal to lastSync, include it
				accumulatedPubs = append(accumulatedPubs, pub)
			}
		}

		// If we got fewer publications than requested, we've reached the end
		if len(publications) < perPage {
			break
		}

		page++
	}

	log.Printf("Total publications fetched: %d", len(accumulatedPubs))

	return accumulatedPubs, nil
}
