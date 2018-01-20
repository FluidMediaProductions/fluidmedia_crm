package main

import (
	"github.com/fluidmediaproductions/fluidmedia_crm/model"
	"net/http"
	"fmt"
	"log"
	"github.com/gorilla/mux"
	"strconv"
	"database/sql"
)

func handleUsers(m *model.Model, page *Page, user *model.User, w http.ResponseWriter, r *http.Request) {
	type UsersContext struct {
		Users []*model.User
	}
	users, err := m.Users()
	if err != nil {
		log.Printf("Error getting users: %v", err)
		display500(w)
		return
	}
	displayWithContext(w, "users", page, user, &UsersContext{Users: users})
}

func handleUsersEdit(m *model.Model, page *Page, user *model.User, w http.ResponseWriter, r *http.Request) {
	type UserContext struct {
		User *model.User
	}
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	user, err := m.User(id)
	if err == sql.ErrNoRows {
		display404(w)
		return
	} else if err != nil {
		log.Printf("Error getting user: %v", err)
		display500(w)
		return
	}
	if r.Method == "GET" {
		displayWithContext(w, "users-edit", page, user, &UserContext{User: user})
	} else if r.Method == "POST" {
		r.ParseForm()
		newUser := &model.User{
			ID: id,
			Name: r.Form.Get("name"),
			Email: r.Form.Get("email"),
			Phone: r.Form.Get("phone"),
			IsAdmin: r.Form.Get("isadmin") == "checked",
			Login: r.Form.Get("login"),
			Pass: r.Form.Get("pass"),
		}
		m.SaveUser(newUser)
		http.Redirect(w, r, "/users", 302)
	}
}

func handleUsersNew(m *model.Model, page *Page, user *model.User, w http.ResponseWriter, r *http.Request) {
	userId, err := m.NewUser()
	if err != nil {
		log.Printf("Error creating new user: %v", err)
		display500(w)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/users/%d", userId), 302)
}

func handleUsersDelete(m *model.Model, page *Page, user *model.User, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	err := m.DeleteUser(id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		display500(w)
		return
	}
	http.Redirect(w, r, "/users", 302)
}
