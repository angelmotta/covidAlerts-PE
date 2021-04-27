package handler

import (
	"errors"
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
	fmt.Println("**** Create daily New cases record ****")
	// Read CSV and return a report struct
	reportNewCases  := getReportCases(filePathPositive, "")
	// Insert into DB (using Interface)
	if reportNewCases.NewCases != 0 {
		_, err = newCases.repo.Create(&reportNewCases)
		if err != nil {	// if SQL insertion fail
			return dateCases, dailyCases, err	// return null values and error value
		}
		dateCases = reportNewCases.Date
		dailyCases = reportNewCases.NewCases
	} else {
		log.Println("Data in CSV File inconsistent. No insertion executed")
		log.Println("File is supposed updated to date: ", reportNewCases.Date, " but...")
		log.Println("New cases value in that date is: ", reportNewCases.NewCases)
		// Return null values and error
		return 	dateCases, dailyCases, errors.New("zero value new positive cases")
	}
	return dateCases, dailyCases, nil
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
	fmt.Println("**** Create daily Deceased cases record ****")
	// Read CSV and return a report
	reportNewDeceased  := getReportDeceased(filePathDeceased, "")
	// Insert into DB (using Interface)
	if reportNewDeceased.NewDeceased != 0 {
		_, err = deceasedCases.repo.Create(&reportNewDeceased)
		if err != nil {
			return dateDeceased, numDeceased, err // return null values and error value
		}
		dateDeceased = reportNewDeceased.Date
		numDeceased = reportNewDeceased.NewDeceased
	} else {
		log.Println("Data in CSV File inconsistent. No insertion executed")
		log.Println("File is supposed updated to date: ", reportNewDeceased.Date, " but...")
		log.Println("New cases value in that date is: ", reportNewDeceased.NewDeceased)
		return dateDeceased, numDeceased, errors.New("zero value deceased")
	}
	return dateDeceased, numDeceased, nil
}

// New Post Tweet
func PostTweet(config *util.Config, listTweets []string) error {
	fmt.Println("**** Post Tweet ****")
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
	log.Println("HTTP Twitter Credential StatusCode:", httpResCode.StatusCode)
	log.Println("Twitter User's Account: ", user.Name)

	// Post tweet
	err = sendTweetMsg(client, listTweets)		// Send Tweets
	if err != nil {
		log.Println("sendTweetMsg() error:", err)
	}
	return nil
}