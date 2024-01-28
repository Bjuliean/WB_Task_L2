package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

const(
	DirName = "save_dir"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("error: need address\n")
	}

	err := Wget(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}

func Wget(url string) error {
	if !strings.Contains(url, "http://") && !strings.Contains(url, "https://") {
		url = "http://"+url
	}
	
	dirPath, err := os.Getwd()
	if err != nil {
		return err
	}
	dirPath += "/"+DirName

	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fileName := path.Base(url)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	file, err := os.Create(dirPath+"/"+fileName)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}