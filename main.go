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

func handleIndex(m *model.Model, page *Page, w http.ResponseWriter, r *http.Request) {
	display(w, "index", page)
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
		{
			InMenu: false,
			Path: "/contacts/new",
			Methods: []string{"GET"},
			Handler: handleContactsNew,
		},
		{
			InMenu: false,
			Path: "/contacts/del/{id:[0-9]+}",
			Methods: []string{"GET"},
			Handler: handleContactsDelete,
		},


		{
			Title: "Organisations",
			Icon: "business",
			InMenu: true,
			Path: "/organisations",
			Methods: []string{"GET"},
			Handler: handleOrganisations,
		},
		{
			Title: "Edit organisation",
			InMenu: false,
			Path: "/organisations/{id:[0-9]+}",
			Methods: []string{"GET", "POST"},
			Handler: handleOrganisationsEdit,
		},
		{
			InMenu: false,
			Path: "/organisations/new",
			Methods: []string{"GET"},
			Handler: handleOrganisationsNew,
		},
		{
			InMenu: false,
			Path: "/organisations/del/{id:[0-9]+}",
			Methods: []string{"GET"},
			Handler: handleOrganisationsDelete,
		},
	}

	dbInst, err := db.InitDb(cfg.Db)
	if err != nil {
		log.Fatalf("Can't initalize database: %v", err)
	}

	modelInst := model.New(dbInst)

	serveHttp(modelInst, pages, cfg.ListenSpec)
}
