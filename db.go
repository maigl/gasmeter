package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// todo use flag/env params
var (
	host     = "192.168.0.42"
	port     = 49155
	user     = "maigl"
	password = "dreggn"
	dbname   = "ts"
	schema   = "gasmeter"
	table    = "impulses"
	db       *sql.DB
)

func connectDB() (*sql.DB, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error

	db, err = sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}

	createSchemaStmt := `create schema if not exists ` + schema

	_, err = db.Exec(createSchemaStmt)
	if err != nil {
		db.Close()
		return nil, err
	}

	_, err = db.Exec(`set timezone = 'Europe/Berlin'`)
	if err != nil {
		db.Close()
		return nil, err
	}

	createTable := `create table if not exists ` + table + ` (value_in_m3 real, comment varchar(255), time timestamptz not null default current_timestamp primary key)`

	_, err = db.Exec(createTable)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func insertImpulseIntoDB(value float64) error {
	return insertValueIntoDB(value, "impulse")
}

func updateValueInDB(value float64) error {
	return insertValueIntoDB(value, "manual update")
}

func insertValueIntoDB(value float64, comment string) error {
	// check length of comment
	if len(comment) > 255 {
		comment = comment[:255]
	}

	insertStmt := `insert into ` + table + ` values ($1, $2)`

	_, err := db.Exec(insertStmt, value, comment)
	if err != nil {
		return fmt.Errorf("error inserting value into db: %w", err)
	}

	return nil
}

func lastValueFromDB() (time.Time, float64, error) {
	rows, err := db.Query(`select value_in_m3, time from ` + table + ` order by time desc limit 1`)
	if err != nil {
		return time.Time{}, 0, err
	}

	defer rows.Close()

	if !rows.Next() {
		return time.Time{}, 0, nil
	}

	rows.Next()
	var ts time.Time
	var value float64

	err = rows.Scan(&value, &ts)
	if err != nil {
		return time.Time{}, 0, err
	}
	return ts, value, nil
}
