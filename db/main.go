package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
	"github.com/alexedwards/scs/stores/pgstore"
	"time"
)

type Config struct {
	ConnectString string
}

type pgDb struct {
	dbConn *sqlx.DB

	sqlSelectContacts *sqlx.Stmt
	sqlSelectContact  *sqlx.Stmt
	sqlUpdateContact  *sqlx.NamedStmt
	sqlInsertContact  *sqlx.NamedStmt
	sqlDeleteContact  *sqlx.Stmt

	sqlSelectOrganisations *sqlx.Stmt
	sqlSelectOrganisation  *sqlx.Stmt
	sqlUpdateOrganisation  *sqlx.NamedStmt
	sqlInsertOrganisation  *sqlx.NamedStmt
	sqlDeleteOrganisation  *sqlx.Stmt

	sqlSelectUsers    *sqlx.Stmt
	sqlSelectUser     *sqlx.Stmt
	sqlUpdateUser     *sqlx.NamedStmt
	sqlUpdateUserPass *sqlx.NamedStmt
	sqlInsertUser     *sqlx.NamedStmt
	sqlDeleteUser     *sqlx.Stmt
}

func InitDb(cfg Config) (*pgDb, error) {
	if dbConn, err := sqlx.Connect("mysql", cfg.ConnectString); err != nil {
		return nil, err
	} else {
		p := &pgDb{dbConn: dbConn}
		if err := p.dbConn.Ping(); err != nil {
			return nil, err
		}
		if err := p.migrate(); err != nil {
			return nil, err
		}
		if err := p.prepareSqlStatements(); err != nil {
			return nil, err
		}
		return p, nil
	}
}

func (p *pgDb) migrate() error {
	driver, err := postgres.WithInstance(p.dbConn.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
	}
	return nil
}

func (p *pgDb) prepareSqlStatements() (err error) {
	if err := p.prepareContactsSqlStatements(); err != nil {
		return err
	}
	if err := p.prepareOrganisationsSqlStatements(); err != nil {
		return err
	}
	if err := p.prepareUsersSqlStatements(); err != nil {
		return err
	}
	return nil
}

func (p *pgDb) SessionStore() *pgstore.PGStore {
	return pgstore.New(p.dbConn.DB, time.Minute)
}
