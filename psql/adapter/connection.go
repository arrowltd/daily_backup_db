package adapter

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/fatih/color"
)

type Connection struct {
	*sql.DB
	*Adapter
}

func (conn *Connection) CreateDatabaseIfNotExists() error {
	if conn.doesDatabaseExist() {
		color.Red("Database '%s' already exists. \n", conn.Database)
		return errors.New("err:database_already_exist")
	}

	if _, err := conn.DB.Exec(fmt.Sprintf("CREATE DATABASE %s;", conn.Database)); err != nil {
		panic(err)
	}

	color.Green("Created '%s' database. \n", conn.Database)
	return nil
}
func (conn *Connection) doesDatabaseExist() bool {
	row := conn.DB.QueryRow(`
		SELECT EXISTS(	SELECT datname FROM pg_catalog.pg_database WHERE datname = $1);`, conn.Database)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		panic(err)
	}

	return exists
}
