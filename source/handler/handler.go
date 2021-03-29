package handler

import (
	"github.com/angelmotta/covidAlerts-PE/source/driver"
	"github.com/angelmotta/covidAlerts-PE/source/repository"
	"github.com/angelmotta/covidAlerts-PE/source/repository/deceasedCases"
	"github.com/angelmotta/covidAlerts-PE/source/repository/newCases"
)

// New Daily Cases
type NewCasesRepo struct {
	repo repository.NewCasesRepo // interface
}

// Return struct 'NewCasesRepo' with repository Interface
func NewCasesHandler(db *driver.DB) *NewCasesRepo {
	return &NewCasesRepo{
		repo: newCases.NewSQLNewCasesRepo(db.SQL),
	}
}

func (newCases *NewCasesRepo) Create() error {
	fileNameCases := "dataFiles/positivos_covid_3_2_2021.csv"
	reportNewCases  := getReportCases(fileNameCases)
	err := newCases.repo.Create(&reportNewCases)	// insert into DB (using Interface)
	if err != nil {
		return err
	}
	return nil
}

// New Deceased Cases
type DeceasedCasesRepo struct {
	repo repository.DeceasedCasesRepo // interface
}

// Return struct 'NewCasesRepo' with repository Interface
func NewDeceasedCasesHandler(db *driver.DB) *DeceasedCasesRepo {
	return &DeceasedCasesRepo{
		repo: deceasedCases.NewSQLDeceasedCasesRepo(db.SQL),
	}
}

func (deceasedCases *DeceasedCasesRepo) Create() error {
	fileNameDeceased := "dataFiles/fallecidos_covid_3_2_2021.csv"
	reportNewDeceased  := getReportDeceased(fileNameDeceased)
	err := deceasedCases.repo.Create(&reportNewDeceased)	// insert into DB (using Interface)
	if err != nil {
		return err
	}
	return nil
}
