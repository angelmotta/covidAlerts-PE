package model

import "fmt"

type DeceasedReport struct {
	Date           string
	NewDeceased    int
	TotalDeceased  int
	DeceasesByDept map[string]int
}

func (report *DeceasedReport) Display() {
	// Display Results
	fmt.Println("\n*** Report NewDeceased ***")
	fmt.Printf("Número de fallecidos (el día %s): %d personas\n", report.Date, report.NewDeceased)
	fmt.Println("Fallecidos por departamento", report.DeceasesByDept)
	fmt.Printf("Total de fallecidos a la fecha: %d personas\n", report.TotalDeceased)
}
