package models

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/arrowltd/daily_backup_db/date"
	"github.com/arrowltd/daily_backup_db/env"
	"github.com/arrowltd/daily_backup_db/psql"
	"github.com/arrowltd/daily_backup_db/psql/adapter"
	"github.com/fatih/color"
	"github.com/mplulu/renv"
)

func (models *Models) NewModel() {

}
func (model *Models) restoreDatabase(dateStr string) {
	log.Println("Start the restore process")
	dateStringFormat := date.TimeToDateStringFileFormat(date.DateStringToTime(dateStr))

	filePath := fmt.Sprintf("/tmp/auto_*_%v.dump", dateStringFormat)
	dumpFiles, err := filepath.Glob(filePath)

	if err != nil {
		panic(err)
	}
	if len(dumpFiles) < 1 {
		color.Red("The backup file doesnot exist in %s", dateStr)
	} else {
		var env *env.ENV
		renv.Parse("", &env)

		adapter := adapter.Initialize(env.DatabaseConfigFilePath, env.Environment)
		for _, f := range dumpFiles {

			fileNames := strings.Split(strings.ReplaceAll(f, ".dump", ""), "_")
			dbName := fileNames[1] + fileNames[2]
			//assign database name
			adapter.Database = dbName
			//create database if it doesn't exist
			psql.CreateDb(adapter)

			dbLink := fmt.Sprintf("%v://%v:%v@%v:%v/%v",
				adapter.Type, adapter.Username, adapter.Password,
				adapter.Host, adapter.Port,
				dbName)
			cmdArgs := []string{
				"-d",
				dbLink,
				f,
			}
			cmd := exec.Command("pg_restore", cmdArgs...)
			_, err := cmd.CombinedOutput()
			if err != nil {
				panic(err)
			}
			color.Green("Restore successfully to database: %s", adapter.Database)
		}
	}
	log.Println("End the restore process")
}
