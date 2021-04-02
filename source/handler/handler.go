package handler

import (
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
func (newCases *newCasesRepo) Create() (dateCases string, dailyCases int, err error) {
	// TODO: Method HTTP to get CSV file
	fileNameCases := "dataFiles/positivos_covid_2_4_2021.csv"
	reportNewCases  := getReportCases(fileNameCases)		// Read CSV and return a report
	_, err = newCases.repo.Create(&reportNewCases)	// insert into DB (using Interface)
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
func (deceasedCases *deceasedCasesRepo) Create() (dateDeceased string, numDeceased int, err error) {
	fileNameDeceased := "dataFiles/fallecidos_covid_2_4_2021.csv"
	reportNewDeceased  := getReportDeceased(fileNameDeceased)	// Read CSV and return a report
	_, err = deceasedCases.repo.Create(&reportNewDeceased)	// insert into DB (using Interface)
	if err != nil {
		return
	}

	dateDeceased = reportNewDeceased.Date
	numDeceased = reportNewDeceased.NewDeceased
	return
}


// New Post Tweet
func NewPostTweet(config *util.Config, tweetMsg string) (int, error) {
	// Config Post Request
	configTwitter := oauth1.NewConfig(config.TApiKey, config.TApiSecretKey)
	token := oauth1.NewToken(config.TAccessToken, config.TAccessTokenSecret)

	httpClient := configTwitter.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// Post tweet
	tweet, resp, err := client.Statuses.Update(tweetMsg, nil)
	if err != nil {
		log.Println("Remote Twitter API Server not responding")
		return 0, err
	}

	codeResp := resp.StatusCode
	if codeResp != 200 {
		log.Printf("Response:\n%v\n",resp)
		return codeResp, nil
	}

	log.Printf("Posted Tweet\n%v\n", tweet)
	return codeResp, nil
}