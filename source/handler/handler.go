package handler

import (
	"fmt"
	"github.com/angelmotta/covidAlerts-PE/source/driver"
	"github.com/angelmotta/covidAlerts-PE/source/repository"
	"github.com/angelmotta/covidAlerts-PE/source/repository/deceasedCases"
	"github.com/angelmotta/covidAlerts-PE/source/repository/newCases"
	"github.com/angelmotta/covidAlerts-PE/source/util"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"log"
)

// New Daily Cases
type newCasesRepo struct {
	repo repository.NewCasesRepo // interface
}

// Return struct 'newCasesRepo' with repository Interface
func NewCasesHandler(db *driver.DB) *newCasesRepo {
	return &newCasesRepo{
		repo: newCases.NewSQLNewCasesRepo(db.SQL),	// 'newCases' (interface implementation)
	}
}

// Create daily newCases record
func (newCases *newCasesRepo) Create(filePathPositive string) (dateCases string, dailyCases int, err error) {
	// Read CSV and return a report struct
	reportNewCases  := getReportCases(filePathPositive, "")
	// Insert into DB (using Interface)
	_, err = newCases.repo.Create(&reportNewCases)
	if err != nil {	// if SQL insertion fail
		return 		// return null date, null dailyCases and error value
	}
	dateCases = reportNewCases.Date
	dailyCases = reportNewCases.NewCases
	return
}

// New Deceased Cases
type deceasedCasesRepo struct {
	repo repository.DeceasedCasesRepo // interface
}

// Return struct 'newCasesRepo' with repository Interface
func NewDeceasedCasesHandler(db *driver.DB) *deceasedCasesRepo {
	return &deceasedCasesRepo{
		repo: deceasedCases.NewSQLDeceasedCasesRepo(db.SQL),
	}
}

// Create daily deceased record
func (deceasedCases *deceasedCasesRepo) Create(filePathDeceased string) (dateDeceased string, numDeceased int, err error) {
	// Read CSV and return a report
	reportNewDeceased  := getReportDeceased(filePathDeceased, "")
	// Insert into DB (using Interface)
	_, err = deceasedCases.repo.Create(&reportNewDeceased)
	if err != nil {
		return
	}
	dateDeceased = reportNewDeceased.Date
	numDeceased = reportNewDeceased.NewDeceased
	return
}

func sendTweetMsg(client *twitter.Client, listTweets []string) error {
	for _, tweet := range listTweets {		// Post Tweet
		// Tweeting
		fmt.Println("Sending the following tweet")
		fmt.Println(tweet)
		fmt.Println()
		// Send Tweet
		_, respHttp, err := client.Statuses.Update(tweet, nil)
		if err != nil {
			log.Println("client.Statuses.Update(tweetMsg) error:", err)
			log.Println("client.Statuses.Update(tweetMsg) response HTTP:", respHttp.StatusCode)
			return err
		}
		fmt.Println("Tweet successfully sent")
	}
	return nil
}

// New Post Tweet
func PostTweet(config *util.Config, listTweets []string) error {
	// Config Post Request
	configTwitter := oauth1.NewConfig(config.TApiKey, config.TApiSecretKey)
	token := oauth1.NewToken(config.TAccessToken, config.TAccessTokenSecret)

	httpClient := configTwitter.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}
	user, httpResCode, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		log.Println("VerifyCredentials() Error:", err)
		log.Println("HTTP Response Code:", httpResCode.StatusCode)
		return err
	}
	fmt.Println("HTTP Twitter Credential StatusCode:", httpResCode.StatusCode)
	fmt.Println("Twitter User's Account: ", user.Name)

	// Post tweet
	err = sendTweetMsg(client, listTweets)		// Send Tweets
	if err != nil {
		log.Println("sendTweetMsg() error:", err)
	}
	return nil
}