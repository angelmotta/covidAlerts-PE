package driver

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
)

type DB struct {
	SQL *sql.DB
	//Mongo *mgo.database
}

// Pointer to DB struct
var dbConn = &DB{}

func ConnectSQL(host, port, user, pass, dbname string) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	dbConn.SQL = db
	return dbConn, err
}