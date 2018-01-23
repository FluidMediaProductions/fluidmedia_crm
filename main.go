package main

import (
	"net/http"
	"github.com/fluidmediaproductions/fluidmedia_crm/model"
	"github.com/fluidmediaproductions/fluidmedia_crm/db"
	"flag"
	"log"
	"os"
	"fmt"
)

type Config struct {
	ListenSpec string

	Db db.Config
}

func parseFlags() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.ListenSpec, "listen", ":8080", "HTTP listen spec")
	flag.StringVar(&cfg.Db.ConnectString, "db-connect", "user=postgres password=%s host=127.0.0.1 dbname=fluidmedia_crm", "DB Connect String")

	flag.Parse()

	pass := os.Getenv("POSTGRES_PASSWORD")
	if pass == "" {
		pass = "Rwbwreia123&"
	}
	cfg.Db.ConnectString = fmt.Sprintf(cfg.Db.ConnectString, pass)

	return cfg
}

func handleIndex(m *model.Model, page *Page, user *model.User, w http.ResponseWriter, r *http.Request) {
	display(w, "index", page, user)
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
			Title: "Search",
			InMenu: false,
			Path: "/search",
			Methods: []string{"GET"},
			Handler: handleSearch,
		},
		{
			Title: "Profile",
			InMenu: false,
			Path: "/profile",
			Methods: []string{"GET", "POST"},
			Handler: handleProfile,
		},
		{
			Title: "2 Factor Authentication",
			InMenu: false,
			Path: "/profile/2fa",
			Methods: []string{"GET", "POST"},
			Handler: handleProfile2FA,
		},

		{
			Title: "Contacts",
			Icon: "contacts",
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

		{
			Title: "Users",
			Icon: "person",
			InMenu: true,
			Path: "/users",
			Methods: []string{"GET"},
			Handler: handleUsers,
			AdminRequired: true,
		},
		{
			Title: "Edit user",
			InMenu: false,
			Path: "/users/{id:[0-9]+}",
			Methods: []string{"GET", "POST"},
			Handler: handleUsersEdit,
		},
		{
			InMenu: false,
			Path: "/users/new",
			Methods: []string{"GET"},
			Handler: handleUsersNew,
		},
		{
			InMenu: false,
			Path: "/users/del/{id:[0-9]+}",
			Methods: []string{"GET"},
			Handler: handleUsersDelete,
		},
	}

	dbInst, err := db.InitDb(cfg.Db)
	if err != nil {
		log.Fatalf("Can't initalize database: %v", err)
	}

	modelInst := model.New(dbInst)

	serveHttp(modelInst, pages, cfg.ListenSpec)
}
