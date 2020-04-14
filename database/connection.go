package database

import (
	"database/sql"
	"fmt"
	"github.com/ranabd36/project-qa/config"
	"log"
)

var (
	Connection *sql.DB
)

func Connect() {
	var err error
	Connection, err = sql.Open(config.Database.Driver, getDBString())
	if err != nil {
		log.Fatalf("Failed while connecting to database, %v\n", err)
	}
	
	if err = Connection.Ping(); err != nil {
		log.Fatalf("Failed to connect with database, %v\n", err)
	}
}

func getDBString() string {
	if config.Database.Driver == "mysql" {
		return getMysqlDBString()
	} else if config.Database.Driver == "postgres" {
		return getPostgresDBString()
	}
	return ""
}

func getPostgresDBString() string {
	return fmt.Sprintf("user=%v password=%v dbname=%v sslmode=disable",
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
	)
}

func getMysqlDBString() string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
	)
}
