package adapter

import (
	"database/sql"
	"fmt"

	"github.com/fatih/color"
)

type Connection struct {
	*sql.DB
	*Adapter
}

func (conn *Connection) CreateDatabase() error {
	if _, err := conn.DB.Exec(fmt.Sprintf("CREATE DATABASE %s;", conn.Database)); err != nil {
		panic(err)
	}
	color.Green("Created '%s' database. \n", conn.Database)
	return nil
}
func (conn *Connection) DropDatabase() error {
	if _, err := conn.DB.Exec(fmt.Sprintf("DROP DATABASE %s;", conn.Database)); err != nil {
		panic(err)
	}

	color.Yellow("Database '%s' dropped. \n", conn.Database)
	return nil
}
func (conn *Connection) DoesDatabaseExist() bool {
	row := conn.DB.QueryRow(`
		SELECT EXISTS(	SELECT datname FROM pg_catalog.pg_database WHERE datname = $1);`, conn.Database)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		panic(err)
	}

	return exists
}
