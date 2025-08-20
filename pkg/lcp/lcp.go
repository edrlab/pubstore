// Copyright 2023 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package lcp

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/edrlab/pubstore/pkg/conf"
	"github.com/edrlab/pubstore/pkg/stor"
)

// LicenseRequest contains every parameter required for requesting a
// license to the License Server v2
type LicenseRequest struct {
	PublicationID string     `json:"publication_id"`
	UserID        string     `json:"user_id,omitempty"`
	UserName      string     `json:"user_name,omitempty"`
	UserEmail     string     `json:"user_email,omitempty"`
	UserEncrypted []string   `json:"user_encrypted,omitempty"`
	Start         *time.Time `json:"start,omitempty"`
	End           *time.Time `json:"end,omitempty"`
	Copy          *int       `json:"copy,omitempty"`
	Print         *int       `json:"print,omitempty"`
	Profile       string     `json:"profile"`
	TextHint      string     `json:"text_hint"`
	PassHash      string     `json:"pass_hash"`
}

/* Example of License Request V1
`{
	  "user": {
	    "id": "d9f298a7-7f34-49e7-8aae-4378ecb1d597",
	    "email": "user@mymail.com",
	    "encrypted": ["email"]
	  },
	  "encryption": {
	    "user_key": {
	      "text_hint": "The title of the first book you ever read",
	      "hex_value": "4981AA0A50D563040519E9032B5D74367B1D129E239A1BA82667A57333866494"
	    }
	  },
	  "rights": {
	    "print": 10,
	    "copy": 2048,
	    "start": "2023-06-14T01:08:15+01:00",
	    "end": "2024-11-25T01:08:15+01:00"
	  }
	}`
*/

type LicenceRequestV1 struct {
	Provider   string     `json:"provider"`
	User       User       `json:"user"`
	Encryption Encryption `json:"encryption"`
	Rights     Rights     `json:"rights,omitempty"`
}

type User struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	Encrypted []string `json:"encrypted"`
}

type Encryption struct {
	UserKey UserKey `json:"user_key"`
}

type UserKey struct {
	TextHint string `json:"text_hint"`
	HexValue string `json:"hex_value"`
}

type Rights struct {
	Print *int       `json:"print,omitempty"`
	Copy  *int       `json:"copy,omitempty"`
	Start *time.Time `json:"start,omitempty"`
	End   *time.Time `json:"end,omitempty"`
}

// HashPassphrase generates the hash of a passphrase
func HashPassphrase(passphrase string) string {

	hash := sha256.Sum256([]byte(passphrase))
	hashString := hex.EncodeToString(hash[:])

	return hashString
}

// v1Request prepares a license request in the form expected by the License Server V1
func v1Request(licenseReq LicenseRequest) LicenceRequestV1 {
	user := User{
		ID:        licenseReq.UserID,
		Email:     licenseReq.UserEmail,
		Encrypted: []string{"email"},
	}

	userKey := UserKey{
		TextHint: licenseReq.TextHint,
		HexValue: licenseReq.PassHash,
	}

	encryption := Encryption{
		UserKey: userKey,
	}

	rights := Rights{
		Print: licenseReq.Print,
		Copy:  licenseReq.Copy,
		Start: licenseReq.Start,
		End:   licenseReq.End,
	}

	licence := LicenceRequestV1{
		Provider:   "https://edrlab.org",
		User:       user,
		Encryption: encryption,
		Rights:     rights,
	}

	return licence
}

