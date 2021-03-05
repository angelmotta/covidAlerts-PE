package main

import (
	"AlertsProject/source/model"
	"encoding/csv"
	"io"
	"log"
	"os"
)

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

func getReportCases(fileName string) model.CasesReport {
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

	// vars for Struct attributes
	newCasesLastDay := 0
	totalCases := 0
	casesByDept := make(map[string]int)

	// iterate over the entire file
	for {
		record, err := csvReader.Read() // get a record string[]
		if err == io.EOF { break }
		checkForError(err)

		// Count total covid cases from the beginning of the pandemic
		totalCases++

		// Looking for new cases in last day
		if record[0] == record[8] {  // Compare 'FECHA_CORTE' with 'FECHA_RESULTADO'
			// new case
			newCasesLastDay++
			// Get cases by 'DEPARTAMENTO'
			dept := record[2]
			value, isPresent := casesByDept[dept]
			if isPresent {
				casesByDept[dept]++
			} else {
				value = 1
				casesByDept[dept] = value
			}
		}
	}

	myNewReport := model.CasesReport{Date: lastDay, NewCases:newCasesLastDay, TotalCases: totalCases, NewCasesByDept: casesByDept}
	return myNewReport
}

func getReportDeceased(fileName string) model.DeceasedReport {
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
	deceaseByDept := make(map[string]int)
	for {
		record, err := csvReader.Read() // get a record string[]
		if err == io.EOF { break }
		checkForError(err)

		// Count total deceases from the beginning of pandemic
		totalDeceased++

		// Looking for deceases in the last day ('FECHA_CORTE' field)
		if record[0] == record[2] {  // Compare 'FECHA_CORTE' with 'FECHA_FALLECIMIENTO'
			deceasesLastDay++
			// Get decease by department
			dept := record[6]
			val, isPresent := deceaseByDept[dept]
			if isPresent {
				deceaseByDept[dept]++
			} else {
				val = 1
				deceaseByDept[dept] = val
			}
		}
	}

	myReportDeceases := model.DeceasedReport{Date: lastDay, Deceased: deceasesLastDay, TotalDeceased: totalDeceased, DeceasesByDept: deceaseByDept}
	return myReportDeceases
}


func main() {
	// process CSV File
	fileNameCases := "dataFiles/positivos_covid_3_2_2021.csv"
	fileNameDeceased := "dataFiles/fallecidos_covid_3_2_2021.csv"

	// Parse and receive a report
	newReportCases := getReportCases(fileNameCases)
	newReportDeceases := getReportDeceased(fileNameDeceased)

	newReportCases.Display()
	newReportDeceases.Display()
}
