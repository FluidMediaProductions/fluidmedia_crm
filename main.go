package main

import (
	"net/http"
	"github.com/fluidmediaproductions/fluidmedia_crm/model"
	"github.com/fluidmediaproductions/fluidmedia_crm/db"
	"flag"
	"log"
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


func handleIndex(model *model.Model, page *Page, w http.ResponseWriter, r *http.Request) {
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
	}

	dbInst, err := db.InitDb(cfg.Db)
	if err != nil {
		log.Fatalf("Can't initalize database: %v", err)
	}

	modelInst := model.New(dbInst)

	serveHttp(modelInst, pages, cfg.ListenSpec)
}
