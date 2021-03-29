package handler

import (
	"github.com/angelmotta/covidAlerts-PE/source/driver"
	"github.com/angelmotta/covidAlerts-PE/source/repository"
	"github.com/angelmotta/covidAlerts-PE/source/repository/deceasedCases"
	"github.com/angelmotta/covidAlerts-PE/source/repository/newCases"
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
	reportNewCases  := getReportCases(fileNameCases)
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
