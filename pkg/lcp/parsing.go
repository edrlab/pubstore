package lcp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

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

type LicenceLCPL struct {
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

func ParseLicenceLCPL(data []byte) (string, string, string, int, int, time.Time, time.Time, error) {
	var lcp LicenceLCPL
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

/*
{
	"id": "aea91a67-b1de-4761-97fa-9d2f038a20ba",
	"status": "ready",
	"updated": {
	  "license": "2023-06-22T12:07:51Z",
	  "status": "2023-06-22T12:07:51Z"
	},
	"message": "ready",
	"links": [
	  {
		"rel": "license",
		"href": "https://front-prod.edrlab.org/frontend/api/v1/licenses/aea91a67-b1de-4761-97fa-9d2f038a20ba",
		"type": "application/vnd.readium.lcp.license.v1.0+json"
	  },
	  {
		"rel": "register",
		"href": "https://front-prod.edrlab.org/lsdserver/licenses/aea91a67-b1de-4761-97fa-9d2f038a20ba/register{?id,name}",
		"type": "application/vnd.readium.license.status.v1.0+json",
		"templated": true
	  },
	  {
		"rel": "return",
		"href": "https://front-prod.edrlab.org/lsdserver/licenses/aea91a67-b1de-4761-97fa-9d2f038a20ba/return{?id,name}",
		"type": "application/vnd.readium.license.status.v1.0+json",
		"templated": true
	  },
	  {
		"rel": "renew",
		"href": "https://front-prod.edrlab.org/lsdserver/licenses/aea91a67-b1de-4761-97fa-9d2f038a20ba/renew{?end,id,name}",
		"type": "application/vnd.readium.license.status.v1.0+json",
		"templated": true
	  }
	],
	"potential_rights": {
	  "end": "2099-12-31T23:00:00Z"
	}
  }
*/

type Lsd struct {
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

func getLsdStatusDocument(url string) (Lsd, error) {

	response, err := http.Get(url)
	if err != nil {
		return Lsd{}, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return Lsd{}, err
	}

	var data Lsd
	err = json.Unmarshal(body, &data)
	if err != nil {
		return Lsd{}, err
	}

	return data, nil
}
