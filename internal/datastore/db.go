package datastore

import (
	"fmt"

	"github.com/ali-shokoohi/micro-gopia/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbhost     = config.Confs.Postgres.Host
	dbport     = config.Confs.Postgres.Port
	dbuser     = config.Confs.Postgres.Username
	dbpassword = config.Confs.Postgres.Password
	dbname     = config.Confs.Postgres.DB
	sslmode    = config.Confs.Postgres.SslMode
)

// Database type
type Database struct{}

// GetDatabase () *gorm.DB {...} Return a valid database client
func (database *Database) GetDatabase() *gorm.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s TimeZone=Asia/Tehran",
		dbhost, dbport, dbuser, dbpassword, dbname, sslmode)

	//db, err := sql.Open("postgres", psqlInfo)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to database!")

	return db
}
