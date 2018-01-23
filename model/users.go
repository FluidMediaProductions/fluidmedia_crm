package model

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

type User struct {
	ID         int
	Name       string
	Login      string
	Pass       string
	Email      string
	Phone      string
	IsAdmin    bool
	Disabled   bool
	TotpSecret string `db:"totp_secret"`
}

func (m *Model) Users() ([]*User, error) {
	return m.db.SelectUsers()
}

func (m *Model) User(id int) (*User, error) {
	return m.db.SelectUser(id)
}

func (m *Model) SaveUser(user *User) error {
	if user.Pass != "" {
		log.Printf("Generating password for user %s", user.Login)
		pass, err := bcrypt.GenerateFromPassword([]byte(user.Pass), 14)
		if err != nil {
			log.Printf("Error generating password for user %s: %v", user.Login, err)
			return err
		}
		user.Pass = string(pass)
		err = m.db.UpdateUserPass(user)
		if err != nil {
			log.Printf("Error updating password for user %s: %v", user.Login, err)
			return err
		}
	}
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
