package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type FlagData struct {
	Filename string
	Pattern  string
	A_       bool
	B_       bool
	C_       bool
	C        bool
	I        bool
	V        bool
	F_       bool
	N        bool
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

	grepData := SearchForPattern(textData, flagData)

	ApplyFlags(textData, grepData, flagData)

	Output(grepData, flagData)
}

func ApplyFlags(textData, grepData []string, flagData *FlagData) {
	switch {
	case flagData.N:
		FlagN(textData, grepData)
	}
}

func SearchForPattern(textData []string, flagData *FlagData) []string {
	res := make([]string, 0, len(textData))
	tmp := make([]string, 0, len(textData))

	for i := 0; i < len(textData); i++ {
		if flagData.I {
			tmp = append(tmp, strings.ToLower(textData[i]))
		} else {
			tmp = append(tmp, textData[i])
		}
	}

	if flagData.V {
		for i := 0; i < len(textData); i++ {
			if !strings.Contains(tmp[i], flagData.Pattern) {
				res = append(res, textData[i])
			}
		}
	} else {
		for i := 0; i < len(textData); i++ {
			if strings.Contains(tmp[i], flagData.Pattern) {
				res = append(res, textData[i])
			}
		}
	}

	return res
}

func Output(grepData []string, flagData *FlagData) {
	if flagData.C {
		fmt.Println(len(grepData))
	} else {
		for _, v := range grepData {
			fmt.Println(v)
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
		return nil, fmt.Errorf("format: [EXECUTABLE] [PATTERN] [FILE] [FLAGS]\n")
	}

	res.Pattern = os.Args[1]
	res.Filename = os.Args[2]

	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-A":
			res.A_ = true
		case "-B":
			res.B_ = true
		case "-C":
			res.C_ = true
		case "-c":
			res.C = true
		case "-i":
			res.I = true
		case "-v":
			res.V = true
		case "-f":
			res.F_ = true
		case "-n":
			res.N = true
		default:
			return nil, fmt.Errorf("incorrect option entered\n")
		}
	}

	return res, nil
}

func FlagN(textData, grepData []string) {
	z := 0
	for i := 0; i < len(textData); i++ {
		if z < len(grepData) && textData[i] == grepData[z] {
			grepData[z] = fmt.Sprintf("%d:%s", i+1, grepData[z])
			z++
		}
	}
}
