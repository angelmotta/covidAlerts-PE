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
		repo: newCases.NewSQLNewCasesRepo(db.SQL),
	}
}

// Create daily newCases record
func (newCases *newCasesRepo) Create() error {
	fileNameCases := "dataFiles/positivos_covid_3_2_2021.csv"
	reportNewCases  := getReportCases(fileNameCases)	// Read and process CSV
	err := newCases.repo.Create(&reportNewCases)	// insert into DB (using Interface)
	if err != nil {
		return err
	}
	return nil
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
func (deceasedCases *deceasedCasesRepo) Create() error {
	fileNameDeceased := "dataFiles/fallecidos_covid_3_2_2021.csv"
	reportNewDeceased  := getReportDeceased(fileNameDeceased)
	err := deceasedCases.repo.Create(&reportNewDeceased)	// insert into DB (using Interface)
	if err != nil {
		return err
	}
	return nil
}

// New Post Tweet
func NewPostTweet(config *util.Config, numNewCases, numDeceased int) (int, error) {
	// Config Post Request
	configTwitter := oauth1.NewConfig(config.TApiKey, config.TApiSecretKey)
	token := oauth1.NewToken(config.TAccessToken, config.TAccessTokenSecret)

	httpClient := configTwitter.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// Post New Tweet
	// TODO: get message
	msgTweet := "Hola mundo: testing :)"
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