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

	// Insert new daily cases
	newCasesHandler := handler.NewCasesHandler(dbConn)
	dateDailyCases, numNewCases, err := newCasesHandler.Create()	// get-read csv and insert into DB
	if err != nil {
		log.Println("No Insertion of new cases record, ", err)
	}

	// Insert new deceased cases
	deceasedCasesHandler := handler.NewDeceasedCasesHandler(dbConn)
	dateDeceased, numDeceased, err := deceasedCasesHandler.Create()	// get-read csv and insert into DB
	if err != nil {
		log.Println("No Insertion of new deceased Cases record, ", err)
	}

	// Twitter Post publication
	fmt.Printf("\nDatos oficiales:\n%v nuevos casos (%v)\n%v fallecidos (%v)\n",numNewCases, dateDailyCases, numDeceased, dateDeceased)

	// Read Tweets Msg
	listTweets, isOK := handler.GetTweetMsg(dateDailyCases, numNewCases, dateDeceased, numDeceased)
	if isOK != true {
		os.Exit(0)
	}
	// Send Tweets
	for _, tweet := range listTweets {
		fmt.Println(tweet)
		fmt.Println()
	}
	/*
	tweetMsg, isValid := handler.GetTweetMsg(dateDailyCases, numNewCases, dateDeceased, numDeceased)
	if isValid != true {
		os.Exit(0)
	}
	log.Printf("Prepared Tweet: \n%v\n", tweetMsg)

	// Post Tweet
	codResp, err := handler.NewPostTweet(&config, tweetMsg)
	if err != nil {
		log.Println("API not responding. Post tweet failed, ", err)
	}
	if codResp != 200 {
		log.Println("Post tweet failed, check response message from API Server")
	}
	log.Println("Post tweet Done!")
	 */
}
