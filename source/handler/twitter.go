package handler

import (
	"fmt"
	"log"
)

func GetTweetMsg(dateNewCases string, numNewCases int, dateDeceased string, numDeceased int) ([]string, bool) {
	fmt.Printf("\n**** getTweetMsg ***\n")
	tweetList := make([]string,0)
	// Check base case
	if dateNewCases == "" && dateDeceased == "" {
		log.Println("'dateNewCases' and 'dateDeceased' are blank values")
		return tweetList, false // slice with zero elements
	}

	// Prepare Tweet
	var tweetMsg string
	if dateNewCases == dateDeceased { // Ideal case: 1 Tweet
		tweetMsg = fmt.Sprintf("📊 MINSA publica hoy los datos correspondientes al: %v\n\U0001F9A0 %v nuevos casos\n\U0001FAA6  %v fallecidos", dateNewCases, numNewCases, numDeceased)
		tweetList = append(tweetList, tweetMsg)
	} else { // 2 Tweets
		log.Println("Different dates detected for New_Positive_Cases and Deceased_Cases")
		log.Printf("dateNewCases: %v, dateDeceased: %v\n", dateNewCases, dateDeceased)
		// Prepare Tweet for New Positive Cases
		tweetMsg = fmt.Sprintf("📊 MINSA publica hoy los datos correspondientes al %v\n\U0001F9A0 %v nuevos casos", dateNewCases, numNewCases)
		tweetList = append(tweetList, tweetMsg)
		// Prepare Tweet for Deceased Cases
		tweetMsg = fmt.Sprintf("📊 MINSA publica hoy los datos correspondientes al %v\n\U0001FAA6 %v fallecidos", dateDeceased, numDeceased)
		tweetList = append(tweetList, tweetMsg)
	}
	return tweetList, true
}
