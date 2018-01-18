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
	return nil
}

func (p *pgDb) SelectContacts() ([]*model.Contact, error){
	people := make([]*model.Contact, 0)
	if err := p.sqlSelectContacts.Select(&people); err != nil {
		return nil, err
	}
	return people, nil
}