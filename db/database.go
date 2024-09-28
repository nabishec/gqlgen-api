package db

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	dataSourceName string
	DB             *sqlx.DB
}

type DataSourceName struct {
	Protocol string
	Username string
	Password string
	Host     string
	Port     string
	DBname   string
	Options  string
}

func NewDatabase() *Database {
	return &Database{}
}

func (db *Database) ConnectDatabase(config *DataSourceName) error {
	db.dataSourceName = config.Protocol + "://" + config.Username + ":" + config.Password + "@" +
		config.Host + ":" + config.Port + "/" + config.DBname + "?" + config.Options

	var connectError error
	db.DB, connectError = sqlx.Connect("pgx", db.dataSourceName)
	return connectError
}

func (db *Database) PingDatabase() error {
	if db.DB == nil {
		return errors.New("database isn`t established")
	}

	var pingError = db.DB.Ping()
	return pingError
}

func (db *Database) CloseDatabase() error {
	var err = db.DB.Close()
	return err
}

// migrations
func (db *Database) MigrationsUp() error {
	//create migration's driver for sql
	if db.DB == nil {
		return errors.New("database connection not established")
	}

	//create connection for migration
	migrationDB, err := sqlx.Connect("pgx", db.dataSourceName)
	if err != nil {
		log.Println("Failed to create connection for migrations:", err)
		return err
	}

	sqlDatabase := migrationDB.DB
	driver, err := postgres.WithInstance(sqlDatabase, &postgres.Config{})
	if err != nil {
		log.Println("couldn't create driver", err)
		return err
	}
	defer func() {
		if err := driver.Close(); err != nil {
			log.Println("migration's driver couldn't close", err)
		} else {
			log.Println("migration's driver close")
		}
		if err := migrationDB.Close(); err != nil {
			log.Println("migration's connection couldn't close", err)
		} else {
			log.Println("migration's connection close")
		}
	}()

	// create migration's example
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migration",
		"postgres", driver)
	if err != nil {
		log.Println("coudn't create migrate instance", err)
		return err
	}

	//start migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Println("failed to apply migrations", err)
		return err
	}

	if err == migrate.ErrNoChange {
		log.Println("no migrations to apply")
	} else {
		log.Println("migrations applied successfully")
	}

	return nil
}

func (db *Database) MigrationsDown() error {
	if db.DB == nil {
		return errors.New("database connection not established")
	}

	//create connection for migration
	migrationDB, err := sqlx.Connect("pgx", db.dataSourceName)
	if err != nil {
		log.Println("Failed to create connection for migrations:", err)
		return err
	}

	//create migration's driver for sql
	sqlDatabase := migrationDB.DB
	driver, err := postgres.WithInstance(sqlDatabase, &postgres.Config{})
	if err != nil {
		log.Println("couldn't create driver", err)
		return err
	}
	defer func() {
		if err := driver.Close(); err != nil {
			log.Println("migration's driver couldn't close", err)
		} else {
			log.Println("migration's driver close")
		}
		if err := migrationDB.Close(); err != nil {
			log.Println("migration's connection couldn't close", err)
		} else {
			log.Println("migration's connection close")
		}
	}()

	// create migration's example
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migration",
		"postgres", driver)
	if err != nil {
		log.Println("coudn't create migrate instance", err)
		return err
	}

	//start migrations
	err = m.Down()
	if err != nil {
		log.Println("failed to rollback migrations", err)
		return err
	} else {
		log.Println("migrations rllback successfully")
	}

	return nil
}
