// Copyright 2025 European Digital Reading Lab. All rights reserved.
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
	"log"
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
	Copy          *int32     `json:"copy,omitempty"`
	Print         *int32     `json:"print,omitempty"`
	Profile       string     `json:"profile"`
	TextHint      string     `json:"text_hint"`
	PassHash      string     `json:"pass_hash"`
}

// LicenseRequestV1 represents the structure of a license request for the LCP Server version 1
type LicenseRequestV1 struct {
	Provider   string     `json:"provider"`
	User       User       `json:"user"`
	Encryption Encryption `json:"encryption"`
	Rights     Rights     `json:"rights,omitempty"`
}

type LCPLicense struct {
	Provider   string     `json:"provider"`
	ID         string     `json:"id"`
	Issued     string     `json:"issued"`
	Encryption Encryption `json:"encryption"`
	Links      []Link     `json:"links"`
	User       User       `json:"user"`
	Rights     Rights     `json:"rights"`
	Signature  Signature  `json:"signature"`
}

type Encryption struct {
	Profile    string     `json:"profile"`
	ContentKey ContentKey `json:"content_key"`
	UserKey    UserKey    `json:"user_key"`
}

type ContentKey struct {
	Algorithm      string `json:"algorithm"`
	EncryptedValue string `json:"encrypted_value"`
}

type UserKey struct {
	Algorithm string `json:"algorithm"`
	TextHint  string `json:"text_hint"`
	KeyCheck  string `json:"key_check"`
	HexValue  string `json:"hex_value"`
}

type Link struct {
	Rel    string `json:"rel"`
	Href   string `json:"href"`
	Type   string `json:"type"`
	Title  string `json:"title,omitempty"`
	Length int    `json:"length,omitempty"`
	Hash   string `json:"hash,omitempty"`
}

type User struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	Encrypted []string `json:"encrypted"`
}

type Rights struct {
	Print *int32     `json:"print,omitempty"`
	Copy  *int32     `json:"copy,omitempty"`
	Start *time.Time `json:"start,omitempty"`
	End   *time.Time `json:"end,omitempty"`
}

type Signature struct {
	Certificate string `json:"certificate"`
	Value       string `json:"value"`
	Algorithm   string `json:"algorithm"`
}

// HashPassphrase generates the hash of a passphrase
func HashPassphrase(passphrase string) string {

	hash := sha256.Sum256([]byte(passphrase))
	hashString := hex.EncodeToString(hash[:])

	return hashString
}

// v1Request prepares a license request in the form expected by the License Server V1
func v1Request(licenseReq LicenseRequest) LicenseRequestV1 {
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

	license := LicenseRequestV1{
		Provider:   "https://edrlab.org",
		User:       user,
		Encryption: encryption,
		Rights:     rights,
	}

	return license
}

// GenerateLicense sends a request to the License Server and returns a new license to the caller
func GenerateLicense(lcpsv conf.LCPServer, licenseReq LicenseRequest) (LCPLicense, []byte, error) {

	var license LCPLicense
	var url string
	var payload []byte
	var err error

	log.Println("Requesting the generation of a license...")

	// License Server V1
	if lcpsv.Version == "v1" {
		url = fmt.Sprintf(lcpsv.Url+"/contents/%s/license", licenseReq.PublicationID)
		v1req := v1Request(licenseReq)
		payload, err = json.Marshal(v1req)
		if err != nil {
			return license, nil, err
		}

		// License Server V2
	} else {
		url = lcpsv.Url + "/licenses"
		payload, err = json.Marshal(licenseReq)
		if err != nil {
			return license, nil, err
		}
	}

	return executeLicenseRequest(lcpsv, url, payload)
}

// GetFreshLicense sends a request to the License Server and returns the fresh license to the caller
func GetFreshLicense(lcpsv conf.LCPServer, transaction *stor.Transaction) (LCPLicense, []byte, error) {

	var license LCPLicense
	var url string
	var payload []byte
	var err error

	log.Println("Requesting a fresh license...")

	// License Server V1
	if lcpsv.Version == "v1" {
		url = lcpsv.Url + "/licenses/" + transaction.LicenseId

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
		licenseReq := LicenseRequestV1{
			User:       user,
			Encryption: encryption,
		}
		payload, err = json.Marshal(licenseReq)
		if err != nil {
			return license, nil, err
		}

		// License Server V2
	} else {
		url = lcpsv.Url + "/licenses/" + transaction.LicenseId

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
			return license, nil, err
		}
	}
	return executeLicenseRequest(lcpsv, url, payload)
}

// executeLicenseRequest sends a request to the License Server and returns the license
func executeLicenseRequest(lcpsv conf.LCPServer, url string, payload []byte) (LCPLicense, []byte, error) {

	var license LCPLicense
	var err error

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return license, nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(lcpsv.UserName, lcpsv.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return license, nil, errors.New("failed to send a license request to the License Server")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		if resp.StatusCode >= http.StatusBadRequest && resp.StatusCode < http.StatusInternalServerError {
			return license, nil, fmt.Errorf("client error occurred. Status code: %d", resp.StatusCode)
		} else if resp.StatusCode == http.StatusInternalServerError {
			return license, nil, fmt.Errorf("server error occurred. Status code: %d", resp.StatusCode)
		} else {
			return license, nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return license, nil, err
	}

	err = json.Unmarshal(body, &license)
	if err != nil {
		return license, nil, err
	}

	return license, body, nil
}
