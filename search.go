package main

import (
	"net/http"
	"github.com/fluidmediaproductions/fluidmedia_crm/model"
	"database/sql"
)

type SearchContext struct {
	Results []*SearchResult
	Term string
}

type SearchResult struct {
	Obj interface{}
	Type struct {
		Name string
		URL string
	}
}

func handleSearch(m *model.Model, page *Page, user *model.User, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	term := r.Form.Get("term")

	results := make([]*SearchResult, 0)

	contacts, err := m.SearchContacts(term)
	if err != nil && err != sql.ErrNoRows {
		display500(w, err)
		return
	}
	for _, contact := range contacts {
		results = append(results, &SearchResult{
			Obj: contact,
			Type: struct {
				Name string
				URL string
			}{
				Name: "Contact",
				URL: "contacts",
			},
		})
	}

	organisations, err := m.SearchOrganisations(term)
	if err != nil && err != sql.ErrNoRows {
		display500(w, err)
		return
	}
	for _, organisation := range organisations {
		results = append(results, &SearchResult{
			Obj: organisation,
			Type: struct {
				Name string
				URL string
			}{
				Name: "Organisation",
				URL: "organisations",
			},
		})
	}

	displayWithContext(w, "search", page, user, &SearchContext{Term: term, Results: results})
}
