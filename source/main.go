package main

import (
	"fmt"
	"github.com/angelmotta/covidAlerts-PE/source/driver"
	"github.com/angelmotta/covidAlerts-PE/source/handler"
	"github.com/angelmotta/covidAlerts-PE/source/util"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
)

func main() {
	// Load Config values
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Can not load configuration", err)
	}

	dbConn, err := driver.ConnectSQL(config.DBHost, config.DBPort, config.DBUser, config.DBPass, config.DBName)
	if err != nil {
		log.Fatal("DB connection error", err)
	}
	fmt.Println("Successfully connected to DB")

	newCasesHandler := handler.NewCasesHandler(dbConn)
	err = newCasesHandler.Create()	// get-read csv and insert into DB
	if err != nil {
		log.Println("Insertion newCases record failed, ", err)
	}

	deceasedCasesHandler := handler.NewDeceasedCasesHandler(dbConn)
	err = deceasedCasesHandler.Create()	// get-read csv and insert into DB
	if err != nil {
		log.Println("Insertion deceasedCases record failed, ", err)
	}
}
