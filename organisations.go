package main

import (
	"github.com/fluidmediaproductions/fluidmedia_crm/model"
	"net/http"
	"log"
	"fmt"
	"github.com/gorilla/mux"
	"strconv"
	"database/sql"
)

func handleOrganisations(m *model.Model, page *Page, w http.ResponseWriter, r *http.Request) {
	type OrganisationsContext struct {
		Organisations []*model.Organisation
	}
	organisations, err := m.Organisations()
	if err != nil {
		log.Printf("Error getting organisations: %v", err)
		display500(w)
		return
	}
	displayWithContext(w, "organisations", page, &OrganisationsContext{Organisations: organisations})
}

func handleOrganisationsEdit(m *model.Model, page *Page, w http.ResponseWriter, r *http.Request) {
	type OrganisationContext struct {
		Organisation *model.Organisation
	}
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	organisation, err := m.Organisation(id)
	if err == sql.ErrNoRows {
		display404(w)
		return
	} else if err != nil {
		log.Printf("Error getting organisation: %v", err)
		display500(w)
		return
	}
	if r.Method == "GET" {
		displayWithContext(w, "organisations-edit", page, &OrganisationContext{Organisation: organisation})
	} else if r.Method == "POST" {
		r.ParseForm()
		newOrganisation := &model.Organisation{
			ID: id,
			Name: r.Form.Get("name"),
			Image: organisation.Image,
			Email: r.Form.Get("email"),
			Phone: r.Form.Get("phone"),
			Website: r.Form.Get("website"),
			Twitter: r.Form.Get("twitter"),
			Address: r.Form.Get("address"),
			Description: r.Form.Get("desc"),
		}
		m.SaveOrganisation(newOrganisation)
		http.Redirect(w, r, "/organisations", 302)
	}
}

func handleOrganisationsNew(m *model.Model, page *Page, w http.ResponseWriter, r *http.Request) {
	organisationId, err := m.NewOrganisation()
	if err != nil {
		log.Printf("Error creating new organisation: %v", err)
		display500(w)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/organisations/%d", organisationId), 302)
}

func handleOrganisationsDelete(m *model.Model, page *Page, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	err := m.DeleteOrganisation(id)
	if err != nil {
		log.Printf("Error deleting organisation: %v", err)
		display500(w)
		return
	}
	http.Redirect(w, r, "/organisations", 302)
}
