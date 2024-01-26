package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type FlagData struct {
	Filename string
	F        int
	D        string
	S        bool
}

func main() {
	flagData, err := InitFlagData()
	if err != nil {
		log.Fatal(err)
	}

	textData, err := LoadFile(flagData.Filename)
	if err != nil {
		log.Fatal(err)
	}

	Cut(textData, flagData)
}

func Cut(textData []string, flagData *FlagData) {
	for i := 0; i < len(textData); i++ {
		if flagData.S && !strings.Contains(textData[i], flagData.D) {
			continue
		} else if !flagData.S && !strings.Contains(textData[i], flagData.D) {
			fmt.Println(textData[i])
			continue
		}
		tmp := strings.Split(textData[i], flagData.D)
		if flagData.F < len(tmp) {
			fmt.Println(tmp[flagData.F])
		}
	}
}

func LoadFile(path string) ([]string, error) {
	res := make([]string, 0, 30)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	for sc.Scan() {
		res = append(res, sc.Text())
	}

	return res, nil
}

func InitFlagData() (*FlagData, error) {
	res := &FlagData{}

	if len(os.Args) < 3 {
		return nil, fmt.Errorf("format: [EXECUTABLE] [FILE] [FLAGS]\n")
	}

	res.Filename = os.Args[1]

	for i := 2; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-d":
			if i == len(os.Args)-1 {
				return nil, fmt.Errorf("cut: option requires an argument -- 'd'")
			}
			if len(os.Args[i+1]) > 1 {
				return nil, fmt.Errorf("cut: the delimiter must be a single character")
			}

			res.D = os.Args[i+1]
			i++
		case "-f":
			if i == len(os.Args)-1 {
				return nil, fmt.Errorf("cut: option requires an argument -- 'f'")
			}

			val, err := strconv.Atoi(os.Args[i+1])
			if err != nil {
				return nil, err
			}
			if val <= 0 {
				return nil, fmt.Errorf("cut: fields are numbered from 1")
			}
			res.F = val - 1
		case "-s":
			res.S = true
		}
	}

	return res, nil
}
