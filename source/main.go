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

	// Get Positive Cases CSV File
	positiveFilePath := config.DirPositiveFiles+"positivos_covid.csv"
	//positiveFilePath := config.DirPositiveFiles+"positivos_covid_3_4_2021.csv"
	err = handler.DownloadFile(config.UrlNewCases, positiveFilePath) 	// Download CSV File
	isPositiveFileOk := true

	if err != nil {
		log.Println("DownloadFile() positive cases error: ", err)
		isPositiveFileOk = false
	}
	if isPositiveFileOk {
		log.Println("Positive file stored OK in:", positiveFilePath)
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

	// Get Deceased Cases CSV File
	//deceasedFilePath := config.DirDeceasedFiles+"fallecidos_covid_3_4_2021.csv"
	deceasedFilePath := config.DirDeceasedFiles+"fallecidos_covid.csv"
	err = handler.DownloadFile(config.UrlDeceased, deceasedFilePath)	// Download CSV File
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
	
	// Twitter Post publication
	fmt.Printf("\nDatos oficiales obtenidos:\n%v nuevos casos (%v)\n%v fallecidos (%v)\n",numNewCases, dateDailyCases, numDeceased, dateDeceased)

	// Read Tweets Msg
	listTweets, isOK := handler.GetTweetMsg(dateDailyCases, numNewCases, dateDeceased, numDeceased)
	if isOK != true {
		log.Println("Not required creation of new Tweets")
		os.Exit(0)
	}
	// Send Tweets
	// Only print tweets to the console
	log.Println("Display new Tweets")
	for _, tweet := range listTweets {
		fmt.Println(tweet)
		fmt.Println()
	}

	// Post Tweet
	/*
	err = handler.PostTweet(&config, listTweets)
	if err != nil {
		log.Println("PostTweet() failed, ", err)
	}
	log.Println("Post tweet Done!")
	*/
}
