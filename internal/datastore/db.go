package datastore

import (
	"fmt"
	"log"

	"github.com/ali-shokoohi/micro-gopia/config"
	"gorm.io/gorm"
)

// Database type
type Database struct {
	Gorm *gorm.DB
}

// NewDatabase(db *gorm.DB) Database Return a valid database client
func NewDatabase(dialector gorm.Dialector) *Database {
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(err)
	}

	log.Println("Successfully connected to database!")
	return &Database{
		Gorm: db,
	}
}

// GeneratePsqlInfo Generates psql info by config
func GeneratePsqlInfo() string {
	var (
		dbhost     = config.Confs.Postgres.Host
		dbport     = config.Confs.Postgres.Port
		dbuser     = config.Confs.Postgres.Username
		dbpassword = config.Confs.Postgres.Password
		dbname     = config.Confs.Postgres.DB
		sslmode    = config.Confs.Postgres.SslMode
	)
	return fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s TimeZone=Asia/Tehran",
		dbhost, dbport, dbuser, dbpassword, dbname, sslmode)
}
