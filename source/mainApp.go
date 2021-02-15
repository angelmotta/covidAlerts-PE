package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

type Report struct {
	date string
	newCases int
	deceases int
	// newCasesByDept map
	// deceasesByDept map
}

func checkForError(e error) {
	if e != nil {
		log.Fatalln("Error opening csv file: ", e)
	}
}

func getLastDay(f *os.File) string {
	// setup csvReader
	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'

	// discard first line
	_, err := csvReader.Read()
	checkForError(err)

	// get from the first line 'FECHA_CORTE'
	record, err := csvReader.Read()
	return record[0]
}

func parseCsvFile(fileName string) Report {
	// Try Open file
	csvFile, err := os.Open(fileName)
	checkForError(err)
	defer csvFile.Close()

	lastDay := getLastDay(csvFile)
	fmt.Printf("last day: %s\n", lastDay)
	// setup a new csvReader and restore the pointer of the File descriptor
	csvFile.Seek(0, 0) // This is needed to avoid bugs about the pointer position
	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = ';'

	// discard first line
	_, err = csvReader.Read()
	checkForError(err)

	// iterate over the entire file
	newCasesLastDay := 0
	totalCases := 0
	deceasedLastDay := 0
	for {
		record, err := csvReader.Read() // get a record string[]
		if err == io.EOF { break }
		checkForError(err)

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

	// Display Results
	// TODO: get 'deceases' attribute
	fmt.Printf("Nuevos casos COVID (registrados el %s): %d contagiados\n", lastDay, newCasesLastDay)
	fmt.Printf("Número de fallecidos el día (%s): %d personas\n", lastDay, deceasedLastDay)
	fmt.Printf("Total casos positivos a la fecha: %d contagiados \n", totalCases)


	myNewReport := Report{lastDay, newCasesLastDay, deceasedLastDay}
	return myNewReport
}

func testModifiedStruct(newReport *Report) {
	newReport.date = "20210131"
	newReport.newCases = 100
	newReport.deceases = 30
}

func getStructReport() Report {
	newReport := Report{"20210131", 10, 3} // Struct literal
	return newReport
}


func main() {
	// process CSV File
	fileName := "dataFiles/positivos_covid_3_2_2021.csv"

	// Create a new Struct and send it by reference
	//report1 := Report{}
	//fmt.Println(report1)
	//testModifiedStruct(&report1)
	//fmt.Println(report1)

	// Invoke a function and receive a struct
	//report2 := getStructReport()
	//fmt.Println(report2)

	// Parse and receive a report
	report1 := parseCsvFile(fileName)
	fmt.Println(report1)
	fmt.Printf("{date: %s, newCases: %d, decease: %d}", report1.date, report1.newCases, report1.deceases)
}
