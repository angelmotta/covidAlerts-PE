package deceasedCases

import (
	"database/sql"
	"github.com/angelmotta/covidAlerts-PE/source/model"
	"github.com/angelmotta/covidAlerts-PE/source/repository"
	"log"
	"time"
)

// pgsqlDeceasedCasesRepo implements the interface 'repository.DeceasedCasesRepo'
type pgsqlDeceasedCasesRepo struct {
	Conn *sql.DB
}

// Return interface implementation
func NewSQLDeceasedCasesRepo(Conn *sql.DB) repository.DeceasedCasesRepo {
	return &pgsqlDeceasedCasesRepo{Conn}	// return interface
}

func (pgRepo *pgsqlDeceasedCasesRepo) Create(report *model.DeceasedReport) (int, error) {
	// Prepare statement
	stmtSQL := "INSERT INTO dailydeceased (deceasedcases_date, newdeceased_amount, totaldeceased, tsrecord) VALUES ($1, $2, $3, $4)"
	stmt, err := pgRepo.Conn.Prepare(stmtSQL)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	// Execute Sql statement
	dateTime := time.Now()
	tsRecord := dateTime.Format("2006-01-02 15:04:05")
	res, err := stmt.Exec(report.Date, report.NewDeceased, report.TotalDeceased, tsRecord)
	if err != nil {
		log.Println("SQL INSERT Execution failed, ", err)
		log.Printf("SQL values were, deceasedcases_date: %v, newdeceased_amount: %v, totaldeceased: %v \n", report.Date, report.NewDeceased, report.TotalDeceased)
		return 0, err
	}

	// Validate executed statement
	log.Println("SQL INSERT into dailydeceased table, Successfully Executed!")
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	log.Printf("DB Metada: #Rows affected = %d\n", rowCnt)

	return int(rowCnt), nil
}