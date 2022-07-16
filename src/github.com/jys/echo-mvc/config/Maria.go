package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"mvc/domain"
	"os"
)

type config struct {
	dbConnection string
	dbUrl        string
	dbDatabase   string
	dbUser       string
	dbPassword   string
}

func Init() {
	config := createConfig()

	dsn := config.dbUser + ":" + config.dbPassword + "@tcp(" + config.dbUrl + ")/" + config.dbDatabase + "?charset=utf8&parseTime=True&loc=Local"

	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	_ = db.AutoMigrate(&domain.Music{})
}

func createConfig() config {
	c := config{
		dbConnection: os.Getenv("DB_CONNECTION"),
		dbUrl:        os.Getenv("DB_URL"),
		dbDatabase:   os.Getenv("DB_DATABASE"),
		dbUser:       os.Getenv("DB_USER"),
		dbPassword:   os.Getenv("DB_PASSWORD"),
	}

	return c
}
