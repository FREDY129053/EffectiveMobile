package database

import (
	"fmt"
	_ "fmt"
	"subscriptions/rest-service/internal/models"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDBConnect() (db *gorm.DB, err error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		viper.GetString("DB_HOST"),
		viper.GetString("POSTGRES_USER"),
		viper.GetString("POSTGRES_PASSWORD"),
		viper.GetString("POSTGRES_DB"),
		viper.GetString("DB_PORT"),
	)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	var isDBExist bool
	t := conn.Raw(fmt.Sprintf("SELECT datname FROM pg_catalog.pg_database WHERE lower(datname) = lower('%s');", viper.GetString("POSTGRES_DB")))
	
	if err = t.Row().Scan(&isDBExist); err != nil {
		conn.Exec(fmt.Sprintf("CREATE DATABASE %s", viper.GetString("POSTGRES_DB")))
	}

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	db.AutoMigrate(&models.Subscription{})
	
	return
}
