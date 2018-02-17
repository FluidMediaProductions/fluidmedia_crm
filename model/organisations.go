package model

import (
	"github.com/renstrom/fuzzysearch/fuzzy"
	"strings"
)

type Organisation struct {
	ID          int `db:"id"`
	Name        string `db:"name"`
	Image       string `db:"image"`
	Email       string `db:"email"`
	Phone       string `db:"phone"`
	Website     string `db:"website"`
	Twitter     string `db:"twitter"`
	Youtube     string `db:"youtube"`
	Instagram   string `db:"instagram"`
	Facebook    string `db:"facebook"`
	Address     string `db:"address"`
	Description string `db:"description"`
}

func (m *Model) Organisations() ([]*Organisation, error) {
	return m.db.SelectOrganisations()
}

func (m *Model) Organisation(id int) (*Organisation, error) {
	return m.db.SelectOrganisation(id)
}

func (m *Model) SaveOrganisation(organisation *Organisation) error {
	return m.db.UpdateOrganisation(organisation)
}

func (m *Model) NewOrganisation() (int, error) {
	return m.db.NewOrganisation()
}


func (m *Model) SearchOrganisations(search string) ([]*Organisation, error) {
	search = strings.ToLower(search)
	organisations, err := m.Organisations()
	foundOrganisations := make([]*Organisation, 0)
	if err != nil {
		return nil, err
	}
	for _, organisation := range organisations {
		match := false
		if fuzzy.Match(search, strings.ToLower(organisation.Name)) { match = true }
		if fuzzy.Match(search, strings.ToLower(organisation.Email)) { match = true }
		if fuzzy.Match(search, strings.ToLower(organisation.Phone)) { match = true }
		if fuzzy.Match(search, strings.ToLower(organisation.Website)) { match = true }
		if fuzzy.Match(search, strings.ToLower(organisation.Twitter)) { match = true }
		if fuzzy.Match(search, strings.ToLower(organisation.Facebook)) { match = true }
		if fuzzy.Match(search, strings.ToLower(organisation.Instagram)) { match = true }
		if fuzzy.Match(search, strings.ToLower(organisation.Youtube)) { match = true }
		if fuzzy.Match(search, strings.ToLower(organisation.Address)) { match = true }
		if match {
			foundOrganisations = append(foundOrganisations, organisation)
		}
	}
	return foundOrganisations, nil
}


func (m *Model) DeleteOrganisation(id int) error {
	return m.db.DeleteOrganisation(id)
}
