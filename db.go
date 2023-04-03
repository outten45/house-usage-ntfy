package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
)

var dbLog *sqlx.DB

var schema = `
  CREATE TABLE IF NOT EXISTS events (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		key text not null,
		time int,
		value numeric,
        last_updated_at int,
		UNIQUE(key)
  );
  `

var insertStmt = `
	INSERT INTO 
	  events (key, time, value, last_updated_at)
	VALUES
	  (:key, :time, :value, :last_updated_at)	
	
	`

func createLogDB(dbFile string) error {
	var err error
	dbConnStr := fmt.Sprintf("%s?_pragma=journal_mode(WAL)", dbFile)
	dbLog, err = sqlx.Open("sqlite", dbConnStr)
	if err != nil {
		log.Fatalf("db create error: %s\n", err)
		return err
	}

	_, err = dbLog.Exec(schema)
	if err != nil {
		log.Fatalf("Unable to create the schema: %s\n", err)
	}
	return nil
}

func saveEvent(key string, value float64) error {
	now := time.Now().Unix()
	e := struct {
		Key         string
		Time        int64
		Value       float64
		LastUpdated int64 `db:"last_updated_at"`
	}{
		Key:         key,
		Time:        now,
		Value:       value,
		LastUpdated: now,
	}

	tx := dbLog.MustBegin()
	_, err := tx.NamedExec(insertStmt, e)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}
	err = tx.Commit()

	return nil
}
