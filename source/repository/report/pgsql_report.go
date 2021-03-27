package report

import (
	"database/sql"
	"github.com/angelmotta/covidAlerts-PE/source/model"
	"github.com/angelmotta/covidAlerts-PE/source/repository"
)

// pgsqlReportRepository will implement the interface 'repository.Repository'
type pgsqlReportRepository struct {
	Conn *sql.DB
}

// Return interface implementation
func NewPgSqlReportRepository (Conn *sql.DB) repository.Repository {
	return &pgsqlReportRepository{Conn}
}

// pgsqlReportRepository implements 'pgsqlReportRepository' interface
func (pgRepo *pgsqlReportRepository) Create(report *model.CasesReport) error {
	// TODO
	//sqlStatement := `INSERT INTO dailycases (date, newcases) VALUES ($1, $2)`
	return nil
}