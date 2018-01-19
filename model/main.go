package model

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
}

type Model struct {
	db db
}

func New(db db) *Model {
	return &Model{
		db: db,
	}
}

