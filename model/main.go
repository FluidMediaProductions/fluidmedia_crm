package model

import "github.com/alexedwards/scs/stores/pgstore"

type db interface {
	SelectContacts() ([]*Contact, error)
	SelectContact(int) (*Contact, error)
	UpdateContact(*Contact) error
	NewContact() (int, error)
	DeleteContact(int) error

	SelectOrganisations() ([]*Organisation, error)
	SelectOrganisation(int) (*Organisation, error)
	UpdateOrganisation(*Organisation) error
	NewOrganisation() (int, error)
	DeleteOrganisation(int) error

	SelectUsers() ([]*User, error)
	SelectUser(int) (*User, error)
	UpdateUser(*User) error
	UpdateUserPass(*User) error
	NewUser() (int, error)
	DeleteUser(int) error

	SessionStore() *pgstore.PGStore
}

type Model struct {
	db db
}

func New(db db) *Model {
	return &Model{
		db: db,
	}
}

func (m *Model) NewSessionStore() *pgstore.PGStore {
	return m.db.SessionStore()
}

