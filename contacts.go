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

type ContactsContext struct {
	ContactStates   map[int][2]string
	ContactedStates map[int]string
	Contact         *model.Contact
	Contacts        []*model.Contact
	Organisations   []*model.Organisation
}

func handleContacts(m *model.Model, page *Page, user *model.User, w http.ResponseWriter, r *http.Request) {
	contacts, err := m.Contacts()
	if err != nil {
		log.Printf("Error getting contacts: %v", err)
		display500(w)
		return
	}
	organisations, err := m.Organisations()
	if err != nil {
		log.Printf("Error getting organisationss: %v", err)
		display500(w)
		return
	}
	displayWithContext(w, "contacts", page, user, &ContactsContext{Contacts: contacts, ContactStates: m.ContactStates(),
		ContactedStates: m.ContactedStates(), Organisations: organisations})
}

func handleContactsEdit(m *model.Model, page *Page, user *model.User, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	contact, err := m.Contact(id)
	if err == sql.ErrNoRows {
		display404(w)
		return
	} else if err != nil {
		log.Printf("Error getting contact: %v", err)
		display500(w)
		return
	}
	if r.Method == "GET" {
		organisations, err := m.Organisations()
		if err != nil {
			log.Printf("Error getting organisationss: %v", err)
			display500(w)
			return
		}
		displayWithContext(w, "contacts-edit", page, user, &ContactsContext{Contact: contact,
			ContactStates: m.ContactStates(), ContactedStates: m.ContactedStates(), Organisations: organisations})
	} else if r.Method == "POST" {
		r.ParseForm()
		state, err := strconv.Atoi(r.Form.Get("state"))
		if err != nil {
			state = 1
		}
		contactedState, err := strconv.Atoi(r.Form.Get("contacted_state"))
		if err != nil {
			contactedState = 1
		}
		organisation, err := strconv.Atoi(r.Form.Get("organisation"))
		if err != nil {
			organisation = 0
		}
		newContact := &model.Contact{
			ID:             id,
			Name:           r.Form.Get("name"),
			State:          state,
			ContactedState: contactedState,
			Image:          contact.Image,
			Email:          r.Form.Get("email"),
			Phone:          r.Form.Get("phone"),
			Mobile:         r.Form.Get("mobile"),
			Website:        r.Form.Get("website"),
			Twitter:        r.Form.Get("twitter"),
			Youtube:        r.Form.Get("youtube"),
			Instagram:      r.Form.Get("instagram"),
			Facebook:       r.Form.Get("facebook"),
			Address:        r.Form.Get("address"),
			Description:    r.Form.Get("desc"),
			OrganisationId: organisation,
		}
		err = m.SaveContact(newContact)
		if err != nil {
			log.Printf("Error updating contact: %v", err)
			display500(w)
			return
		}
		http.Redirect(w, r, "/contacts", 302)
	}
}

func handleContactsNew(m *model.Model, page *Page, user *model.User, w http.ResponseWriter, r *http.Request) {
	contactId, err := m.NewContact()
	if err != nil {
		log.Printf("Error creating new contact: %v", err)
		display500(w)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/contacts/%d", contactId), 302)
}

func handleContactsDelete(m *model.Model, page *Page, user *model.User, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	err := m.DeleteContact(id)
	if err != nil {
		log.Printf("Error deleting contact: %v", err)
		display500(w)
		return
	}
	http.Redirect(w, r, "/contacts", 302)
}
