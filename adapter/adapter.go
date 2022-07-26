package adapter

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/fatih/color"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

type Adapter struct {
	Type              string `yaml:"type"` // currently support psql
	Database          string `yaml:"database"`
	Username          string `yaml:"username"`
	Password          string `yaml:"password"`
	Host              string `yaml:"host"`
	Port              string `yaml:"port"`
	MaxIdleConnection int    `yaml:"maxIdleConnection"`
	MaxOpenConnection int    `yaml:"maxOpenConnection"`
}

func Initialize(path string, env string) *Adapter {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Can't not read %v .err #%v ", path, err)))
	}

	envConfig := make(map[string]*Adapter)
	if err := yaml.Unmarshal(yamlFile, envConfig); err != nil {
		panic(errors.New(fmt.Sprintf("Unmarshal: %v", err)))
	}

	adapter, found := envConfig[env]
	if !found {
		panic(errors.New(fmt.Sprintf(" ========== Can not read configurations for '%s' database. =========", env)))
	}

	return adapter
}
func (adapter *Adapter) ConnectToPostgres() *Connection {
	db, err := sql.Open(adapter.Type, fmt.Sprintf("%s://%s:%s@%s:%s?sslmode=disable",
		adapter.Type,
		adapter.Username,
		adapter.Password,
		adapter.Host,
		adapter.Port))
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(adapter.MaxIdleConnection)
	db.SetMaxOpenConns(adapter.MaxOpenConnection)

	if err := db.Ping(); err != nil {
		panic(err)
	}
	color.Green(`Open database connection.`)

	return &Connection{DB: db, Adapter: adapter}
}
func (adapter *Adapter) ConnectToDatabase() *Connection {
	db := connectToDatabase(
		adapter.Type, adapter.Username, adapter.Password, adapter.Host, adapter.Port, adapter.Database,
		adapter.MaxIdleConnection, adapter.MaxOpenConnection)

	return &Connection{DB: db, Adapter: adapter}
}

func connectToDatabase(dbType, username, password, host, port, database string, maxIdleConnection, maxOpenConnection int) *sql.DB {
	db, err := sql.Open(dbType, fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		dbType,
		username,
		password,
		host,
		port,
		database))
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(maxIdleConnection)
	db.SetMaxOpenConns(maxOpenConnection)

	if err := db.Ping(); err != nil {
		panic(err)
	}
	color.Green("Connected to '%s' database at %s:%s\n", database, host, port)

	return db
}
