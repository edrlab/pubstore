package lcp

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

/**

`{
	  "provider": "http://www.imaginaryebookretailer.com",
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

type Licence struct {
	Provider   string     `json:"provider"`
	ID         string     `json:"id"`
	User       User       `json:"user"`
	Encryption Encryption `json:"encryption"`
	Rights     Rights     `json:"rights"`
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
	Print int       `json:"print"`
	Copy  int       `json:"copy"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

/**

	{
  "provider": "https://pubstore.edrlab.org",
  "id": "aea91a67-b1de-4761-97fa-9d2f038a20ba",
  "issued": "2023-06-22T12:07:51Z",
  "encryption": {
    "profile": "http://readium.org/lcp/profile-1.0",
    "content_key": {
      "algorithm": "http://www.w3.org/2001/04/xmlenc#aes256-cbc",
      "encrypted_value": "AqbaQUVuhI1VehuBsJ5uGjtDJLCOiuhhP/WvjKFi7BlBV/0mZNo+x5BX/3jAuMsmv+6+YT01pmJ7Pr+rQIbBDw=="
    },
    "user_key": {
      "algorithm": "http://www.w3.org/2001/04/xmlenc#sha256",
      "text_hint": "Hint",
      "key_check": "N0pZvdznClFwMfdsFimjYZehUum7tRtd0jVnuSH8rlasBmivCBQ/QlhIoZfNG9qb5UZQbYhM07g9g7yNYqUBdQ=="
    }
  },
  "links": [
    {
      "rel": "hint",
      "href": "https://front-prod.edrlab.org/frontend/static/hint.html",
      "type": "text/html"
    },
    {
      "rel": "publication",
      "href": "https://front-prod.edrlab.org/lcpserver/contents/6b4b7eb4-7630-4819-a813-e8276a83f78f",
      "type": "application/epub+zip",
      "title": "1342.encrypted.epub",
      "length": 24850367,
      "hash": "2b76a1fd895f05883db9da1e16d25f712af57d8d1a98ba8586cc645c3808f1d8"
    },
    {
      "rel": "status",
      "href": "https://front-prod.edrlab.org/lsdserver/licenses/aea91a67-b1de-4761-97fa-9d2f038a20ba/status",
      "type": "application/vnd.readium.license.status.v1.0+json"
    }
  ],
  "user": {
    "id": "7e1ac2b9-ed71-4180-bfb2-4a07bce46477",
    "email": "Qpz++L4gjQ0UolZcDnSg9vG5bKfX+rhLakGo+9JsvALrY1r+N3znlyLSPksgC9Wg",
    "encrypted": [
      "email"
    ]
  },
  "rights": {
    "print": 43,
    "copy": 12,
    "start": "2023-06-14T00:08:15Z",
    "end": "2099-12-31T23:00:00Z"
  },
  "signature": {
    "certificate": "MIIDKTCCAhGgAwIBAgIIEjZ5lJEfgf8wDQYJKoZIhvcNAQELBQAwQjETMBEGA1UEChMKZWRybGFiLm9yZzEXMBUGA1UECxMOZWRybGFiLm9yZyBMQ1AxEjAQBgNVBAMTCUVEUkxhYiBDQTAeFw0yMjA3MDgxNDMzMDFaFw0yNDA3MDcyMzU5NTlaME4xEzARBgNVBAoTCmVkcmxhYi5vcmcxJjAkBgNVBAsTHVJlYWRpdW0gTENQIExpY2Vuc2UgUHJvdmlkZXJzMQ8wDQYDVQQDEwZFRFJMYWIwgZswEAYHKoZIzj0CAQYFK4EEACMDgYYABAEU5BNfVLta4fz4MtmfHROMkLLThwyuKIKFeysg4cHjxBx0GAU+hGD3+rNTj7cDHa9FMQlE+sFNdbGfd2s4c3D8fAFErQI6QQ52MDuFSClaN0aWdVpjznc5V6Y6qvWTgh4P4V4gv40ot+QiVzwTevBWtsSbrw8nltCySIRGs66D7BEiU6OBnjCBmzAfBgNVHSMEGDAWgBTcXPyT5B+f7rC66lILK8pSXODJhzAdBgNVHQ4EFgQUCc2GhhgYvSOYvXkdYOxNqxs4FywwDgYDVR0PAQH/BAQDAgeAMAkGA1UdEwQCMAAwPgYDVR0fBDcwNTAzoDGgL4YtaHR0cDovL2NybC5lZHJsYWIudGVsZXNlYy5kZS9ybC9FRFJMYWJfQ0EuY3JsMA0GCSqGSIb3DQEBCwUAA4IBAQBxCl511aMdpIl56NKI0VW2tTM3FvhN717kNWsdr6Mj4xm2HXZ+BLfhGqFkm1iYwkM45o2unuVqe8zkIfEE24ghBd37aXrmS8IxY9t6gVFZKUGP/Q7NeA0EzUKru086mDDAuOgC05EAMlL6hgk+8IXw/BiD3hROAaop781UAkG3usU46n8w3meDqjjseLFfLTlGCU4njTGWZe3U8bOM0iz52LPcJGGT+fOPm2wGMdLL7aicxF166qWD05xC4UpdARwjopBGj7qkw6LVrM0E2mGpF0SyCyM4tdQH4PkHTtZ06vjipbzvE6TlFJP9/9M4HTUDDUevH4rPUip+de8wKvGF",
    "value": "AHaPyrqHzS+S2CGQrLZvyUVIsA8wabf+P4rnovY8PZHhtW3vhH/QdCXSN/r1ptF9W66NwMV3PFZySSFe56ihRhkvAK8Htk7zOY7UJrPqTj7KqmrTnJPVipJizeoW/p/1AKfiQ39tzybrj3oWO71whpLX1YCNBRO4GIOOUNXqRSdLkjYR",
    "algorithm": "http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha256"
  }
}
*/

