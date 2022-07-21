package models

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fatih/color"

	"github.com/arrowltd/daily_backup_db/date"
	"github.com/arrowltd/daily_backup_db/env"
	"github.com/arrowltd/daily_backup_db/psql/adapter"
	"github.com/arrowltd/daily_backup_db/utils"
	"github.com/mplulu/renv"
)

func (models *Models) dailyBackupDatabase(dbName string) {
	var env *env.ENV
	renv.Parse("", &env)
	adapter := adapter.Initialize(env.DatabaseConfigFilePath, env.Environment)
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
	for _, f := range oldDumpFiles {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}

	dateString := date.TimeToDateStringFileFormat(time.Now())
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
