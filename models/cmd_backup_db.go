package models

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/arrowltd/daily_backup_db/adapter"
	"github.com/arrowltd/daily_backup_db/date"
	"github.com/arrowltd/daily_backup_db/utils"
	"github.com/fatih/color"
)

func (models *Models) dailyBackupDatabase(host, port, username, password, dbName string) {
	adapter := &adapter.Adapter{
		Type:              "postgres",
		Database:          dbName,
		Username:          username,
		Password:          password,
		Host:              host,
		Port:              port,
		MaxIdleConnection: 80,
		MaxOpenConnection: 40,
	}
	adapter.Database = dbName
	adapter.ConnectToDatabase()

	go utils.IntervalEverydayAt(fmt.Sprintf("Daily backup database %v", dbName), 23, 0, func() {
		models.backupDatabase(adapter)
	})

}

func (models *Models) backupDatabase(adapter *adapter.Adapter) {
	log.Println("Database backup daily is starting......")
	oldDumpFiles, err := filepath.Glob(fmt.Sprintf("/tmp/auto_%v_*", adapter.Database))
	if err != nil {
		panic(err)
	}
	//delete all database dump file
	afterDateString := date.TimeToDateStringFileFormat(date.Now().AddDate(0, 0, -3))
	afterDateInt, err := strconv.Atoi(afterDateString)
	if err != nil {
		panic(err)
	}
	for _, f := range oldDumpFiles {
		dateFileString := strings.Split(strings.ReplaceAll(f, ".dump", ""), "_")[2]
		dateFileInt, err := strconv.Atoi(dateFileString)
		if err != nil {
			panic(err)
		}
		if dateFileInt <= afterDateInt {
			if err := os.Remove(f); err != nil {
				panic(err)
			}
		}

	}
	dateString := date.TimeToDateStringFileFormat(date.Now())
	dumpFile := fmt.Sprintf("/tmp/auto_%v_%v.dump", adapter.Database, dateString)
	dbLink := fmt.Sprintf("%v://%v:%v@%v:%v/%v",
		adapter.Type, adapter.Username, adapter.Password,
		adapter.Host, adapter.Port,
		adapter.Database)

	cmdArgs := []string{
		"--no-owner",
		"--dbname", dbLink,
		"-f", dumpFile,
		"-Fc",
	}
	cmd := exec.Command("pg_dump", cmdArgs...)
	_, err = cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	color.Green("Backup %s to dumpfile %s successfully", adapter.Database, dumpFile)
}
