// Copyright 2023 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package lcp

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type contentKey struct {
	Algorithm      string `json:"algorithm"`
	EncryptedValue string `json:"encrypted_value"`
}

type userKey struct {
	Algorithm string `json:"algorithm"`
	TextHint  string `json:"text_hint"`
	KeyCheck  string `json:"key_check"`
	HexValue  string `json:"hex_value"`
}

type encryption struct {
	Profile    string     `json:"profile"`
	ContentKey contentKey `json:"content_key"`
	UserKey    userKey    `json:"user_key"`
}

type link struct {
	Rel    string `json:"rel"`
	Href   string `json:"href"`
	Type   string `json:"type"`
	Title  string `json:"title,omitempty"`
	Length int    `json:"length,omitempty"`
	Hash   string `json:"hash,omitempty"`
}

type signature struct {
	Certificate string `json:"certificate"`
	Value       string `json:"value"`
	Algorithm   string `json:"algorithm"`
}

type LCPLicense struct {
	Provider   string     `json:"provider"`
	ID         string     `json:"id"`
	Issued     string     `json:"issued"`
	Encryption encryption `json:"encryption"`
	Links      []link     `json:"links"`
	User       user       `json:"user"`
	Rights     rights     `json:"rights"`
	Signature  signature  `json:"signature"`
}

type user struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	Encrypted []string `json:"encrypted"`
}

type rights struct {
	Print int       `json:"print"`
	Copy  int       `json:"copy"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func ParseLicense(data []byte) (string, string, string, int, int, time.Time, time.Time, error) {
	var lcp LCPLicense
	err := json.Unmarshal(data, &lcp)
	if err != nil {
		return "", "", "", 0, 0, time.Now(), time.Now(), err
	}

	// Extracting ID
	id := lcp.ID

	// Extracting link information
	var publicationLink link
	for _, l := range lcp.Links {
		if l.Rel == "publication" {
			publicationLink = l
			break
		}
	}

	// Extracting status link information
	var publicationStatus link
	for _, l := range lcp.Links {
		if l.Rel == "status" {
			publicationStatus = l
			break
		}
	}

	// Extracting publication link information
	// publicationType := publicationLink.Type
	publicationTitle := publicationLink.Title
	publicationStatusHref := publicationStatus.Href
	printRights := lcp.Rights.Print
	copyRights := lcp.Rights.Copy
	startDate := lcp.Rights.Start
	endDate := lcp.Rights.End
	// publicationLength := publicationLink.Length

	return id, publicationTitle, publicationStatusHref, printRights, copyRights, startDate, endDate, err
}

type StatusDoc struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Updated struct {
		License time.Time `json:"license"`
		Status  time.Time `json:"status"`
	} `json:"updated"`
	Message         string `json:"message"`
	Links           []link `json:"links"`
	PotentialRights struct {
		End time.Time `json:"end"`
	} `json:"potential_rights"`
}

// getStatusDocFromUrl gets a status document from a url
func getStatusDocFromUrl(url string) (StatusDoc, error) {

	response, err := http.Get(url)
	if err != nil {
		return StatusDoc{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return StatusDoc{}, err
	}

	var data StatusDoc
	err = json.Unmarshal(body, &data)
	if err != nil {
		return StatusDoc{}, err
	}

	return data, nil
}
