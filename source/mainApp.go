package main

import (
	"fmt"
	"github.com/angelmotta/covidAlerts-PE/source/driver"
	//"database/sql"
	//"fmt"
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
	newCasesHandler.Create()	// get-read csv and insert into DB

	// Process CSV File
	/*
	fileNameCases := "dataFiles/positivos_covid_3_2_2021.csv"
	fileNameDeceased := "dataFiles/fallecidos_covid_3_2_2021.csv"

	newReportCases := handler.GetReportCases(fileNameCases)
	newReportDeceases := handler.GetReportDeceased(fileNameDeceased)

	newReportCases.Display()
	newReportDeceases.Display()

	// Db Connection
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.DBHost, config.DBPort, config.DBUser, config.DBPass, config.DBName)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	*/


}
