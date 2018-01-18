package main


import (
	"log"
	"os"
	"strings"
	"net/http"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"github.com/gorilla/mux"
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
}

func display(w http.ResponseWriter, tmpl string, data interface{}) {
	getTemplate().ExecuteTemplate(w, tmpl, data)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	display(w, "index", &Page{Title: "Home"})
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

func main() {
	r := mux.NewRouter()

	r.Methods("GET").Path("/").HandlerFunc(handleIndex)

	http.HandleFunc("/static/", staticHandler)
	http.Handle("/", r)

	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}
