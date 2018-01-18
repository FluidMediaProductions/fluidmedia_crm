package main

import (
	"net/http"
)

func handleIndex(page *Page, w http.ResponseWriter, r *http.Request) {
	display(w, "index", page)
}

type Contact struct {
	ID int
	Name string
	Image string
	Email string
}

var contacts []*Contact

func handleContacts(page *Page, w http.ResponseWriter, r *http.Request) {
	type ContactsContext struct {
		Contacts []*Contact
	}
	displayWithContext(w, "contacts", page, &ContactsContext{Contacts: contacts})
}

func main() {
	pages = []*Page{
		{
			Title: "Home",
			Icon: "home",
			Path: "/",
			Methods: []string{"GET"},
			Handler: handleIndex,
		},
		{
			Title: "Contacts",
			Icon: "person",
			Path: "/contacts",
			Methods: []string{"GET"},
			Handler: handleContacts,
		},
	}

	contacts = []*Contact{
		{
			ID: 0,
			Name: "Mark Otto",
			Image: "user.png",
			Email: "mdo@example.com",
		},
		{
			ID: 1,
			Name: "Jacob Thornton",
			Image: "user.png",
			Email: "fat@example.com",
		},
		{
			ID: 2,
			Name: "Larry the Bird",
			Image: "user.png",
			Email: "twitter@example.com",
		},
		{
			ID: 3,
			Name: "Larry Jellybean",
			Image: "user.png",
			Email: "lajelly@example.com",
		},
		{
			ID: 4,
			Name: "Larry Kikat",
			Image: "user.png",
			Email: "lakitkat@example.com",
		},
	}

	serveHttp(pages)
}
