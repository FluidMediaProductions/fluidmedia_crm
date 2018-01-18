package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/fluidmediaproductions/fluidmedia_crm/model"
)

type Config struct {
	ConnectString string
}

type pgDb struct {
	dbConn *sqlx.DB

	sqlSelectContacts *sqlx.Stmt
	sqlSelectContact *sqlx.Stmt
}

func InitDb(cfg Config) (*pgDb, error) {
	if dbConn, err := sqlx.Connect("postgres", cfg.ConnectString); err != nil {
		return nil, err
	} else {
		p := &pgDb{dbConn: dbConn}
		if err := p.dbConn.Ping(); err != nil {
			return nil, err
		}
		if err := p.createTablesIfNotExist(); err != nil {
			return nil, err
		}
		if err := p.prepareSqlStatements(); err != nil {
			return nil, err
		}
		return p, nil
	}
}

func (p *pgDb) createTablesIfNotExist() error {
	createSql := `
       CREATE TABLE IF NOT EXISTS contacts (
       id SERIAL NOT NULL PRIMARY KEY,
       image TEXT NOT NULL,
       name TEXT NOT NULL,
       email TEXT);
    `
	if rows, err := p.dbConn.Query(createSql); err != nil {
		return err
	} else {
		rows.Close()
	}
	return nil
}

func (p *pgDb) prepareSqlStatements() (err error) {
	if p.sqlSelectContacts, err = p.dbConn.Preparex("SELECT * FROM contacts"); err != nil {
		return err
	}
	if p.sqlSelectContact, err = p.dbConn.Preparex(`SELECT * FROM contacts WHERE id=$1`); err != nil {
		return err
	}
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