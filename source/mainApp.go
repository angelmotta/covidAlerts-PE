package main

import (
	"database/sql"
	"fmt"
	"github.com/angelmotta/covidAlerts-PE/source/usecase"
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

	// TODO: HTTP get CSV Files

	// Process CSV File
	fileNameCases := "dataFiles/positivos_covid_3_2_2021.csv"
	fileNameDeceased := "dataFiles/fallecidos_covid_3_2_2021.csv"

	newReportCases := usecase.GetReportCases(fileNameCases)
	newReportDeceases := usecase.GetReportDeceased(fileNameDeceased)

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

	fmt.Println("\nDB Successfully connected!\n")

}
