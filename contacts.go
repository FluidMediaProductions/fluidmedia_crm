package main

import (
	"github.com/fluidmediaproductions/fluidmedia_crm/model"
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"strconv"
	"database/sql"
	"fmt"
)

func handleContacts(m *model.Model, page *Page, w http.ResponseWriter, r *http.Request) {
	type ContactsContext struct {
		ContactStates map[int][2]string
		Contacts []*model.Contact
	}
	contacts, err := m.Contacts()
	if err != nil {
		log.Printf("Error getting contacts: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		display500(w)
		return
	}
	displayWithContext(w, "contacts", page, &ContactsContext{Contacts: contacts, ContactStates: m.ContactStates()})
}

func handleContactsEdit(m *model.Model, page *Page, w http.ResponseWriter, r *http.Request) {
	type ContactContext struct {
		ContactStates map[int][2]string
		Contact *model.Contact
	}
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	contact, err := m.Contact(id)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		display404(w)
		return
	} else if err != nil {
		log.Printf("Error getting contact: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		display500(w)
		return
	}
	if r.Method == "GET" {
		displayWithContext(w, "contacts-edit", page, &ContactContext{Contact: contact, ContactStates: m.ContactStates()})
	} else if r.Method == "POST" {
		r.ParseForm()
		state, _ := strconv.Atoi(r.Form.Get("state"))
		newContact := &model.Contact{
			ID: id,
			Name: r.Form.Get("name"),
			State: state,
			Image: contact.Image,
			Email: r.Form.Get("email"),
			Phone: r.Form.Get("phone"),
			Mobile: r.Form.Get("mobile"),
			Website: r.Form.Get("website"),
			Twitter: r.Form.Get("twitter"),
			Address: r.Form.Get("address"),
			Description: r.Form.Get("desc"),
		}
		m.SaveContact(newContact)
		http.Redirect(w, r, "/contacts", 302)
	}
}

func handleContactsNew(m *model.Model, page *Page, w http.ResponseWriter, r *http.Request) {
	contactId, err := m.NewContact()
	if err != nil {
		log.Printf("Error creating new contact: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		display500(w)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/contacts/%d", contactId), 302)
}

func handleContactsDelete(m *model.Model, page *Page, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	err := m.DeleteContact(id)
	if err != nil {
		log.Printf("Error deleting contact: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		display500(w)
		return
	}
	http.Redirect(w, r, "/contacts", 302)
}
