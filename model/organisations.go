package model

type Organisation struct {
	ID          int
	Name        string
	Image       string
	Email       string
	Phone       string
	Website     string
	Twitter     string
	Youtube     string
	Instagram   string
	Facebook    string
	Address     string
	Description string
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

func (m *Model) DeleteOrganisation(id int) error {
	return m.db.DeleteOrganisation(id)
}
