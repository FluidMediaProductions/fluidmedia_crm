package db

import "github.com/fluidmediaproductions/fluidmedia_crm/model"

func (p *pgDb) SelectContacts() ([]*model.Contact, error) {
	contacts := make([]*model.Contact, 0)
	if err := p.dbConn.Select(&contacts, "SELECT * FROM contacts ORDER BY id"); err != nil {
		return nil, err
	}
	return contacts, nil
}

func (p *pgDb) SelectContact(id int) (*model.Contact, error) {
	var contact model.Contact
	if err := p.dbConn.Get(&contact, "SELECT * FROM contacts WHERE id=?", id); err != nil {
		return nil, err
	}
	return &contact, nil
}

func (p *pgDb) UpdateContact(contact *model.Contact) error {
	_, err := p.dbConn.NamedExec("UPDATE contacts SET name=:name, email=:email," +
		" image=:image, state=:state, contacted_state=:contacted_state, phone=:phone, mobile=:mobile, website=:website," +
		" twitter=:twitter, youtube=:youtube, instagram=:instagram, facebook=:facebook, address=:address, description=:description," +
		" organisation_id=:organisation_id WHERE id=:id", contact)
	return err
}

func (p *pgDb) NewContact() (int, error) {
	res, err := p.dbConn.NamedExec("INSERT INTO contacts (name, email, image, state, contacted_state," +
		" phone, mobile, website, twitter, youtube, instagram, facebook, address, description, organisation_id)" +
		" VALUES (:name, :email, :image, :state, :contacted_state, :phone, :mobile, :website, :twitter, :youtube, :instagram," +
		" :facebook, :address, :description, :organisation_id)", &model.Contact{})
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (p *pgDb) DeleteContact(id int) error {
	_, err := p.dbConn.Exec("DELETE FROM contacts WHERE id=?", id)
	return err
}
