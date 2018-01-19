package model

type Contact struct {
	ID int
	Name string
	Image string
	State int
	Email string
	Phone string
	Mobile string
	Website string
	Twitter string
	Address string
	Description string
	OrganisationId int `db:"organisation_id"`
	Organisation *Organisation
}

var contactStates = map[int][2]string{
	0: {"Lead", "Leads"},
	1: {"Opportunity", "Opportunities"},
	2: {"Customer", "Customers"},
}

func (m *Model) Contacts() ([]*Contact, error) {
	contacts, err := m.db.SelectContacts()
	if err != nil {
		return nil, err
	}
	for _, contact := range contacts {
		organisation, err := m.Organisation(contact.OrganisationId)
		if err != nil {
			organisation = &Organisation{}
		}
		contact.Organisation = organisation
	}
	return contacts, nil
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