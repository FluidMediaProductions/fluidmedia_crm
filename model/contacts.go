package model

type Contact struct {
	ID int
	Name string
	Image string
	Email string
}

func (m *Model) Contacts() ([]*Contact, error) {
	return m.db.SelectContacts()
}

func (m *Model) Contact(id int) (*Contact, error) {
	return m.db.SelectContact(id)
}

func (m *Model) SaveContact(contact *Contact) error {
	return m.db.UpdateContact(contact)
}
