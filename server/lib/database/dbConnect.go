package database

import (
	"fmt"
	"log"

	"gRPC-PostgreSQL-REST/server/model"

	"github.com/jinzhu/gorm"

	//pq is a blank import for database driver
	_ "github.com/lib/pq" //blank import contains initialization code
)

//DBRepo has an gorm object for connecting to DB
type DBRepo struct {
	GormDB *gorm.DB
}

//DBConnect Method to connect to Db
func (dc *DBRepo) DBConnect(dbConfig model.DBConfig) error {
	var err error
	// Format DB configs to required format for connecting to DB
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.Port)
	dc.GormDB, err = gorm.Open("postgres", dbinfo)
	if err != nil {
		log.Printf("Unable to connect DB %v", err)
		return err
	}
	log.Printf("Postgres started at %s PORT", "5432")
	return err
}
