package model

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID int
	Name string
	Login string
	Pass string
	Email string
	Phone string
	IsAdmin bool
	Disabled bool
}

func (m *Model) Users() ([]*User, error) {
	return m.db.SelectUsers()
}

func (m *Model) User(id int) (*User, error) {
	return m.db.SelectUser(id)
}

func (m *Model) SaveUser(user *User) error {
	return m.db.UpdateUser(user)
}

func (m *Model) NewUser() (int, error) {
	return m.db.NewUser()
}

func (m *Model) DeleteUser(id int) error {
	return m.db.DeleteUser(id)
}

func (m *Model) UserLogin(name string, pass string) (*User, bool) {
	users, err := m.db.SelectUsers()
	if err != nil {
		return nil, false
	}
	for _, user := range users {
		if user.Login == name {
			err := bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(pass))
			if err == nil {
				return user, true
			}
			return nil, false
		}
	}
	return nil, false
}