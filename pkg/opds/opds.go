// Copyright 2023 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package opds

import (
	"time"

	"github.com/edrlab/pubstore/pkg/conf"
	"github.com/edrlab/pubstore/pkg/stor"
)

type MetadataFeed struct {
	Title string `json:"title"`
}

// TODO: an OPDS Publication supports x languages as an array -> update the model and mapping
type Metadata struct {
	Type       string `json:"@type"`
	Title      string `json:"title"`
	Author     string `json:"author,omitempty"`
	Identifier string `json:"identifier,omitempty"`
	Language   string `json:"language,omitempty"`
	Published  string `json:"published,omitempty"`
}

type Link struct {
	Rel        string      `json:"rel,omitempty"`
	Href       string      `json:"href,omitempty"`
	Type       string      `json:"type,omitempty"`
	Child      []Link      `json:"child,omitempty"`
	Properties *Properties `json:"properties,omitempty"`
}

type Publication struct {
	Metadata Metadata `json:"metadata"`
	Links    []Link   `json:"links,omitempty"`
	Images   []Image  `json:"images,omitempty"`
}

type Image struct {
	Href   string `json:"href"`
	Type   string `json:"type,omitempty"`
	Height int    `json:"height,omitempty"`
	Width  int    `json:"width,omitempty"`
}

type Availability struct {
	Status    string     `json:"state,omitempty"`
	StartDate *time.Time `json:"since,omitempty"`
	EndDate   *time.Time `json:"until,omitempty"`
}

type Properties struct {
	Availability        *Availability `json:"availability,omitempty"`
	IndirectAcquisition []Link        `json:"indirectAcquisition,omitempty"`
	LcpHashedPassphrase string        `json:"lcp_hashed_passphrase,omitempty"`
}

type Root struct {
	Metadata     MetadataFeed  `json:"metadata"`
	Links        []Link        `json:"links"`
	Publications []Publication `json:"publications"`
}

type Opds struct {
	*conf.Config
	*stor.Store
}

// PublicBaseURL is set with the base URL of the server.
// Creating this avoids passing the whole config to many inner functions of this module.
var publicBaseUrl string

// Init initializes the module
func Init(c *conf.Config, s *stor.Store) Opds {

	publicBaseUrl = c.PublicBaseUrl

	return Opds{
		Config: c,
		Store:  s,
	}
}
