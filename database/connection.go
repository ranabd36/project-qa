package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/ranabd36/project-qa/config"
	"log"
)

func Connect() (*sql.DB, error) {
	
	db, err := sql.Open(config.Database.Driver, GetDBString())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func GetDBString() string {
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
