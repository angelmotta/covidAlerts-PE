package report

import (
	"database/sql"
	"github.com/angelmotta/covidAlerts-PE/source/model"
	"github.com/angelmotta/covidAlerts-PE/source/repository"
	"log"
)

// pgsqlReportRepository will implement the interface 'repository.Repository'
type pgsqlReportRepository struct {
	Conn *sql.DB
}

// Return interface implementation
func NewPgSqlReportRepository (Conn *sql.DB) repository.Repository {
	return &pgsqlReportRepository{Conn}
}

// pgsqlReportRepository implements 'pgsqlReportRepository' Interface
func (pgRepo *pgsqlReportRepository) Create(report *model.CasesReport) error {
	// Prepare statement
	stmt, err := pgRepo.Conn.Prepare("INSERT INTO dailycases (newcases_date, newcases) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute Sql statement
	res, err := stmt.Exec(report.Date, report.NewCases)	// 'Date' string, 'NewCases' int
	if err != nil {
		return err
	}

	// Validate executed statement
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("DBMetada: #Rows affected = %d\n", rowCnt)

	return nil
}