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
	"github.com/pquerna/otp/totp"
)

var sessionManager *scs.Manager

var config struct{
	model *model.Model
	clientName string
}

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

type PageHandler func(*model.Model, *Page, *model.User, http.ResponseWriter, *http.Request)

type Page struct {
	Title string
	Icon string
	InMenu bool
	Path string
	Methods []string
	Handler PageHandler
	AdminRequired bool
}

var pages []*Page

type TemplateContext struct {
	Page *Page
	User *model.User
	Pages []*Page
	Context interface{}
	ClientName string
}

func display404(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNotFound)
	return display(w, "404", &Page{Title: "Not Found"}, nil)
}

func display500(w http.ResponseWriter, err error) error {
	w.WriteHeader(http.StatusInternalServerError)
	return displayWithContext(w, "500", &Page{Title: "Error"}, nil, err)
}

func display(w http.ResponseWriter, tmpl string, page *Page, user *model.User) error {
	return displayWithContext(w, tmpl, page, user, nil)
}

func displayWithContext(w http.ResponseWriter, tmpl string, page *Page, user *model.User, context interface{}) error {
	err := getTemplate().ExecuteTemplate(w, tmpl, &TemplateContext{Page: page, Pages: pages, User: user, Context: context, ClientName: config.clientName})
	if err != nil {
		log.Printf("Error rendering %s: %v", tmpl, err)
		return err
	}
	return nil
}


func handlerWrapper(handler PageHandler, page *Page,) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := sessionManager.Load(r)
		userId, err := session.GetInt("userId")
		if err != nil {
			http.Redirect(w, r, "/login", 302)
			return
		}
		user, err := config.model.User(userId)
		if err != nil {
			http.Redirect(w, r, "/login", 302)
			return
		}
		if !user.Disabled {
			if page.AdminRequired && !user.IsAdmin {
				display404(w)
				return
			}
			log.Printf("%s: %s", r.Method, r.URL.Path)
			handler(config.model, page, user, w, r)
		} else {
			http.Redirect(w, r, "/login", 302)
		}
	}
}

func handle404(w http.ResponseWriter, r *http.Request) {
	display404(w)
}

func staticHandler(base string) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:]
		data, err := ioutil.ReadFile(base + "/" + string(path))

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
			} else if strings.HasSuffix(path, ".jpg") {
				contentType = "image/jpeg"
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
	})
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	type LoginContext struct {
		Message string
	}
	session := sessionManager.Load(r)
	if r.Method == "GET" {
		userId, err := session.GetInt("userId")
		if err == nil {
			user, err := config.model.User(userId)
			if err == nil {
				if !user.Disabled {
					http.Redirect(w, r, "/", 302)
					return
				}
			}
		}
		displayWithContext(w, "login", &Page{Title: "Login"}, nil, &LoginContext{"Sign in to start your session"})
	} else if r.Method == "POST" {
		r.ParseForm()
		user, valid := config.model.UserLogin(r.Form.Get("username"), r.Form.Get("password"))
		if !valid {
			log.Printf("User %s failed authentication", r.Form.Get("username"))
			displayWithContext(w, "login", &Page{Title: "Login"}, nil, &LoginContext{"Login failed"})
			return
		}
		if user.TotpSecret != "" {
			key := r.Form.Get("key")
			valid := totp.Validate(key, user.TotpSecret)
			if !valid {
				log.Printf("User %s failed 2fa", user.Login)
				displayWithContext(w, "login", &Page{Title: "Login"}, nil, &LoginContext{"Login failed"})
				return
			}
		}
		err := session.PutInt(w, "userId", user.ID)
		if err != nil {
			display500(w, err)
		}
		log.Printf("User %s successfuly authentiated", user.Login)
		http.Redirect(w, r, "/", 302)
	}
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	session := sessionManager.Load(r)
	userId, err := session.GetInt("userId")
	if err == nil {
		user, err := config.model.User(userId)
		if err == nil {
			if !user.Disabled {
				session.Destroy(w)
			}
		}
	}
	http.Redirect(w, r, "/login", 302)
}


func serveHttp(model *model.Model, pages []*Page, cfg *Config) {
	sessionManager = scs.NewManager(model.NewSessionStore())
	sessionManager.Lifetime(time.Hour * 24 * 30)
	sessionManager.Persist(true)
	//sessionManager.Secure(true)

	config.model = model
	config.clientName = cfg.ClientName

	r := mux.NewRouter()

	for _, page := range pages {
		r.Methods(page.Methods...).Path(page.Path).HandlerFunc(handlerWrapper(page.Handler, page))
	}
	r.Methods("GET", "POST").Path("/login").HandlerFunc(handleLogin)
	r.Methods("GET").Path("/logout").HandlerFunc(handleLogout)

	r.NotFoundHandler = http.HandlerFunc(handle404)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", staticHandler("static")))
	r.PathPrefix("/data/").Handler(http.StripPrefix("/data", staticHandler("data")))


	log.Println("Listening...")
	http.ListenAndServe(cfg.ListenSpec, sessionManager.Use(r))
}
