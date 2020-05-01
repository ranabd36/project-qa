package main

import (
	"database/sql"
	"flag"
	"github.com/ranabd36/project-qa/config"
	"github.com/ranabd36/project-qa/database"
	"log"
	"os"
	
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", "./database/migrations", "directory with migration files")
)

func main() {
	config.Load()
	flags.Usage = usage
	_ = flags.Parse(os.Args[1:])
	
	args := flags.Args()
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}
	
	switch args[0] {
	case "create":
		if err := goose.Run("create", nil, *dir, args[1:]...); err != nil {
			log.Fatalf("goose run: %v", err)
		}
		return
	case "fix":
		if err := goose.Run("fix", nil, *dir); err != nil {
			log.Fatalf("goose run: %v", err)
		}
		return
	case "-h", "--help":
		flags.Usage()
		return
	}
	
	if len(args) < 1 {
		flags.Usage()
		return
	}
	
	command := args[0]
	driver := config.Database.Driver
	dbString := database.GetDBString()
	
	switch driver {
	case "postgres", "mysql", "sqlite3", "redshift":
		if err := goose.SetDialect(driver); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("%q driver not supported\n", driver)
	}
	
	if dbString == "" {
		log.Fatalf("-dbString=%q not supported\n", dbString)
	}
	
	if driver == "redshift" {
		driver = "postgres"
	}
	
	db, err := sql.Open(driver, dbString)
	if err != nil {
		log.Fatalf("Invalid DB string:%q %v\n", dbString, err)
	}
	
	var arguments []string
	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}
	
	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		log.Fatalf("goose run: %v", err)
	}
}

func usage() {
	log.Print(usagePrefix)
	flags.PrintDefaults()
	log.Print(usageCommands)
}

var (
	usagePrefix = `Usage: migrate [OPTIONS] COMMAND
Drivers:
    postgres
    mysql
    sqlite3
    redshift
Examples:
    goose status
    goose create init sql
    goose create add_some_column sql
    goose create fetch_user_data go
    goose up
    goose down
    goose redo
Options:
	-dir string
		directory with migration files (default ".")
`
	
	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
	fix                  Apply sequential ordering to migrations
`
)

