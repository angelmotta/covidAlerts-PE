package newCases

import (
	"database/sql"
	"github.com/angelmotta/covidAlerts-PE/source/model"
	"github.com/angelmotta/covidAlerts-PE/source/repository"
	"log"
	"time"
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
func (pgRepo *pgsqlNewCasesRepo) Create(report *model.CasesReport) (int, error) {
	// Prepare statement
	stmtSQL := "INSERT INTO dailycases (newcases_date, newcases_amount, totalcases, tsrecord) VALUES ($1, $2, $3, $4)"
	stmt, err := pgRepo.Conn.Prepare(stmtSQL)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	// Execute SQL statement
	// Get timestamp for new record
	dateTime := time.Now()
	tsRecord := dateTime.Format("2006-01-02 15:04:05")
	res, err := stmt.Exec(report.Date, report.NewCases, report.TotalCases, tsRecord) // INSERT
	if err != nil {
		log.Println("SQL INSERT Execution failed, ", err)
		log.Printf("SQL values were, newcases_date: %v, newcases_amount: %v, totalCases: %v \n", report.Date, report.NewCases, report.TotalCases)
		return 0, err
	}

	// Log executed statement
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	log.Println("SQL INSERT into dailycases table, Successfully Executed!")
	log.Printf("DB Metada: #Rows affected = %d\n", rowCnt)

	return int(rowCnt), nil
}