func CreateLcpPassHash(passphrase string) string {

	hash := sha256.Sum256([]byte(passphrase))
	hashString := hex.EncodeToString(hash[:])

	return hashString
}

func generateLicence(provider, userUUID, userEmail, textHint, hexValue string, printRights, copyRights int, start, end time.Time) Licence {
	user := User{
		ID:        userUUID,
		Email:     userEmail,
		Encrypted: []string{"email"},
	}

	userKey := UserKey{
		TextHint: textHint,
		HexValue: hexValue,
	}

	encryption := Encryption{
		UserKey: userKey,
	}

	rights := Rights{
		Print: printRights,
		Copy:  copyRights,
		Start: start,
		End:   end,
	}

	licence := Licence{
		Provider:   provider,
		User:       user,
		Encryption: encryption,
		Rights:     rights,
	}

	return licence
}

func generateLicenceFromLcpServer(pubUUID, userUUID, userEmail, textHint, hexValue string, printRights, copyRights int, start, end time.Time) ([]byte, error) {

	provider := "https://pubstore.edrlab.org"
	licence := generateLicence(provider, userUUID, userEmail, textHint, hexValue, printRights, copyRights, start, end)

	url := fmt.Sprintf("https://front-prod.edrlab.org/lcpserver/contents/%s/license", pubUUID)
	username := "adm_username"
	password := "adm_password"

	payload, err := json.Marshal(licence)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("License created successfully.")
	} else if resp.StatusCode >= http.StatusBadRequest && resp.StatusCode < http.StatusInternalServerError {
		return nil, fmt.Errorf("client error occurred. Status code: %d", resp.StatusCode)
	} else if resp.StatusCode == http.StatusInternalServerError {
		return nil, fmt.Errorf("server error occurred. Status code: %d", resp.StatusCode)
	} else {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func LicenceBuy(pubUUID, userUUID, userEmail, textHint, hexValue string, printRights, copyRights int) ([]byte, error) {
	start := time.Now()
	end := time.Now().AddDate(100, 0, 0)
	return generateLicenceFromLcpServer(pubUUID, userUUID, userEmail, textHint, hexValue, printRights, copyRights, start, end)
}

func LicenceLoan(pubUUID, userUUID, userEmail, textHint, hexValue string, printRights, copyRights int, start, end time.Time) ([]byte, error) {

	return generateLicenceFromLcpServer(pubUUID, userUUID, userEmail, textHint, hexValue, printRights, copyRights, start, end)
}

func GenerateFreshLicenceFromLcpServer(licenceId, email, textHint, hexValue string) ([]byte, error) {

	provider := "https://pubstore.edrlab.org"
	user := User{
		Email:     email,
		Encrypted: []string{"email"},
	}

	userKey := UserKey{
		TextHint: textHint,
		HexValue: hexValue,
	}

	encryption := Encryption{
		UserKey: userKey,
	}

	licence := Licence{
		Provider:   provider,
		User:       user,
		Encryption: encryption,
	}

	payload, err := json.Marshal(licence)
	if err != nil {
		return nil, err
	}

	url := "https://front-prod.edrlab.org/lcpserver/licenses/" + licenceId
	username := "adm_username"
	password := "adm_password"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("License created successfully.")
	} else if resp.StatusCode >= http.StatusBadRequest && resp.StatusCode < http.StatusInternalServerError {
		return nil, fmt.Errorf("client error occurred. Status code: %d", resp.StatusCode)
	} else if resp.StatusCode == http.StatusInternalServerError {
		return nil, fmt.Errorf("server error occurred. Status code: %d", resp.StatusCode)
	} else {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

type LsdStatus struct {
	StatusMessage      string
	StatusCode         string
	EndPotentialRights time.Time
	PrintRights        int
	CopyRights         int
	StartDate          time.Time
	EndDate            time.Time
}

func GetLsdStatus(licenceId, email, textHint, hexValue string) (*LsdStatus, error) {

	licenceBytes, err := GenerateFreshLicenceFromLcpServer(licenceId, email, textHint, hexValue)
	if err != nil {
		return nil, err

	}

	_, _, publicationStatusHref, printRights, copyRights, startDate, endDate, err := ParseLicenceLCPL(licenceBytes)
	if err != nil {
		return nil, err
	}

	// make a request on publicationStatusHref
	lsd, err := getLsdStatusDocument(publicationStatusHref)
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
		PrintRights:        printRights,
		CopyRights:         copyRights,
		StartDate:          startDate,
		EndDate:            endDate,
	}, nil
}
