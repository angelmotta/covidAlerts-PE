package model

import "fmt"

// Daily Report
type CasesReport struct {
	Date           string
	NewCases       int
	TotalCases     int
	NewCasesByDept map[string]int
}

func (report *CasesReport) Display() {
	// Display Results
	fmt.Println("\n*** New Cases Report ***")
	fmt.Printf("Nuevos casos (registrados el %s): %d contagiados\n", report.Date, report.NewCases)
	fmt.Println("Casos por Departamento", report.NewCasesByDept)
	fmt.Printf("Total de casos a la fecha: %d contagiados \n\n", report.TotalCases)
}