// GenerateLicense sends a request to the License Server and returns a new license to the caller
func GenerateLicense(lcpsv conf.LCPServerAccess, licenseReq LicenseRequest) ([]byte, error) {

	var url string
	var payload []byte
	var err error

	// License Server V1
	if lcpsv.Version == "v1" {
		url = fmt.Sprintf(lcpsv.Url+"/contents/%s/license", licenseReq.PublicationID)
		v1req := v1Request(licenseReq)
		payload, err = json.Marshal(v1req)
		if err != nil {
			return nil, err
		}

		// License Server V2
	} else {
		url = lcpsv.Url + "/licenses"
		payload, err = json.Marshal(licenseReq)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(lcpsv.UserName, lcpsv.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, errors.New("failed to send a license request to the License Server")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("License created successfully.")
	} else if resp.StatusCode >= http.StatusBadRequest && resp.StatusCode < http.StatusInternalServerError {
		return nil, fmt.Errorf("a client error occurred. Status code: %d", resp.StatusCode)
	} else if resp.StatusCode == http.StatusInternalServerError {
		return nil, fmt.Errorf("a server error occurred. Status code: %d", resp.StatusCode)
	} else {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// GetFreshLicense sends a request to the License Server and returns the fresh license to the caller
func GetFreshLicense(lcpsv conf.LCPServerAccess, transaction *stor.Transaction) ([]byte, error) {

	var url string
	var payload []byte
	var err error

	// License Server V1
	if lcpsv.Version == "v1" {
		url = lcpsv.Url + "/licenses/" + transaction.LicenceId

		user := User{
			Email:     transaction.User.Email,
			Encrypted: []string{"email"},
		}
		userKey := UserKey{
			TextHint: transaction.User.TextHint,
			HexValue: transaction.User.HPassphrase,
		}
		encryption := Encryption{
			UserKey: userKey,
		}
		licence := LicenceRequestV1{
			User:       user,
			Encryption: encryption,
		}
		payload, err = json.Marshal(licence)
		if err != nil {
			return nil, err
		}

		// License Server V2
	} else {
		url = lcpsv.Url + "/licenses/" + transaction.LicenceId

		licenseReq := LicenseRequest{
			PublicationID: transaction.Publication.UUID,
			UserID:        transaction.User.UUID,
			UserEmail:     transaction.User.Email,
			UserEncrypted: []string{"email"},
			TextHint:      transaction.User.TextHint,
			PassHash:      transaction.User.HPassphrase,
		}
		payload, err = json.Marshal(licenseReq)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(lcpsv.UserName, lcpsv.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Fresh License successfully fetched.")
	} else if resp.StatusCode >= http.StatusBadRequest && resp.StatusCode < http.StatusInternalServerError {
		return nil, fmt.Errorf("client error occurred. Status code: %d", resp.StatusCode)
	} else if resp.StatusCode == http.StatusInternalServerError {
		return nil, fmt.Errorf("server error occurred. Status code: %d", resp.StatusCode)
	} else {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

type LsdStatus struct {
	StatusMessage      string
	StatusCode         string
	EndPotentialRights time.Time
	PrintLimit         int
	CopyLimit          int
	StartDate          time.Time
	EndDate            time.Time
}

// GetStatusDocument sends a request to the License Server and returns a status document to the caller
func GetStatusDocument(lcpsv conf.LCPServerAccess, transaction *stor.Transaction) (*LsdStatus, error) {

	// TODO: avoid fetching the fresh license first, just to get the url to the status document.
	licenceBytes, err := GetFreshLicense(lcpsv, transaction)
	if err != nil {
		return nil, err
	}

	_, _, publicationStatusHref, printRights, copyRights, startDate, endDate, err := ParseLicense(licenceBytes)
	if err != nil {
		return nil, err
	}

	// make a request on publicationStatusHref
	lsd, err := getStatusDocFromUrl(publicationStatusHref)
	if err != nil {
		return nil, err
	}

	statusMessage := lsd.Message
	endPotentialRights := lsd.PotentialRights.End
	statusCode := lsd.Status

	return &LsdStatus{
		StatusMessage:      statusMessage,
		StatusCode:         statusCode,
		EndPotentialRights: endPotentialRights,
		PrintLimit:         printRights,
		CopyLimit:          copyRights,
		StartDate:          startDate,
		EndDate:            endDate,
	}, nil
}
