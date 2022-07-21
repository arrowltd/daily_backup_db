package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/arrowltd/daily_backup_db/models"
)

var cmd = flag.String("cmd", "", "command line")
var name = flag.String("name", "", "Database name")
var date = flag.String("date", "", "Date")
var host = flag.String("host", "127.0.0.1", "host")
var port = flag.String("port", "5432", "port")
var username = flag.String("username", "", "username")
var password = flag.String("password", "", "username")

func main() {
	flag.Parse()
	log.Println("Process is starting...")
	if *cmd != "" {
		models := &models.Models{}
		models.RunCmd(*cmd)
	} else {
		fmt.Println("Invalid Cmd")
	}
}
