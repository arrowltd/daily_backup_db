package psql

import "github.com/arrowltd/daily_backup_db/psql/adapter"

func CreateDb(adapter *adapter.Adapter) error {
	conn := adapter.ConnectToPostgres()
	return create(conn)
}

func create(connection *adapter.Connection) error {
	defer connection.Close()
	return connection.CreateDatabaseIfNotExists()

}
