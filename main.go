package main


import (
	"net/http"
)

func handleIndex(page *Page, w http.ResponseWriter, r *http.Request) {
	display(w, "index", page)
}

func handleContacts(page *Page, w http.ResponseWriter, r *http.Request) {
	display(w, "contacts", page)
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

	serveHttp(pages)
}
