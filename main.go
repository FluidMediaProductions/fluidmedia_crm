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
	ClientName string

	Db db.Config
}

func parseFlags() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.ListenSpec, "listen", ":8080", "HTTP listen spec")

	flag.Parse()

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "root"
	}
	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		pass = "Rwbwreia123&"
	}
	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "fluidmedia_crm"
	}

	connString := "%s:%s@tcp(%s)/%s?charset=utf8"
	cfg.Db.ConnectString = fmt.Sprintf(connString, user, pass, host, dbName)
	log.Printf("Connecting to database: %s", cfg.Db.ConnectString)

	cfg.ClientName = os.Getenv("CLIENT_NAME")

	return cfg
}

func handleIndex(m *model.Model, page *Page, user *model.User, w http.ResponseWriter, r *http.Request) {
	type IndexData struct {
		UncontactedLeads int
		UncontactedOpportunities int
	}
	uncontactedLeads, err := m.UncontactedLeads()
	if err != nil {
		display500(w, err)
	}
	uncontactedOpportunities, err := m.UncontactedOpportunities()
	if err != nil {
		display500(w, err)
	}
	displayWithContext(w, "index", page, user, &IndexData{
		UncontactedLeads: uncontactedLeads,
		UncontactedOpportunities: uncontactedOpportunities,
	})
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

	serveHttp(modelInst, pages, cfg)
}
