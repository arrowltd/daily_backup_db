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
