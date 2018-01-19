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
	"github.com/alexedwards/scs"
	"github.com/fluidmediaproductions/fluidmedia_crm/model"
	"time"
)

var sessionManager *scs.Manager

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

type PageHandler func(*model.Model, *Page, http.ResponseWriter, *http.Request)

type Page struct {
	Title string
	Icon string
	InMenu bool
	Path string
	Methods []string
	Handler PageHandler
}

var pages []*Page

type TemplateContext struct {
	Page *Page
	Pages []*Page
	Context interface{}
}

func display404(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNotFound)
	return display(w, "404", &Page{Title: "Not Found"})
}

func display500(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusInternalServerError)
	return display(w, "500", &Page{Title: "Error"})
}

func display(w http.ResponseWriter, tmpl string, page *Page) error {
	return displayWithContext(w, tmpl, page, nil)
}

func displayWithContext(w http.ResponseWriter, tmpl string, page *Page, context interface{}) error {
	err := getTemplate().ExecuteTemplate(w, tmpl, &TemplateContext{Page: page, Pages: pages, Context: context})
	if err != nil {
		log.Printf("Error rendering %s: %v", tmpl, err)
		return err
	}
	return nil
}


func handlerWrapper(handler PageHandler, model *model.Model, page *Page) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := sessionManager.Load(r)
		authed, _ := session.GetBool("authed")
		if authed {
			handler(model, page, w, r)
		} else {
			http.Redirect(w, r, "/login", 302)
		}
	}
}

func handle404(w http.ResponseWriter, r *http.Request) {
	display404(w)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	log.Println(path)
	data, err := ioutil.ReadFile("static/" + string(path))

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
		display404(w)
	}
}

func serveHttp(model *model.Model, pages []*Page, listenSpec string) {
	sessionManager = scs.NewManager(model.NewSessionStore())
	sessionManager.Lifetime(time.Hour * 24 * 30)
	sessionManager.Persist(true)
	sessionManager.Secure(true)

	r := mux.NewRouter()

	for _, page := range pages {
		r.Methods(page.Methods...).Path(page.Path).HandlerFunc(handlerWrapper(page.Handler, model, page))
	}

	r.NotFoundHandler = http.HandlerFunc(handle404)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.HandlerFunc(staticHandler)))


	log.Println("Listening...")
	http.ListenAndServe(listenSpec, sessionManager.Use(r))
}
