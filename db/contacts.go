package db

import "github.com/fluidmediaproductions/fluidmedia_crm/model"

func (p *pgDb) createContactsTablesIfNotExist() error {
	createSql := `
       CREATE TABLE IF NOT EXISTS contacts (
       id SERIAL NOT NULL PRIMARY KEY,
       image TEXT NOT NULL,
       name TEXT NOT NULL,
       state INTEGER NOT NULL,
       email TEXT NOT NULL,
       phone TEXT NOT NULL,
       mobile TEXT NOT NULL,
       website TEXT NOT NULL,
       twitter TEXT NOT NULL,
       address TEXT NOT NULL,
       description TEXT NOT NULL);
    `
	if rows, err := p.dbConn.Query(createSql); err != nil {
		return err
	} else {
		rows.Close()
	}
	return nil
}

func (p *pgDb) prepareContactsSqlStatements() (err error) {
	if p.sqlSelectContacts, err = p.dbConn.Preparex("SELECT * FROM contacts ORDER BY id"); err != nil { return err}
	if p.sqlSelectContact, err = p.dbConn.Preparex("SELECT * FROM contacts WHERE id=$1"); err != nil { return err }
	if p.sqlUpdateContact, err = p.dbConn.PrepareNamed("UPDATE contacts SET name=:name, email=:email," +
		" image=:image, state=:state, phone=:phone, mobile=:mobile, website=:website, twitter=:twitter, address=:address," +
		" description=:description WHERE id=:id"); err != nil { return err }
	if p.sqlInsertContact, err = p.dbConn.PrepareNamed("INSERT INTO contacts (name, email, image, state, phone," +
		" mobile, website, twitter, address, description)" +
		" VALUES (:name, :email, :image, :state, :phone, :mobile, :website, :twitter, :address, :description) RETURNING id"); err != nil { return err }
	if p.sqlDeleteContact, err = p.dbConn.Preparex("DELETE FROM contacts WHERE id=$1"); err != nil { return err }
    return nil
}

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

func (p *pgDb) NewContact() (int, error) {
	id := 0
	err := p.sqlInsertContact.QueryRow(&model.Contact{}).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (p *pgDb) DeleteContact(id int) error {
	_, err := p.sqlDeleteContact.Exec(id)
	return err
}