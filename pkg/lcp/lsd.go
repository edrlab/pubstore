// Copyright 2025 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package lcp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/edrlab/pubstore/pkg/conf"
	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/jtacoma/uritemplates"
)

type LsdStatus struct {
	StatusMessage string
	StatusCode    string
	Links         map[string]string
	Print         int32
	Copy          int32
	Start         *time.Time
	End           *time.Time
	MaxEnd        *time.Time
}

type StatusDoc struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Updated struct {
		License time.Time `json:"license"`
		Status  time.Time `json:"status"`
	} `json:"updated"`
	Message         string `json:"message"`
	Links           []Link `json:"links"`
	PotentialRights struct {
		End *time.Time `json:"end,omitempty"`
	} `json:"potential_rights"`
}

// See Problem Details for HTTP APIs, rfc 7807 : https://tools.ietf.org/html/rfc7807
type ErrResponse struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

// GetStatusDocument sends a request to the License Server and returns a status document to the caller
func GetStatusDocument(lcpsv conf.LCPServer, transaction *stor.Transaction) (*LsdStatus, error) {

	// fetch the status document
	statusDoc, err := getStatusDocFromUrl(transaction.StatusDocLink)
	if err != nil {
		return nil, err
	}

	// find the links
	var registerLink, renewLink, returnLink string
	for _, l := range statusDoc.Links {
		switch l.Rel {
		case "register":
			registerLink = l.Href
		case "renew":
			renewLink = l.Href
		case "return":
			returnLink = l.Href
		}
	}

	// set status document links
	links := map[string]string{
		"register": registerLink,
		"renew":    renewLink,
		"return":   returnLink,
	}

	// if the license has been updated, get the fresh license
	if statusDoc.Updated.License.After(transaction.LicenseUpdated) {
		license, _, err := GetFreshLicense(lcpsv, transaction)
		if err != nil {
			log.Printf("failed to get a fresh license: %s", err)
			return nil, fmt.Errorf("failed to get a fresh license: %w", err)
		}
		noLimit := int32(-1) // -1 stored for no print/copy limits
		if license.Rights.Copy == nil {
			license.Rights.Copy = &noLimit
		}
		if license.Rights.Print == nil {
			license.Rights.Print = &noLimit
		}
		// update the transaction with fresh license rights
		transaction.Print = *license.Rights.Print
		transaction.Copy = *license.Rights.Copy
		transaction.Start = license.Rights.Start
		transaction.End = license.Rights.End
	}

	// update the transaction
	transaction.Status = statusDoc.Status
	transaction.LicenseUpdated = statusDoc.Updated.License

	return &LsdStatus{
		StatusMessage: statusDoc.Message,
		StatusCode:    statusDoc.Status,
		MaxEnd:        statusDoc.PotentialRights.End,
		Links:         links,
		Print:         transaction.Print,
		Copy:          transaction.Copy,
		Start:         transaction.Start,
		End:           transaction.End,
	}, nil
}

// getStatusDocFromUrl gets a status document from a url
func getStatusDocFromUrl(url string) (StatusDoc, error) {

	log.Println("Get a status document")

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

// ExecuteAction executes register, renew, return and revoke actions
func ExecuteAction(lcpsv conf.LCPServer, transaction *stor.Transaction, action string) error {

	// select the http method
	var method string
	switch action {
	case "register":
		method = http.MethodPost
	case "renew":
		method = http.MethodPut
	case "return":
		method = http.MethodPut
	case "revoke":
		method = http.MethodPut
	default:
		return errors.New("unknown action: " + action)
	}

	var actionLink string
	var body io.Reader
	if action != "revoke" {
		// fetch the status document
		// TODO: should we store the register, renew, return links in the transaction instead ?
		statusDoc, err := GetStatusDocument(lcpsv, transaction)
		if err != nil {
			return err
		}

		// get the action link from the status document
		actionLink = statusDoc.Links[action]
		if actionLink == "" {
			return errors.New("no link found")
		}
		// replace url parameters by actual values
		actionLink = expandURI(actionLink)

		// special case for the revoke action
	} else {
		if lcpsv.Version == "v2" {
			actionLink = lcpsv.Url + "/revoke/" + transaction.LicenseId
			// in case a v1 LCP Server is integrated
		} else {
			actionLink = lcpsv.Url + "/licenses/" + transaction.LicenseId + "/status"
			method = http.MethodPatch
			s := map[string]interface{}{
				"status": "revoked",
			}
			jsonBytes, err := json.Marshal(s)
			if err != nil {
				return fmt.Errorf("failed to marshal payload: %w", err)
			}
			body = bytes.NewBuffer(jsonBytes)
		}
	}

	// execute the action
	req, _ := http.NewRequest(method, actionLink, body)
	// revoke needs authentication
	if action == "revoke" {
		req.SetBasicAuth(lcpsv.UserName, lcpsv.Password)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return errors.New("User unauthorized")
	} else if resp.StatusCode != http.StatusOK {
		// parse the error response
		var apiError ErrResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
			log.Printf("failed to parse error response: %s", err)
			return fmt.Errorf("failed to parse error response: %w", err)
		}
		return fmt.Errorf("%s: %s", apiError.Title, apiError.Detail)
	} else {
		// parse the status document which is sent as a response
		var statusDoc StatusDoc
		if err := json.NewDecoder(resp.Body).Decode(&statusDoc); err != nil {
			log.Printf("failed to parse status document: %s", err)
			return fmt.Errorf("failed to parse status document: %w", err)
		}

		// if the license has been updated, get the fresh license
		if statusDoc.Updated.License.After(transaction.LicenseUpdated) {
			license, _, err := GetFreshLicense(lcpsv, transaction)
			if err != nil {
				log.Printf("failed to get a fresh license: %s", err)
				return fmt.Errorf("failed to get a fresh license: %w", err)
			}
			noLimit := int32(-1) // -1 stored for no print/copy limits
			if license.Rights.Copy == nil {
				license.Rights.Copy = &noLimit
			}
			if license.Rights.Print == nil {
				license.Rights.Print = &noLimit
			}
			// update the transaction with fresh license rights
			transaction.Print = *license.Rights.Print
			transaction.Copy = *license.Rights.Copy
			transaction.Start = license.Rights.Start
			transaction.End = license.Rights.End
		}

		// update the transaction with the new status and the latest license update time
		transaction.Status = statusDoc.Status
		transaction.LicenseUpdated = statusDoc.Updated.License

	}
	return nil
}

// expandURI is a helper function to expand URI templates
func expandURI(url string) string {
	template, _ := uritemplates.Parse(url)
	values := make(map[string]interface{})
	// TODO: deal with "end"?
	values["id"] = "EDRLab-PubStore-ID"
	values["name"] = "PubStore"
	expanded, err := template.Expand(values)
	if err != nil {
		log.Printf("failed to expand the link: %s", template)
		expanded = url // fallback
	}
	return expanded
}
