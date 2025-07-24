package database

import (
	_ "fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDBConnect() (db *gorm.DB, err error) {
	// TODO: .env vars

	dsn := "host=localhost user=fredy password=postgres port=5432 sslmode=disable"
	// conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// var isDBExist bool
	// t := conn.Raw(fmt.Sprintf("SELECT datname FROM pg_catalog.pg_database WHERE lower(datname) = lower('%s');", "test"))
	
	// if err = t.Row().Scan(&isDBExist); err != nil {
	// 	conn.Exec(fmt.Sprintf("CREATE DATABASE %s", "test"))
	// }

	db, err = gorm.Open(postgres.Open(dsn + " dbname=test"))
	
	return
}
