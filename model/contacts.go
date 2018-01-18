package model

type Contact struct {
	ID int
	Name string
	Image string
	Email string
	State int
}

var contactStates = map[int][2]string{
	0: {"Lead", "Leads"},
	1: {"Opportunity", "Opportunities"},
	2: {"Customer", "Customers"},
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

func (m *Model) NewContact() (int, error) {
	return m.db.NewContact()
}

func (m *Model) DeleteContact(id int) error {
	return m.db.DeleteContact(id)
}

func (m *Model) ContactStates() map[int][2]string {
	return contactStates
}

func (c *Contact) StateName() [2]string {
	return contactStates[c.State]
}