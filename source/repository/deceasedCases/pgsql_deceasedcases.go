package deceasedCases

import (
	"database/sql"
	"github.com/angelmotta/covidAlerts-PE/source/model"
	"github.com/angelmotta/covidAlerts-PE/source/repository"
	"log"
)

// pgsqlDeceasedCasesRepo implements the interface 'repository.DeceasedCasesRepo'
type pgsqlDeceasedCasesRepo struct {
	Conn *sql.DB
}

// Return interface implementation
func NewSQLDeceasedCasesRepo(Conn *sql.DB) repository.DeceasedCasesRepo {
	return &pgsqlDeceasedCasesRepo{Conn}	// return interface
}

func (pgRepo *pgsqlDeceasedCasesRepo) Create(report *model.DeceasedReport) error {
	// Prepare statement
	stmt, err := pgRepo.Conn.Prepare("INSERT INTO dailydeceased (deceasedcases_date, newdeceased_amount, totaldeceased) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute Sql statement
	res, err := stmt.Exec(report.Date, report.NewDeceased, report.TotalDeceased)
	if err != nil {
		log.Println("SQL INSERT Execution failed, ", err)
		return err
	}

	// Validate executed statement
	log.Println("SQL INSERT into dailydeceased table, Successfully Executed!")
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("DBMetada: #Rows affected = %d\n", rowCnt)

	return nil
}