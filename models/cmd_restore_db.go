package models

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/arrowltd/daily_backup_db/adapter"
	"github.com/arrowltd/daily_backup_db/date"
	"github.com/fatih/color"
)

func (models *Models) NewModel() {

}
func (model *Models) restoreDatabase(host, port, username, password, dbName, dateStr string) {
	log.Println("Start the restore process")
	dateStringFormat := date.TimeToDateStringFileFormat(date.DateStringToTime(dateStr, "YYY/MM/DD"))
	dumpFile := fmt.Sprintf("auto_%v_%v.dump", dbName, dateStringFormat)
	filePath := fmt.Sprintf("/tmp/%v", dumpFile)
	fmt.Println(filePath)
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		color.Red("Not found any %s backup file in %s", dbName, dateStr)
	} else {
		adapter := &adapter.Adapter{
			Type:              "postgres",
			Database:          dbName,
			Host:              host,
			Port:              port,
			Username:          username,
			Password:          password,
			MaxIdleConnection: 80,
			MaxOpenConnection: 40,
		}
		conn := adapter.ConnectToPostgres()
		//create database if it doesn't exist
		if conn.DoesDatabaseExist() {
			conn.DropDatabase()
		}
		err := conn.CreateDatabase()
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		dbLink := fmt.Sprintf("%v://%v:%v@%v:%v/%v",
			adapter.Type, adapter.Username, adapter.Password,
			adapter.Host, adapter.Port,
			adapter.Database)
		cmdArgs := []string{
			"-d",
			dbLink,
			filePath,
		}
		cmd := exec.Command("pg_restore", cmdArgs...)
		output, err := cmd.CombinedOutput()
		fmt.Println(string(output))
		if err != nil {
			panic(err)
		}
		color.Green("Restore successfully to database: %s", adapter.Database)
	}
	log.Println("End the restore process")
}
