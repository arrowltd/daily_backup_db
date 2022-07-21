package models

import (
	"flag"
	"fmt"
)

type Models struct {
}

func (models *Models) RunCmd(cmd string) {
	switch cmd {
	case "backup":
		dbName := flag.Lookup("dbname").Value.(flag.Getter).Get().(string)
		host := flag.Lookup("host").Value.(flag.Getter).Get().(string)
		port := flag.Lookup("port").Value.(flag.Getter).Get().(string)
		username := flag.Lookup("username").Value.(flag.Getter).Get().(string)
		password := flag.Lookup("password").Value.(flag.Getter).Get().(string)
		if dbName == "" {
			fmt.Println("Database name cannot be empty")
		} else if username == "" {
			fmt.Println("username cannot be empty")
		} else {
			go models.dailyBackupDatabase(host, port, username, password, dbName)
			select {}
		}
	case "restore":
		host := flag.Lookup("host").Value.(flag.Getter).Get().(string)
		port := flag.Lookup("port").Value.(flag.Getter).Get().(string)
		username := flag.Lookup("username").Value.(flag.Getter).Get().(string)
		password := flag.Lookup("password").Value.(flag.Getter).Get().(string)
		dbname := flag.Lookup("dbname").Value.(flag.Getter).Get().(string)
		dateStr := flag.Lookup("date").Value.(flag.Getter).Get().(string)
		if dbname == "" {
			fmt.Println("Database name cannot be empty")
		} else if username == "" {
			fmt.Println("username cannot be empty")
		} else {
			models.restoreDatabase(host, port, username, password, dbname, dateStr)
		}

	default:
		fmt.Println("Invalid the command")
	}
}
