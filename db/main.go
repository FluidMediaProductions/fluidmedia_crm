package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/mysql"
	_ "github.com/mattes/migrate/source/file"
	"time"
	"github.com/alexedwards/scs/stores/mysqlstore"
)

type Config struct {
	ConnectString string
}

type pgDb struct {
	dbConn *sqlx.DB
}

func InitDb(cfg Config) (*pgDb, error) {
	if dbConn, err := sqlx.Open("mysql", cfg.ConnectString); err != nil {
		return nil, err
	} else {
		dbConn.SetMaxIdleConns(0)
		p := &pgDb{dbConn: dbConn}
		if err := p.dbConn.Ping(); err != nil {
			return nil, err
		}
		if err := p.migrate(); err != nil {
			return nil, err
		}
		return p, nil
	}
}

func (p *pgDb) migrate() error {
	driver, err := mysql.WithInstance(p.dbConn.DB, &mysql.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql", driver)
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

func (p *pgDb) SessionStore() *mysqlstore.MySQLStore {
	return mysqlstore.New(p.dbConn.DB, time.Minute)
}
