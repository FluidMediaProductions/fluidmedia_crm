package model

type db interface {
	SelectContacts() ([]*Contact, error)
	SelectContact(int) (*Contact, error)
	UpdateContact(*Contact) error
	NewContact() (int, error)
	DeleteContact(int) error
}

type Model struct {
	db db
}

func New(db db) *Model {
	return &Model{
		db: db,
	}
}

