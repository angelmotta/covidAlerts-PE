package main

import (
	"fmt"
	"github.com/angelmotta/covidAlerts-PE/source/driver"
	"github.com/angelmotta/covidAlerts-PE/source/handler"
	"github.com/angelmotta/covidAlerts-PE/source/util"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"os"
)

func main() {
	// Load Config values
	fmt.Println("\n**** START EXECUTION ****")
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Can not load configuration", err)
	}

	// Prepare DB Connection
	dbConn, err := driver.ConnectSQL(config.DBHost, config.DBPort, config.DBUser, config.DBPass, config.DBName)
	if err != nil {
		log.Fatal("DB connection error", err)
	}
	log.Println("Successfully connected to DB")

	// Download Positive Cases CSV File
	positiveFilePath := config.DirPositiveFiles+"positivos_covid.csv"
	// Test using local file
	//positiveFilePath := config.DirPositiveFiles+"15_04_2021_1300_positivos_covid.csv"
	err = handler.DownloadFile(config.UrlNewCases, positiveFilePath) 	// Get CSV File
	isPositiveFileOk := true

	if err != nil {
		log.Println("DownloadFile() positive cases error: ", err)
		isPositiveFileOk = false
	}
	if isPositiveFileOk {
		log.Println("Positive csv file is stored OK in:", positiveFilePath)
	}

	// Insert new daily cases
	var dateDailyCases string
	var numNewCases int
	if isPositiveFileOk {
		newCasesHandler := handler.NewCasesHandler(dbConn)
		dateDailyCases, numNewCases, err = newCasesHandler.Create(positiveFilePath)	// Read CSV and insert into DB
		if err != nil {
			log.Println("No Insertion of new cases record, ", err)
		}
	}

	// Download Deceased Cases CSV File
	deceasedFilePath := config.DirDeceasedFiles+"fallecidos_covid.csv"
	// Test using local file
	//deceasedFilePath := config.DirDeceasedFiles+"15_04_2021_1300_fallecidos_covid.csv"
	err = handler.DownloadFile(config.UrlDeceased, deceasedFilePath) // Get CSV file
	isDeceasedFileOk := true

	if err != nil {
		log.Println("DownloadFile() deceased error: ", err)
		isDeceasedFileOk = false
	}
	if isDeceasedFileOk {
		log.Println("Deceased File stored OK in:", deceasedFilePath)
	}

	// Insert new deceased cases
	var dateDeceased string
	var numDeceased int
	if isDeceasedFileOk {
		deceasedCasesHandler := handler.NewDeceasedCasesHandler(dbConn)
		dateDeceased, numDeceased, err = deceasedCasesHandler.Create(deceasedFilePath)	// Read CSV and insert into DB
		if err != nil {
			log.Println("No Insertion of new deceased Cases record, ", err)
		}
	}

	// Display information collected from files
	fmt.Printf("\nDatos obtenidos:\n%v nuevos casos correspondiente al '%v'\n%v fallecidos correspondiente al '%v'\n",numNewCases, dateDailyCases, numDeceased, dateDeceased)

	// Twitter Post publication
	// Generate Tweets Messages
	listTweets, isNewInfo := handler.GenerateTweetMsg(dateDailyCases, numNewCases, dateDeceased, numDeceased)
	if isNewInfo != true {
		log.Println("Not required creation of new Tweets")
		log.Println("**** END EXECUTION ***")
		os.Exit(0)
	}

	// Display tweets
	log.Println("Display Tweet message:")
	for _, tweet := range listTweets {
		fmt.Println(tweet)
		fmt.Println()
	}

	// Send Tweet
	/*
	err = handler.PostTweet(&config, listTweets)
	if err != nil {
		log.Println("PostTweet() failed, ", err)
	}
	log.Println("**** END EXECUTION ***")
	*/
}
