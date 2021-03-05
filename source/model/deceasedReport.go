package model

import "fmt"

type DeceasedReport struct {
	Date          	string
	Deceased      	int
	TotalDeceased	int
	DeceasesByDept	map[string]int
}

func (report *DeceasedReport) Display() {
	// Display Results
	fmt.Println("*** Report Deceased ***")
	fmt.Printf("Número de fallecidos (el día %s): %d personas\n", report.Date, report.Deceased)
	fmt.Println("Fallecidos por departamento", report.DeceasesByDept)
	fmt.Printf("Total de fallecidos a la fecha: %d personas\n", report.TotalDeceased)
}
