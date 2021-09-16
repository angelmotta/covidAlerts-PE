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
	log.Println("Searching most recent date in csv file")
	// Open file
	csvFile, err := os.Open(fileName)
	if err != nil {
		log.Println("Can not open csv file: ", err)
		return "", false
	}
	defer csvFile.Close()

	// setup csvReader
	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = ';'

	var idxDateField int
	if tagFile == "positivos" {
		idxDateField = 7
	} else if tagFile == "fallecidos" {
		idxDateField = 1
	} else {
		log.Printf("TagFile not recognized: %v", tagFile)
		return "", false
	}

	// Get most recent Date
	_, _ = csvReader.Read() // read first line
	if err != nil {
		log.Println(err)
		return "", false
	}
	record, err := csvReader.Read() // from second line
	if err != nil {
		log.Println(err)
		return "", false
	}
	fechaCorte := record[0]
	mostRecentDateStr := record[idxDateField] // Get date from column at idxDateField 2 or 8
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
		currDateStr := record[idxDateField]
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

func isValidCsvFile(filename, tagFile string) bool {
	log.Println("Validating fields from CSV File")

	csvFile, err := os.Open(filename)
	if err != nil {
		log.Println("Error opening csv file: ", err)
		return false
	}
	defer csvFile.Close()

	// Setup csvReader
	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = ';'
	// Validate if column header of 'Date' is the expected one
	var idxDateField int
	lineHeaders, _ := csvReader.Read() // get a columnHead string[]

	if tagFile == "positivos" {
		idxDateField = 7
		if lineHeaders[idxDateField] != "FECHA_RESULTADO" {
			log.Printf("Unexpected format in csv file '%v'", filename)
			log.Printf("Expected 'FECHA_RESULTADO' at index 8, but found '%v'", lineHeaders[8])
			return false
		}
	} else if tagFile == "fallecidos" {
		idxDateField = 1
		if lineHeaders[idxDateField] != "FECHA_FALLECIMIENTO" {
			log.Printf("Unexpected format in csv file '%v'", filename)
			log.Printf("Expected 'FECHA_FALLECIMIENTO' at index 2, but found '%v'", lineHeaders[2])
			return false
		}
	} else {
		log.Printf("Tag file not recognized: %v", tagFile)
		return false
	}

	return true
}

func getReportCases(filename, dateRowStr string) model.CasesReport {
	fmt.Println("**** getReportCases ****")

	if dateRowStr == "" {
		isOK := isValidCsvFile(filename, "positivos")
		if isOK {
			dateRowStr, isOK = getLastDay(filename, "positivos")
		} else {
			log.Printf("Unexpected format in CSV File '%v'(review column name)\n", filename)
			return model.CasesReport{}
		}
	}

	// Try Open file
	csvFile, err := os.Open(filename)
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
	log.Printf("Searching data in '%v' for date: %v\n", filename, dateRowStr)
	idxDateResult := 7
	idxCity := 1
	for {
		record, err := csvReader.Read() // get a record string[]
		if err == io.EOF { break }
		checkForError(err)
		// Count total covid cases from the beginning of the pandemic
		totalCases++
		// Looking for new cases in last day
		if dateRowStr == record[idxDateResult] {  // Compare specific 'DATE' with 'FECHA_RESULTADO' field
			// new case
			numNewCasesDate++
			// Get cases by 'DEPARTAMENTO'
			dept := record[idxCity]
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

func getReportDeceased(filename, dateRowStr string) model.DeceasedReport {
	fmt.Println("**** getReportDeceased ***")

	if dateRowStr == "" {
		isOK := isValidCsvFile(filename, "fallecidos")
		if isOK {
			dateRowStr, isOK = getLastDay(filename, "fallecidos")
		} else {
			log.Printf("Unexpected format in CSV File '%v' (review column name)\n", filename)
			return model.DeceasedReport{}
		}
	}

	// Open file
	csvFile, err := os.Open(filename)
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
	log.Printf("Searching data in '%v' for date: %v\n", filename, dateRowStr)
	idxDeceasedDate := 1
	idxCity := 5
	for {
		record, err := csvReader.Read() // get a record string[]
		if err == io.EOF { break }
		checkForError(err)

		// Count total deceases from the beginning of pandemic
		totalDeceased++

		// Looking for deceases in specific date
		if dateRowStr == record[idxDeceasedDate] {  // Compare specific 'DATE' with 'FECHA_FALLECIMIENTO'
			numDeceasedDate++
			// Get decease by department
			dept := record[idxCity]
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