// Copyright 2025 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/edrlab/pubstore/pkg/view"
	"github.com/foolin/goview"
)

func (web *Web) bookshelfHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for the bookshelf handler
	// This function will handle the "/user/bookshelf" route

	user := web.getUserByCookie(r)
	if user == nil {
		fmt.Fprintf(w, "bookshelf error")
		return
	}
	transactions, err := web.Store.FindTransactionsByUser(user.ID)
	if err != nil {
		fmt.Fprintf(w, "bookshelf error")
		return
	}

	formatMap := make(map[string]bool)
	var allTransactionViews []*view.TransactionView

	for _, transactionStor := range *transactions {
		tv := web.View.GetTransactionViewFromTransactionStor(&transactionStor)
		allTransactionViews = append(allTransactionViews, tv)

		if tv.PublicationFormat != "" {
			formatMap[tv.PublicationFormat] = true
		}
	}
	uniqueFormats := make([]string, 0, len(formatMap))
	for format := range formatMap {
		uniqueFormats = append(uniqueFormats, format)
	}

	ready := r.URL.Query().Get("ready") == "on"
	active := r.URL.Query().Get("active") == "on"
	expired := r.URL.Query().Get("expired") == "on"
	returned := r.URL.Query().Get("returned") == "on"
	revoked := r.URL.Query().Get("revoked") == "on"
	cancelled := r.URL.Query().Get("cancelled") == "on"
	buy := r.URL.Query().Get("buy") == "on"
	loan := r.URL.Query().Get("loan") == "on"

	formatParams := r.URL.Query()["format"]
	formatFilters := make(map[string]bool)
	for _, format := range formatParams {
		formatFilters[format] = true
	}

	var filteredTransactions []*view.TransactionView
	for _, tv := range allTransactionViews {
		statusMatch := true
		if ready || active || expired || returned || revoked || cancelled {
			statusMatch = false
			switch tv.LicenseStatus {
			case "ready":
				statusMatch = ready
			case "active":
				statusMatch = active
			case "expired":
				statusMatch = expired
			case "returned":
				statusMatch = returned
			case "revoked":
				statusMatch = revoked
			case "cancelled":
				statusMatch = cancelled
			}
		}

		transactionTypesMatch := true
		if buy || loan {
			transactionTypesMatch = false
			if tv.PublicationEndDate != "unknown" {
				transactionTypesMatch = loan
			} else {
				transactionTypesMatch = buy
			}
		}

		formatMatch := true
		if len(formatFilters) > 0 {
			formatMatch = formatFilters[tv.PublicationFormat]
		}

		if statusMatch && formatMatch && transactionTypesMatch {
			filteredTransactions = append(filteredTransactions, tv)
		}
	}

	activeFilters := []string{}

	if ready {
		activeFilters = append(activeFilters, "Ready")
	}
	if active {
		activeFilters = append(activeFilters, "Active")
	}
	if expired {
		activeFilters = append(activeFilters, "Expired")
	}
	if returned {
		activeFilters = append(activeFilters, "Returned")
	}
	if revoked {
		activeFilters = append(activeFilters, "Revoked")
	}
	if cancelled {
		activeFilters = append(activeFilters, "Cancelled")
	}
	if buy {
		activeFilters = append(activeFilters, "Buy")
	}
	if loan {
		activeFilters = append(activeFilters, "Loan")
	}
	for format, ok := range formatFilters {
		if ok {
			// On peut mettre une majuscule
			activeFilters = append(activeFilters, strings.Title(format))
		}
	}

	// On transforme en string "Buy, Loan"
	activeFiltersStr := strings.Join(activeFilters, ", ")

	goviewModel := goview.M{
		"pageTitle":             "pubstore - bookshelf",
		"userIsAuthenticated":   true,
		"userName":              user.Name,
		"uniqueFormats":         uniqueFormats,
		"transactions":          filteredTransactions,
		"bookshelfFilterActive": bool(len(activeFilters) > 0),
		"transactionsCount":     len(filteredTransactions),
		"filters": map[string]bool{
			"ready":   ready,
			"active":  active,
			"expired": expired,
			"returned": returned,
			"revoked": revoked,
			"cancelled": cancelled,
		},
		"formatFilters": formatFilters,
		"transactionTypes": map[string]bool{
			"buy":  buy,
			"loan": loan,
		},
		"activeFiltersStr": activeFiltersStr,
	}

	err = goview.Render(w, http.StatusOK, "bookshelf", goviewModel)
	if err != nil {
		fmt.Fprintf(w, "Render index error: %v!", err)
	}
}
