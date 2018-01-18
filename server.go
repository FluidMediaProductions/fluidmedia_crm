package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"io/ioutil"
	"strings"
	"path/filepath"
	"html/template"
	"os"
)

func getTemplate() *template.Template {
	searchDir := "templates"

	var fileList []string
	_ = filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})

	return template.Must(template.ParseFiles(fileList...))
}

type Page struct {
	Title string
	Icon string
	Path string
	Methods []string
	Handler func(*Page, http.ResponseWriter, *http.Request)
}

var pages []*Page

type TemplateContext struct {
	Page *Page
	Pages []*Page
}

func display(w http.ResponseWriter, tmpl string, page *Page) {
	err := getTemplate().ExecuteTemplate(w, tmpl, &TemplateContext{Page: page, Pages: pages})
	if err != nil {
		log.Printf("Error rendering %s: %v", tmpl, err)
	}
}


func handlerWrapper(handler func(*Page, http.ResponseWriter, *http.Request), page *Page) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(page, w, r)
	}
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	log.Println(path)
	data, err := ioutil.ReadFile(string(path))

	if err == nil {
		var contentType string

		if strings.HasSuffix(path, ".css") {
			contentType = "text/css"
		} else if strings.HasSuffix(path, ".html") {
			contentType = "text/html"
		} else if strings.HasSuffix(path, ".js") {
			contentType = "application/javascript"
		} else if strings.HasSuffix(path, ".png") {
			contentType = "image/png"
		} else if strings.HasSuffix(path, ".svg") {
			contentType = "image/svg+xml"
		} else {
			contentType = "text/plain"
		}

		w.Header().Add("Content-Type", contentType)
		w.Write(data)
	} else {
		w.WriteHeader(404)
		w.Write([]byte("404 O noes - " + http.StatusText(404)))
	}
}

func serveHttp(pages []*Page) {
	r := mux.NewRouter()

	for _, page := range pages {
		r.Methods(page.Methods...).Path(page.Path).HandlerFunc(handlerWrapper(page.Handler, page))
	}

	http.HandleFunc("/static/", staticHandler)
	http.Handle("/", r)

	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}
