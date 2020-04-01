package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/ranabd36/project-qa/config"
	"log"
	"os"
	
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", "./database/migrations", "directory with migration files")
	conf  = config.Get()
)

func main() {
	_ = flags.Parse(os.Args[1:])
	args := flags.Args()
	
	if len(args) < 1 {
		flags.Usage()
		return
	}
	
	command := args[1]
	
	dbstring := getDBString()
	
	if err := goose.SetDialect(conf.Database.DatabaseDriver); err != nil {
		log.Fatalf("goose: failed to set dialect: %v\n", err)
	}
	
	db, err := sql.Open(conf.Database.DatabaseDriver, dbstring)
	if err != nil {
		log.Fatalf("Failed to open DB %v\n", err)
	}
	
	var arguments []string
	if len(args) > 2 {
		arguments = append(arguments, args[2:]...)
	}
	
	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}

func getDBString() string {
	if conf.Database.DatabaseDriver == "mysql" {
		return getMysqlDBString()
	} else if conf.Database.DatabaseDriver == "postgres" {
		return getPostgresDBString()
	}
	return ""
}

func getPostgresDBString() string {
	return fmt.Sprintf("user=%v password=%v dbname=%v sslmode=disable",
		conf.Database.DatabaseUser,
		conf.Database.DatabasePassword,
		conf.Database.DatabaseName,
	)
}

func getMysqlDBString() string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true",
		conf.Database.DatabaseUser,
		conf.Database.DatabasePassword,
		conf.Database.DatabaseHost,
		conf.Database.DatabasePort,
		conf.Database.DatabaseName,
	)
}
