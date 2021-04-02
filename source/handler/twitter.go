package handler

import (
	"fmt"
	"log"
)

func GetTweetMsg(dateNewCases string, numNewCases int, dateDeceased string, numDeceased int) (string, bool) {
	// Check base case
	if dateNewCases == "" && dateDeceased == "" {
		log.Println("'dateNewCases' and 'dateDeceased' are blank values")
		return "", false
	}

	// Prepare msg
	var tweetMsg string
	if dateNewCases == dateDeceased { // Most probable case
		tweetMsg = fmt.Sprintf("MINSA publica hoy los datos del %v\n%v nuevos casos\n%v fallecidos", dateNewCases, numNewCases, numDeceased)
	} else {
		log.Println("Different dates detected for Tweet msg, please check it out")
		log.Printf("dateNewCases: %v, dateDeceased: %v\n", dateNewCases, dateDeceased)
		return "", false
	}
	return tweetMsg, true
}
