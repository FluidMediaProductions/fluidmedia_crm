package model

type db interface {
	SelectContacts() ([]*Contact, error)
	SelectContact(int) (*Contact, error)
}

type Model struct {
	db
}

func New(db db) *Model {
	return &Model{
		db: db,
	}
}

func (m *Model) Contacts() ([]*Contact, error) {
	return m.SelectContacts()
}

func (m *Model) Contact(id int) (*Contact, error) {
	return m.SelectContact(id)
}

