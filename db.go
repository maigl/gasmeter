package main

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// todo use flag/env params
var (
	host       = "192.168.0.42"
	port       = 49155
	user       = "maigl"
	password   = "dreggn"
	dbname     = "ts"
	schemaName = "gasmeter"
	table      = "impulses"
)

var db *gorm.DB

func connectDB() (func(), error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error

	db, err = gorm.Open(postgres.Open(psqlconn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: schemaName + ".",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	return func() {
		db, err := db.DB()
		if err != nil {
			fmt.Println("error closing db connection: ", err)
		}
		db.Close()
	}, nil
}

func insertImpulseIntoDB(impulse *Impulse) error {
	db.AutoMigrate(&Impulse{})
	return db.Create(impulse).Error
}

func lastValueFromDB() (Impulse, error) {
	i := Impulse{}
	tx := db.Last(&i)
	if tx.Error != nil {
		insertImpulseIntoDB(&Impulse{Timestamp: time.Now(), ValueInM3: 0, Comment: "initial value"})
		return lastValueFromDB()
	}
	return i, tx.Error
}
