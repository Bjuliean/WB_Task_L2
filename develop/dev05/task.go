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
	Filename     string
	Pattern      string
	A_           bool
	B_           bool
	C_           bool
	C            bool
	I            bool
	V            bool
	F_           bool
	N            bool
	CopyDataFunc BufferMaker
	ContainsFunc StringComparator
	StrBefore    int
	StrAfter     int
}

//----------------------------------------------------

type BufferMaker interface {
	Load(buf, textData []string, indx int) []string
}

type BufferMakerDefault struct {
}

func (b *BufferMakerDefault) Load(buf, textData []string, indx int) []string {
	return append(buf, textData[indx])
}

type BufferMakerI struct {
}

func (b *BufferMakerI) Load(buf, textData []string, indx int) []string {
	return append(buf, strings.ToLower(textData[indx]))
}

//----------------------------------------------------

type StringComparator interface {
	Search(str, substr string) bool
}

type StringComparatorDefault struct {
}

func (s *StringComparatorDefault) Search(str, substr string) bool {
	return strings.Contains(str, substr)
}

type StringComparatorV struct {
}

func (s *StringComparatorV) Search(str, substr string) bool {
	return !strings.Contains(str, substr)
}

type StringComparatorF struct {
}

func (s *StringComparatorF) Search(str, substr string) bool {
	return str == substr
}

//----------------------------------------------------

func main() {
	Grep()
}

func Grep() {
	flagData, err := InitFlagData()
	if err != nil {
		log.Fatal(err)
	}

	textData, err := LoadFile(flagData.Filename)
	if err != nil {
		log.Fatal(err)
	}

	grepData := SearchForPattern(textData, flagData)

	ApplyFlags(textData, &grepData, flagData)

	Output(grepData, flagData)
}

func ApplyFlags(textData []string, grepData *[]string, flagData *FlagData) {
	switch {
	case flagData.N:
		FlagN(textData, *grepData)
	}
}

func SearchForPattern(textData []string, flagData *FlagData) []string {
	res := make([]string, 0, len(textData))
	tmp := make([]string, 0, len(textData))

	for i := 0; i < len(textData); i++ {
		tmp = flagData.CopyDataFunc.Load(tmp, textData, i)
	}

	for i := 0; i < len(textData); i++ {
		if flagData.ContainsFunc.Search(tmp[i], flagData.Pattern) {
			if flagData.B_ || flagData.C_ {
				bef := flagData.StrBefore
				for z := i - 1; z >= 0 && bef > 0; z-- {
					if !flagData.ContainsFunc.Search(textData[z], flagData.Pattern) {
						res = append(res, textData[z])
					}
					bef--
				}
			}
			res = append(res, textData[i])
			if flagData.A_ || flagData.C_ {
				aft := flagData.StrAfter
				for z := i + 1; z < len(textData) && aft > 0; z++ {
					if !flagData.ContainsFunc.Search(textData[z], flagData.Pattern) {
						res = append(res, textData[z])
					}
					aft--
				}
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
	res.CopyDataFunc = &BufferMakerDefault{}
	res.ContainsFunc = &StringComparatorDefault{}

	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-A":
			if i == len(os.Args)-1 {
				return nil, fmt.Errorf("grep: option requires an argument -- 'A'\n")
			}
			var err error
			res.StrAfter, err = strconv.Atoi(os.Args[i+1])
			if err != nil || res.StrAfter < 0 {
				return nil, fmt.Errorf("grep: invalid context length argument\n")
			}
			res.A_ = true
		case "-B":
			if i == len(os.Args)-1 {
				return nil, fmt.Errorf("grep: option requires an argument -- 'B'\n")
			}
			var err error
			res.StrBefore, err = strconv.Atoi(os.Args[i+1])
			if err != nil || res.StrAfter < 0 {
				return nil, fmt.Errorf("grep: invalid context length argument\n")
			}
			res.B_ = true
		case "-C":
			if i == len(os.Args)-1 {
				return nil, fmt.Errorf("grep: option requires an argument -- 'C'\n")
			}
			var err error
			res.StrAfter, err = strconv.Atoi(os.Args[i+1])
			if err != nil || res.StrAfter < 0 {
				return nil, fmt.Errorf("grep: invalid context length argument\n")
			}
			res.StrBefore = res.StrAfter
			res.C_ = true
		case "-c":
			res.C = true
		case "-i":
			res.I = true
			res.CopyDataFunc = &BufferMakerI{}
		case "-v":
			res.V = true
			res.ContainsFunc = &StringComparatorV{}
		case "-f":
			res.F_ = true
			res.ContainsFunc = &StringComparatorF{}
		case "-n":
			res.N = true
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

func FlagA(grepData *[]string, flagsData *FlagData) {
	flagsData.StrAfter++
	z := flagsData.StrAfter
	last := 0

	for i := 0; i < len(*grepData); i++ {
		z--
		if z < 0 {
			*grepData = append((*grepData)[:i], append([]string{"--"}, (*grepData)[i:]...)...)
			z = flagsData.StrAfter
			last = i
		}
	}

	*grepData = append((*grepData)[:last], (*grepData)[last+1:]...)
}
