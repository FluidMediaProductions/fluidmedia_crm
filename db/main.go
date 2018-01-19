package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	ConnectString string
}

type pgDb struct {
	dbConn *sqlx.DB

	sqlSelectContacts *sqlx.Stmt
	sqlSelectContact *sqlx.Stmt
	sqlUpdateContact *sqlx.NamedStmt
	sqlInsertContact *sqlx.NamedStmt
	sqlDeleteContact *sqlx.Stmt
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
	if err := p.createContactTablesIfNotExist(); err != nil { return err }
	return nil
}

func (p *pgDb) prepareSqlStatements() (err error) {
	if err := p.prepareContactsSqlStatements(); err != nil { return err }
	return nil
}