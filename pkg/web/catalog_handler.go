// Copyright 2025 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/edrlab/pubstore/pkg/view"
	"github.com/foolin/goview"
)

func (web *Web) catalogHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for the catalog handler
	// This function will handle the "/catalog" route

	format := r.URL.Query().Get("format")
	author := r.URL.Query().Get("author")
	publisher := r.URL.Query().Get("publisher")
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("pageSize")

	q := r.URL.Query().Get("q")
	query := r.URL.Query().Get("query")
	queryStr := query
	if len(query) == 0 {
		queryStr = q
	}

	pageInt, _ := strconv.Atoi(page)
	if pageInt < 1 || pageInt > 1000 {
		pageInt = 1
	}
	pageSizeInt, _ := strconv.Atoi(pageSize)
	if pageSizeInt < 1 || pageSizeInt > 1000 {
		pageSizeInt = web.Config.PageSize
	}

	var facet string = ""
	var value string = ""
	if len(queryStr) > 0 {
		facet = "search"
		value = queryStr
	} else if len(format) > 0 {
		facet = "format"
		value = format
	} else if len(author) > 0 {
		facet = "author"
		value = author
	} else if len(publisher) > 0 {
		facet = "publisher"
		value = publisher
	}

	facetsView := web.View.GetCatalogFacetsView()
	pubsView, count := web.View.GetCatalogPublicationsView(facet, value, pageInt, pageSizeInt)
	catalogView := view.GetCatalogView(pubsView, facetsView)

	var pageRange []string = make([]string, pageInt)
	for i := 0; i < pageInt; i++ {
		pageRange[i] = fmt.Sprintf("%d", i+1)
	}
	userStor := web.getUserByCookie(r)
	userName := ""
	if userStor != nil {
		userName = userStor.Name
	}

	goviewModel := goview.M{
		"pageTitle":           "pubstore - catalog",
		"userIsAuthenticated": web.userIsAuthenticated(r),
		"userName":            userName,
		"currentFacetType":    facet,
		"currentFacetValue":   value,
		"currentPageSize":     fmt.Sprintf("%d", pageSizeInt),
		"currentPage":         fmt.Sprintf("%d", pageInt),
		"pageRange":           pageRange,
		"publicationCount":    fmt.Sprintf("%d", count),
		"formats":             (*catalogView).Formats,
		"publications":        (*catalogView).Publications,
	}

	err := goview.Render(w, http.StatusOK, "catalog", goviewModel)
	if err != nil {
		fmt.Fprintf(w, "Render index error: %v!", err)
	}
}
