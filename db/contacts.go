package db

import "github.com/fluidmediaproductions/fluidmedia_crm/model"

func (p *pgDb) SelectContacts() ([]*model.Contact, error){
	contacts := make([]*model.Contact, 0)
	if err := p.sqlSelectContacts.Select(&contacts); err != nil {
		return nil, err
	}
	return contacts, nil
}

func (p *pgDb) SelectContact(id int) (*model.Contact, error){
	var contact model.Contact
	if err := p.sqlSelectContact.Get(&contact, id); err != nil {
		return nil, err
	}
	return &contact, nil
}

func (p *pgDb) UpdateContact(contact *model.Contact) error {
	_, err := p.sqlUpdateContact.Exec(contact)
	return err
}