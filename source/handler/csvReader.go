package handler

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/angelmotta/covidAlerts-PE/source/model"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

// Return date format 'YYYY-MM-DD'
func convertDateFormat(dateVal string) string {
	dateFormat, _ := time.Parse("20060102", dateVal)
	newDateStr := dateFormat.Format("2006-01-02")
	return newDateStr
}

// Return false if csv column headers are not the expected otherwise true
func getLastDay(fileName, tagFile string) (string, error) {
	log.Println("Searching most recent date in csv file")
	// Open file
	csvFile, err := os.Open(fileName)
	if err != nil {
		log.Println("Can not open csv file: ", err)
		return "", err
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
		log.Printf("tag file not recognized: %v", tagFile)
		return "", errors.New("tag file not recognized")
	}

	// Get most recent Date
	_, err = csvReader.Read() // read first line
	if err != nil {
		log.Println(err)
		return "", err
	}

	mostRecentDate := 0
	for {
		// Read record
		record, err := csvReader.Read() // get a record string[]
		if err == io.EOF { break }	// EOF
		if err != nil {
			log.Println("error reading csv File:", err)
		}

		dateField := record[idxDateField]
		fechaCorte := record[0]

		if dateField == "" { continue }
		if fechaCorte == dateField {
			log.Println("'FECHA_CORTE' field in csv file is valid as most recent date", fechaCorte)
			return fechaCorte, nil
		}
		dateFieldRecord, err := strconv.Atoi(dateField)
		if err != nil {
			log.Printf("error casting date field in line with tokens: %v\n", dateField)
			log.Println("error casting date string to int", err)
			return "", err
		}

		if dateFieldRecord > mostRecentDate {
			mostRecentDate = dateFieldRecord
		}
	}

	mostRecentDateStr := strconv.Itoa(mostRecentDate)
	log.Println("'FECHA_CORTE' field has not records associated in csv file")
	log.Println("Record with most recent date in csv file is:", mostRecentDateStr)
	return mostRecentDateStr, nil
}

func isValidCsvFile(filename, tagFile string) bool {
	log.Println("Validating fields header in CSV File")

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

func getReportCases(filename, dateRowStr string) (model.CasesReport, error) {
	fmt.Println("**** getReportCases ****")
	var err error
	if dateRowStr == "" {
		isOK := isValidCsvFile(filename, "positivos")
		if isOK {
			dateRowStr, err = getLastDay(filename, "positivos")
			if err != nil {
				return model.CasesReport{}, err
			}
		} else {
			log.Printf("unexpected format in CSV File '%v'(review column name)\n", filename)
			return model.CasesReport{}, errors.New("unexpected format in CSV File")
		}
	}

	// Open file
	csvFile, err := os.Open(filename)
	if err != nil {
		return model.CasesReport{}, err
	}
	defer csvFile.Close()

	// Setup a csv reader
	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = ';'
	// discard first line
	_, err = csvReader.Read()
	if err != nil {
		return model.CasesReport{}, err
	}

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
		if err != nil {
			return model.CasesReport{}, err
		}
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

	dateStr := convertDateFormat(dateRowStr) // convert to date format 'YYYY-MM-DD'
	myNewReport := model.CasesReport{Date: dateStr, NewCases: numNewCasesDate, TotalCases: totalCases, NewCasesByDept: casesByDept}
	return myNewReport, nil
}

func getReportDeceased(filename, dateRowStr string) (model.DeceasedReport, error) {
	fmt.Println("**** getReportDeceased ***")
	var err error
	if dateRowStr == "" {
		isOK := isValidCsvFile(filename, "fallecidos")
		if isOK {
			dateRowStr, err = getLastDay(filename, "fallecidos")
			if err != nil {
				return model.DeceasedReport{}, err
			}
		} else {
			log.Printf("unexpected format in CSV File '%v' (review column name)\n", filename)
			return model.DeceasedReport{}, errors.New("unexpected format in CSV File")
		}
	}

	// Open file
	csvFile, err := os.Open(filename)
	if err != nil {
		return model.DeceasedReport{}, err
	}
	defer csvFile.Close()

	// Setup a csv reader
	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = ';'
	// discard first line of column headers
	_, err = csvReader.Read()
	if err != nil {
		return model.DeceasedReport{}, err
	}

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
		if err != nil {
			return model.DeceasedReport{}, err
		}

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

	dateStr := convertDateFormat(dateRowStr) // convert to date format 'YYYY-MM-DD'
	myReportDeceases := model.DeceasedReport{Date: dateStr, NewDeceased: numDeceasedDate, TotalDeceased: totalDeceased, DeceasesByDept: deceaseByDept}
	return myReportDeceases, nil
}