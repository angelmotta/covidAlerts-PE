package newCases

import (
	"database/sql"
	"github.com/angelmotta/covidAlerts-PE/source/model"
	"github.com/angelmotta/covidAlerts-PE/source/repository"
	"log"
)

// pgsqlNewCasesRepo implements the interface 'repository.NewCasesRepo'
type pgsqlNewCasesRepo struct {
	Conn *sql.DB
}

// Return interface implementation
func NewSQLNewCasesRepo(Conn *sql.DB) repository.NewCasesRepo {
	return &pgsqlNewCasesRepo{Conn}	// return interface
}

// pgsqlNewCasesRepo implements 'pgsqlNewCasesRepo' Interface
func (pgRepo *pgsqlNewCasesRepo) Create(report *model.CasesReport) error {
	// Prepare statement
	stmt, err := pgRepo.Conn.Prepare("INSERT INTO dailycases (newcases_date, newcases_amount) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute Sql statement
	res, err := stmt.Exec(report.Date, report.NewCases)
	if err != nil {
		log.Println("SQL INSERT Execution failed, ", err)
		return err
	}

	// Validate executed statement
	log.Println("SQL INSERT into dailycases table, Successfully Executed!")
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("DBMetada: #Rows affected = %d\n", rowCnt)

	return nil
}