package model

import (
	"github.com/renstrom/fuzzysearch/fuzzy"
	"strings"
)

type Contact struct {
	ID             int `db:"id"`
	Name           string `db:"name"`
	Image          string `db:"image"`
	State          int `db:"state"`
	ContactedState int `db:"contacted_state"`
	Email          string `db:"email"`
	Phone          string `db:"phone"`
	Mobile         string `db:"mobile"`
	Website        string `db:"website"`
	Twitter        string `db:"twitter"`
	Youtube        string `db:"youtube"`
	Instagram      string `db:"instagram"`
	Facebook       string `db:"facebook"`
	Address        string `db:"address"`
	Description    string `db:"description"`
	OrganisationId int `db:"organisation_id"`
	Organisation   *Organisation
}

var contactStates = map[int][2]string{
	0: {"Lead", "Leads"},
	1: {"Opportunity", "Opportunities"},
	2: {"Customer", "Customers"},
}

var contactedStates = map[int]string{
	0: "Not contacted",
	1: "Attempted contact",
	2: "Contacted",
	3: "Disqualified",
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

func (m *Model) ContactedStates() map[int]string {
	return contactedStates
}


func (m *Model) SearchContacts(search string) ([]*Contact, error) {
	search = strings.ToLower(search)
	contacts, err := m.Contacts()
	foundContacts := make([]*Contact, 0)
	if err != nil {
		return nil, err
	}
	for _, contact := range contacts {
		match := false
		if fuzzy.Match(search, strings.ToLower(contact.Name)) { match = true }
		if fuzzy.Match(search, strings.ToLower(contact.Email)) { match = true }
		if fuzzy.Match(search, strings.ToLower(contact.Phone)) { match = true }
		if fuzzy.Match(search, strings.ToLower(contact.Mobile)) { match = true }
		if fuzzy.Match(search, strings.ToLower(contact.Website)) { match = true }
		if fuzzy.Match(search, strings.ToLower(contact.Twitter)) { match = true }
		if fuzzy.Match(search, strings.ToLower(contact.Facebook)) { match = true }
		if fuzzy.Match(search, strings.ToLower(contact.Instagram)) { match = true }
		if fuzzy.Match(search, strings.ToLower(contact.Youtube)) { match = true }
		if fuzzy.Match(search, strings.ToLower(contact.Address)) { match = true }
		if match {
			foundContacts = append(foundContacts, contact)
		}
	}
	return foundContacts, nil
}

func (m *Model) UncontactedLeads() (int, error) {
	var count int
	contacts, err := m.Contacts()
	if err != nil {
		return 0, err
	}
	for _, c := range contacts {
		if c.State == 0 && c.ContactedState == 0 {
			count += 1
		}
	}
	return count, nil
}

func (m *Model) UncontactedOpportunities() (int, error) {
	var count int
	contacts, err := m.Contacts()
	if err != nil {
		return 0, err
	}
	for _, c := range contacts {
		if c.State == 1 && c.ContactedState == 0 {
			count += 1
		}
	}
	return count, nil
}

func (c *Contact) StateName() [2]string {
	return contactStates[c.State]
}

func (c *Contact) ContactedStateName() string {
	return contactedStates[c.ContactedState]
}
