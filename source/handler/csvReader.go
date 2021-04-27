package handler

import (
	"encoding/csv"
	"fmt"
	"github.com/angelmotta/covidAlerts-PE/source/model"
	"io"
	"log"
	"os"
	"strconv"
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

// Return false if csv column headers are not the expected otherwise true
func getLastDay(fileName, tagFile string) (string, bool) {
	log.Println("Searching most recent date in csv file...")
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

	// Check if column header is the expected one
	idx := 0
	columnHead, _ := csvReader.Read() // get a columnHead string[]
	if tagFile == "positivos" {
		idx = 8
		if columnHead[idx] != "FECHA_RESULTADO" {
			log.Printf("format file '%v' unexpected. Something has changed", fileName)
			log.Printf("Expected 'FECHA_RESULTADO' at index 8, found '%v'", columnHead[8])
			return "", false
		}
	} else if tagFile == "fallecidos" {
		idx = 2
		if columnHead[idx] != "FECHA_FALLECIMIENTO" {
			log.Printf("format file '%v' unexpected. Something has changed", fileName)
			log.Printf("Expected 'FECHA_FALLECIMIENTO' at index 2, found '%v'", columnHead[2])
			return "", false
		}
	} else {
		log.Printf("TagFile not recognized: %v", tagFile)
		return "", false
	}

	// Get most recent Date
	record, err := csvReader.Read() // from second line
	if err != nil {
		log.Println(err)
		return "", false
	}
	fechaCorte := record[0]
	mostRecentDateStr := record[idx] // Get date from column at idx 2 or 8
	mostRecentDate, err := strconv.Atoi(mostRecentDateStr)
	if err != nil {
		log.Println("Error casting date string to int", err)
		return "", false
	}
	// If 'FECHA_CORTE' has a record in csv return it as a most recent date
	if fechaCorte == mostRecentDateStr {
		log.Println("'FECHA_CORTE' is valid as most recent date")
		return fechaCorte, true
	}
	// Read the rest of file and get the most recent Date
	for {
		record, err := csvReader.Read() // get a record string[]
		if err == io.EOF { break }	// EOF
		if err != nil {
			log.Println("Error reading csv File:", err)
		}
		currDateStr := record[idx]
		if currDateStr == "" { continue }
		// If 'FECHA_CORTE' has at least one record return it as a valid most recent date
		if fechaCorte == currDateStr {
			log.Println("'FECHA_CORTE' in csv file is valid. Most recent date", fechaCorte)
			return fechaCorte, true
		}
		currDate, err := strconv.Atoi(currDateStr)
		if err != nil {
			log.Println(record)
			log.Println(currDateStr)
			log.Println("Error casting date string to int", err)
			return "", false
		}
		if currDate > mostRecentDate {
			mostRecentDate = currDate
		}
	}
	mostRecentDateStr = strconv.Itoa(mostRecentDate)
	log.Println("'FECHA_CORTE' has not records in csv file")
	log.Println("Most recent date found in csv file is:", mostRecentDateStr)
	return mostRecentDateStr, true
}

func getReportCases(fileName, dateRowStr string) model.CasesReport {
	fmt.Println("**** getReportCases ****")
	isOK := true
	if dateRowStr == "" {
		dateRowStr, isOK = getLastDay(fileName, "positivos")
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

	// vars to return in Struct attributes
	numNewCasesDate := 0
	totalCases := 0
	casesByDept := make(map[string]int)

	// iterate over the CSV file
	log.Printf("Searching data in '%v' for date: %v\n", fileName, dateRowStr)
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

	dateStr := getDateFormat(dateRowStr) // convert to date format 'YYYY-MM-DD'
	myNewReport := model.CasesReport{Date: dateStr, NewCases: numNewCasesDate, TotalCases: totalCases, NewCasesByDept: casesByDept}
	return myNewReport
}

func getReportDeceased(fileName, dateRowStr string) model.DeceasedReport {
	fmt.Println("**** getReportDeceased ***")
	isOK := true
	if dateRowStr == "" {
		dateRowStr, isOK = getLastDay(fileName, "fallecidos")
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

	// vars to return in Struct attributes
	totalDeceased := 0
	numDeceasedDate := 0
	deceaseByDept := make(map[string]int)

	// Iterate over the CSV file
	log.Printf("Searching data in '%v' for date: %v\n", fileName, dateRowStr)
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