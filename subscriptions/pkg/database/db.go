package database

import (
	"fmt"
	_ "fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDBConnect() (db *gorm.DB, err error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s port=%s sslmode=disable",
		viper.GetString("DB_HOST"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_PORT"),
	)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	conn.Exec(fmt.Sprintf("CREATE DATABASE %s", viper.GetString("DB_NAME")))

	db, err = gorm.Open(postgres.Open(dsn + " dbname=" + viper.GetString("DB_NAME")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	
	return
}
