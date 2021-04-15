package handler

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"log"
)

func GenerateTweetMsg(dateNewCases string, numNewCases int, dateDeceased string, numDeceased int) ([]string, bool) {
	fmt.Printf("\n**** Generate Tweet Messages ****\n")
	tweetList := make([]string,0)
	// Check base case
	if dateNewCases == "" && dateDeceased == "" {
		log.Println("'dateNewCases' and 'dateDeceased' are blank values")
		return tweetList, false // return slice with zero elements
	}

	// Prepare Tweet
	var tweetMsg string
	if dateNewCases == dateDeceased { // Ideal case: 1 Tweet with full information
		tweetMsg = fmt.Sprintf("📊 MINSA publica hoy información correspondiente al %v\n\U0001F9A0 %v nuevos casos\n\U0001FAA6 %v fallecidos", dateNewCases, numNewCases, numDeceased)
		tweetList = append(tweetList, tweetMsg)
	} else { // 1 or 2 Tweets each one with specific information (num new cases or num deceased)
		log.Println("Different dates detected for New_Positive_Cases and Deceased_Cases")
		log.Printf("dateNewCases: '%v', dateDeceased: '%v'\n", dateNewCases, dateDeceased)
		// Prepare Tweet for New Positive Cases
		if dateNewCases != "" {
			tweetMsg = fmt.Sprintf("📊 MINSA publica hoy información correspondiente al %v\n\U0001F9A0 %v nuevos casos", dateNewCases, numNewCases)
			tweetList = append(tweetList, tweetMsg)
		}
		if dateDeceased != "" {
			// Prepare Tweet for Deceased Cases
			tweetMsg = fmt.Sprintf("📊 MINSA publica hoy información correspondiente al %v\n\U0001FAA6 %v fallecidos", dateDeceased, numDeceased)
			tweetList = append(tweetList, tweetMsg)
		}
	}
	return tweetList, true
}

func sendTweetMsg(client *twitter.Client, listTweets []string) error {
	for _, tweet := range listTweets {		// Post Tweet
		// Tweeting
		log.Println("Sending the following tweet:")
		fmt.Println(tweet)
		// Send Tweet
		_, respHttp, err := client.Statuses.Update(tweet, nil)
		if err != nil {
			log.Println("client.Statuses.Update(tweetMsg) error:", err)
			log.Println("client.Statuses.Update(tweetMsg) response HTTP:", respHttp.StatusCode)
			return err
		}
		log.Println("Tweet successfully sent")
	}
	return nil
}