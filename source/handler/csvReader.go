package handler

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/angelmotta/covidAlerts-PE/source/model"
	"io"
	"log"
	"os"
	"time"
)

func checkForError(e error) {
	if e != nil {
		log.Fatalln("Error reading csv file: ", e)
	}
}

// Return date format 'YYYY-MM-DD'
func getDateFormat(dateVal string) string {
	dateFormat, _ := time.Parse("20060102", dateVal)
	newDateStr := dateFormat.Format("2006-01-02")
	return newDateStr
}

// Return most recent date field from csv file
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

// Return false if csv column headers are not the expected otherwise true
func getLastDayV2(fileName, tagFile string) (string, bool) {
	// Try Open file
	csvFile, err := os.Open(fileName)
	if err != nil {
		log.Println("Can not open csv file: ", err)
		return "", false
	}
	defer csvFile.Close()

	// setup csvReader
	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = ';'

	// Check if column headers are the expected ones
	columnHead, _ := csvReader.Read() // get a columnHead string[]
	if tagFile == "positivos" {
		if columnHead[8] != "FECHA_RESULTADO" {
			log.Printf("format file '%v' unexpected. Something has changed", fileName)
			log.Printf("Expected 'FECHA_RESULTADO' at index 8, found '%v'", columnHead[8])
			return "", false
		}
	} else if tagFile == "fallecidos" {
		if columnHead[2] != "FECHA_FALLECIMIENTO" {
			log.Printf("format file '%v' unexpected. Something has changed", fileName)
			log.Printf("Expected 'FECHA_FALLECIMIENTO' at index 2, found '%v'", columnHead[2])
			return "", false
		}
	} else {
		log.Printf("TagFile not recognized: %v", tagFile)
		return "", false
	}

	// Return most recent Date
	record, err := csvReader.Read() // from second line get 'FECHA_CORTE' field
	if err != nil {
		log.Println(err)
		return "", false
	}
	return record[0], true
}

func getReportCases(fileName, dateRowStr string) model.CasesReport {
	fmt.Printf("\n**** getReportCases ****\n")
	isOK := true
	if dateRowStr == "" {
		dateRowStr, isOK = getLastDayV2(fileName, "positivos")
	}

	if isOK != true {
		log.Println("CSV File with unexpected column headers")
		return model.CasesReport{}
	}

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
	numNewCasesDate := 0
	totalCases := 0
	casesByDept := make(map[string]int)
	fileReport := fmt.Sprintf("dataFiles/report_positivos_%v.txt", dateRowStr)
	fileHand, e := os.Create(fileReport)
	if e != nil {
		panic(e)
	}
	defer fileHand.Close()
	wFile := bufio.NewWriter(fileHand)
	log.Println("Searching data for date:", dateRowStr)
	// iterate over the entire file
	for {
		record, err := csvReader.Read() // get a record string[]
		if err == io.EOF { break }
		checkForError(err)

		// Count total covid cases from the beginning of the pandemic
		totalCases++

		// Looking for new cases in last day
		if dateRowStr == record[8] {  // Compare specific 'DATE' with 'FECHA_RESULTADO' field
			// new case
			numNewCasesDate++
			// Write file
			line := fmt.Sprintln(record[0],record[1],record[2],record[3],record[4],record[5],record[6],record[7],record[8])
			_, e = wFile.WriteString(line)
			if e != nil {
				log.Println("Error Writing line to file")
				panic(e)
			}
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
	e = wFile.Flush()
	if e != nil {
		log.Println("Flush error")
		panic(e)
	}
	dateStr := getDateFormat(dateRowStr) // convert to date format 'YYYY-MM-DD'
	myNewReport := model.CasesReport{Date: dateStr, NewCases: numNewCasesDate, TotalCases: totalCases, NewCasesByDept: casesByDept}
	return myNewReport
}

func getReportDeceased(fileName, dateRowStr string) model.DeceasedReport {
	fmt.Printf("\n**** getReportDeceased ***\n")
	//lastDay := getLastDay(fileName)
	//lastDay, isOK := getLastDayV2(fileName, "fallecidos")
	//var lastDay string
	isOK := true
	if dateRowStr == "" {
		dateRowStr, isOK = getLastDayV2(fileName, "fallecidos")
	}

	if isOK != true {
		log.Println("CSV File with unexpected column headers")
		return model.DeceasedReport{}
	}

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
	numDeceasedDate := 0
	deceaseByDept := make(map[string]int)
	log.Println("Searching data for date:", dateRowStr)
	for {
		record, err := csvReader.Read() // get a record string[]
		if err == io.EOF { break }
		checkForError(err)

		// Count total deceases from the beginning of pandemic
		totalDeceased++

		// Looking for deceases in specific date
		if dateRowStr == record[2] {  // Compare specific 'DATE' with 'FECHA_FALLECIMIENTO'
			numDeceasedDate++
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

	dateStr := getDateFormat(dateRowStr) // convert to date format 'YYYY-MM-DD'
	myReportDeceases := model.DeceasedReport{Date: dateStr, NewDeceased: numDeceasedDate, TotalDeceased: totalDeceased, DeceasesByDept: deceaseByDept}
	return myReportDeceases
}