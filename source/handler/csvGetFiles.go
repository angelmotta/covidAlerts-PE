package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func DownloadFile(url, filepath string) error {
	fmt.Println("**** DownloadFile ****")
	outFile, err := os.Create(filepath + ".tmp")
	if err != nil {
		log.Println("os.Create(): ", err)
		return err
	}
	defer outFile.Close()

	// Get CSV File Data
	resp, err := http.Get(url)
	if err != nil {
		log.Println("http.Get() connection failed, ", err)
		return err
	}
	defer resp.Body.Close()
	log.Println("Download File Started")

	// Write to file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		log.Println("io.Copy(): ", err)
		return err
	}
	log.Println("Download File Completed")

	// Rename temp file to original file
	err = os.Rename(filepath+".tmp", filepath)
	return nil
}