package models

import (
	"flag"
	"fmt"
)

type Models struct {
}

func (models *Models) RunCmd(cmd string) {
	switch cmd {
	case "dailybackupdb":
		dbName := flag.Lookup("name").Value.(flag.Getter).Get().(string)
		if dbName != "" {
			go models.dailyBackupDatabase(dbName)
			select {}
		} else {
			fmt.Println("Database name cannot be empty")
		}
	case "restoredb":
		dateStr := flag.Lookup("date").Value.(flag.Getter).Get().(string)
		models.restoreDatabase(dateStr)
	default:
		fmt.Println("Invalid the command")
	}
}
