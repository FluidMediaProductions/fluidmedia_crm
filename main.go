package main

import (
	"net/http"
	"github.com/fluidmediaproductions/fluidmedia_crm/model"
	"github.com/fluidmediaproductions/fluidmedia_crm/db"
	"flag"
	"log"
	"github.com/gorilla/mux"
	"strconv"
	"database/sql"
)

type Config struct {
	ListenSpec string

	Db db.Config
}

func parseFlags() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.ListenSpec, "listen", "localhost:8080", "HTTP listen spec")
	flag.StringVar(&cfg.Db.ConnectString, "db-connect", "user=postgres host=127.0.0.1 dbname=fluidmedia_crm", "DB Connect String")

	flag.Parse()
	return cfg
}


func handleIndex(m *model.Model, page *Page, w http.ResponseWriter, r *http.Request) {
	display(w, "index", page)
}

func handleContacts(m *model.Model, page *Page, w http.ResponseWriter, r *http.Request) {
	type ContactsContext struct {
		Contacts []*model.Contact
	}
	contacts, err := m.Contacts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	displayWithContext(w, "contacts", page, &ContactsContext{Contacts: contacts})
}

func handleContactsEdit(m *model.Model, page *Page, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	contact, err := m.Contact(id)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		display404(w)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		display500(w)
		return
	}
	if r.Method == "GET" {
		displayWithContext(w, "contacts-edit", page, contact)
	} else if r.Method == "POST" {
		r.ParseForm()
		newContact := &model.Contact{
			ID: id,
			Name: r.Form.Get("name"),
			Email: r.Form.Get("email"),
			Image: contact.Image,
		}
		m.SaveContact(newContact)
		http.Redirect(w, r, "/contacts", 302)
	}
}

func main() {
	cfg := parseFlags()

	pages = []*Page{
		{
			Title: "Home",
			Icon: "home",
			InMenu: true,
			Path: "/",
			Methods: []string{"GET"},
			Handler: handleIndex,
		},
		{
			Title: "Contacts",
			Icon: "person",
			InMenu: true,
			Path: "/contacts",
			Methods: []string{"GET"},
			Handler: handleContacts,
		},
		{
			Title: "Edit contact",
			InMenu: false,
			Path: "/contacts/{id:[0-9]+}",
			Methods: []string{"GET", "POST"},
			Handler: handleContactsEdit,
		},
	}

	dbInst, err := db.InitDb(cfg.Db)
	if err != nil {
		log.Fatalf("Can't initalize database: %v", err)
	}

	modelInst := model.New(dbInst)

	serveHttp(modelInst, pages, cfg.ListenSpec)
}
