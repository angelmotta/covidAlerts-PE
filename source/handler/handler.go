package handler

import (
	"github.com/angelmotta/covidAlerts-PE/source/driver"
	"github.com/angelmotta/covidAlerts-PE/source/repository"
	"github.com/angelmotta/covidAlerts-PE/source/repository/report"
)

type NewCasesRepo struct {
	repo repository.Repository	// interface
}

// Return struct 'NewCasesRepo' with repository Interface
func NewCasesHandler(db *driver.DB) *NewCasesRepo {
	return &NewCasesRepo{
		repo: report.NewPgSqlReportRepository(db.SQL),
	}
}

func (newCases *NewCasesRepo) Create() error {
	fileNameCases := "dataFiles/positivos_covid_3_2_2021.csv"
	reportNewCases  := getReportCases(fileNameCases)
	err := newCases.repo.Create(&reportNewCases)	// insert into DB (using Interface)
	if err != nil {
		//log.Fatal(err)
		return err
	}
	return nil
}
