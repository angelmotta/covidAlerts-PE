package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

type ReportCases struct {
	date		string
	newCases 	int
	totalCases	int
	//newCasesByDept map
}

func (report *ReportCases) display() {
	// Display Results
	fmt.Printf("Nuevos casos (registrados el %s): %d contagiados\n", report.date, report.newCases)
	fmt.Printf("Total de casos a la fecha: %d contagiados \n", report.totalCases)
}

type ReportDeceased struct {
	date          string
	deceased      int
	totalDeceased int
	//deceasesByDept map
}

func (report *ReportDeceased) display() {
	// Display Results
	fmt.Printf("Número de fallecidos (el día %s): %d personas\n", report.date, report.deceased)
	fmt.Printf("Total de fallecidos a la fecha: %d personas\n", report.totalDeceased)
}

func checkForError(e error) {
	if e != nil {
		log.Fatalln("Error reading csv file: ", e)
	}
}

func getLastDay(fileName string) string {
	// Try Open file
	csvFile, err := os.Open(fileName)
	checkForError(err)
	defer csvFile.Close()

	// setup csvReader
	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = ';'

	// discard first line (column headers)
	_, err = csvReader.Read()
	checkForError(err)

	// get from the first line the 'FECHA_CORTE' field
	record, err := csvReader.Read()
	checkForError(err)
	return record[0]
}

func getReportCases(fileName string) ReportCases {
	// get lastDay based on 'FECHA_CORTE' field
	lastDay := getLastDay(fileName)

	// Try Open file
	csvFile, err := os.Open(fileName)
	checkForError(err)
	defer csvFile.Close()

	// Setup a csv reader
	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = ';'

	// discard first line
	_, err = csvReader.Read()
	checkForError(err)

	// iterate over the entire file
	newCasesLastDay := 0
	totalCases := 0
	for {
		record, err := csvReader.Read() // get a record string[]
		if err == io.EOF { break }
		checkForError(err)

		// Count total covid cases from the beginning of the pandemic
		totalCases++

		// CSV Sanity check
		if lastDay != record[0] {	// assure all records has the same 'FECHA_CORTE'
			fmt.Printf("Record #%d UUID:%s has a different 'FECHA_CORTE' attribute\n", totalCases, record[1])
		}

		// Looking for cases in last day
		if record[0] == record[8] {  // Compare 'FECHA_CORTE' with 'FECHA_RESULTADO'
			//fmt.Printf("Caso en el ultimo día\n")
			newCasesLastDay++
		}
	}

	myNewReport := ReportCases{lastDay, newCasesLastDay, totalCases}
	return myNewReport
}

func getReportDeceased(fileName string) ReportDeceased {
	// get lastDay based on 'FECHA_CORTE' field
	lastDay := getLastDay(fileName)

	// Try Open file
	csvFile, err := os.Open(fileName)
	checkForError(err)
	defer csvFile.Close()

	// Setup a csv reader
	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = ';'

	// discard first line of column headers
	_, err = csvReader.Read()
	checkForError(err)

	// iterate over the csv file
	totalDeceased := 0
	deceasesLastDay := 0

	for {
		record, err := csvReader.Read() // get a record string[]
		if err == io.EOF { break }
		checkForError(err)

		// Count total deceases from the beginning of pandemic
		totalDeceased++

		// Looking for more recent deceases
		if record[0] == record[2] {  // Compare 'FECHA_CORTE' with 'FECHA_FALLECIMIENTO'
			//fmt.Printf("Caso en el ultimo día\n")
			deceasesLastDay++
		}
	}

	// Dummy data
	myReportDeceases := ReportDeceased{lastDay, deceasesLastDay, totalDeceased}
	return myReportDeceases
}


func main() {
	// process CSV File
	fileNameCases := "dataFiles/positivos_covid_3_2_2021.csv"
	fileNameDeceased := "dataFiles/fallecidos_covid_3_2_2021.csv"

	// Parse and receive a report
	newReportCases := getReportCases(fileNameCases)
	newReportDeceases := getReportDeceased(fileNameDeceased)

	newReportCases.display()
	newReportDeceases.display()
}
