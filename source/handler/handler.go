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
func (newCases *newCasesRepo) Create() (dateCases string, dailyCases int, err error) {
	fileNameCases := "dataFiles/positivos_covid_3_2_2021.csv"
	reportNewCases  := getReportCases(fileNameCases)		// Read CSV and return a report
	_, err = newCases.repo.Create(&reportNewCases)	// insert into DB (using Interface)
	if err != nil {
		return
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
	fileNameDeceased := "dataFiles/positivos_covid_1_4_2021.csv"
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
func NewPostTweet(config *util.Config, dateNewCases string, numNewCases int, dateDeceased string, numDeceased int) (int, error) {
	// Config Post Request
	configTwitter := oauth1.NewConfig(config.TApiKey, config.TApiSecretKey)
	token := oauth1.NewToken(config.TAccessToken, config.TAccessTokenSecret)

	httpClient := configTwitter.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// Post New Tweet
	msgTweet := fmt.Sprintf("MINSA publica hoy los datos del %v\nNuevos casos: %v\nNúmero de fallecidos:%v", dateNewCases, numNewCases, numDeceased)
	fmt.Println(msgTweet)
	// Post tweet
	tweet, resp, err := client.Statuses.Update(msgTweet, nil)
	if err != nil {
		log.Println("Remote Twitter API Server not responding")
		return 0, err
	}

	codeResp := resp.StatusCode
	if codeResp != 200 {
		log.Printf("Response:\n%v\n",resp)
		return 1, nil
	}

	log.Printf("Posted Tweet\n%v\n", tweet)
	return 0, nil
}