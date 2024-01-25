package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type FlagsValue struct {
	K int
	N bool
	R bool
	U bool
}

// Поскольку в задании не сказано про то, какие символы могут быть в файле,
// это вынуждает использовать руны.
type stringBufferEl struct {
	StrNum int
	Str    [][]rune
}

type Comparator interface {
	Compare(i, j int, s *StringBuffer) bool
}

type CompareK struct {
}

// При изучении работы оригинального sort, было выяснено, что если
// символы столбца равны, то неявно сравниваются следующие столбцы.
// Если одна из строчек состоит из 1 столбца, то именно она будет считаться
// меньше, в случае равенства символов
func (c *CompareK) Compare(i, j int, s *StringBuffer) bool {
	a := unicode.ToLower(s.Arr[i].Str[s.SortColNum][0])
	b := unicode.ToLower(s.Arr[j].Str[s.SortColNum][0])

	if a == b {
		if CountSymbols(s.Arr[i].Str) == CountSymbols(s.Arr[j].Str) && len(s.Arr[j].Str) > 1 {
			for z := 1; z < len(s.Arr[i].Str) && z < len(s.Arr[j].Str); z++ {
				fv := unicode.ToLower(s.Arr[i].Str[z][0])
				sv := unicode.ToLower(s.Arr[j].Str[z][0])

				if fv != sv {
					return fv < sv
				}
			}
		} else {
			return CountSymbols(s.Arr[i].Str) < CountSymbols(s.Arr[j].Str)
		}
	}

	return a < b
}

func CountSymbols(str [][]rune) int {
	res := 0

	for _, v := range str {
		res += len(v)
	}

	return res + len(str)
}

type CompareN struct {
}

func (c *CompareN) Compare(i, j int, s *StringBuffer) bool {
	a, b := 0, 0

	a, errA := strconv.Atoi(string(s.Arr[i].Str[s.SortColNum]))
	if errA != nil {
		_, errB := strconv.Atoi(string(s.Arr[j].Str[s.SortColNum]))
		if errB != nil {
			funcK := &CompareK{}
			return funcK.Compare(i, j, s)
		}
		return true
	}

	b, errB := strconv.Atoi(string(s.Arr[j].Str[s.SortColNum]))
	if errB != nil {
		return false
	}

	return a < b
}

type StringBuffer struct {
	Arr        []stringBufferEl
	SortColNum int
	LongestCol int
	Comp       Comparator
}

func (s *StringBuffer) Len() int {
	return len(s.Arr)
}

func (s *StringBuffer) Less(i int, j int) bool {
	return s.Comp.Compare(i, j, s)
}

func (s *StringBuffer) Swap(i int, j int) {
	s.Arr[i], s.Arr[j] = s.Arr[j], s.Arr[i]
}

func main() {
	os.Stdout.WriteString(Sort(os.Args[1]))
}

func Sort(path string) string {
	fv := InitFlags()

	data, err := ParseFile(path)
	if err != nil {
		log.Fatal(err)
	}

	SortProcess(data, fv)

	return Output(data)
}

func Output(buffer *StringBuffer) string {
	str := strings.Builder{}
	str.Grow(len(buffer.Arr) * buffer.LongestCol)

	for i := 0; i < len(buffer.Arr); i++ {
		for z := 0; z < len(buffer.Arr[i].Str); z++ {
			str.WriteString(string(buffer.Arr[i].Str[z]))
			if z != len(buffer.Arr[i].Str)-1 {
				str.WriteString(" ")
			}
		}
		str.WriteString("\n")
	}

	return str.String()
}

func SortProcess(buffer *StringBuffer, fl *FlagsValue) {
	buffer.SortColNum = fl.K - 1
	if buffer.SortColNum > buffer.LongestCol {
		buffer.SortColNum = 0
	}
	buffer.Comp = &CompareK{}
	sort.Sort(buffer)

	if fl.N {
		buffer.Comp = &CompareN{}
		sort.Sort(buffer)
	}

	if fl.R {
		Reverse(buffer)
	}

	if fl.U {
		Unique(buffer, fl)
	}

}

func InitFlags() *FlagsValue {
	fv := &FlagsValue{
		K: 1,
		N: false,
		R: false,
		U: false,
	}

	for i := 0; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-k":
			if i < len(os.Args)-1 {
				num, err := strconv.Atoi(os.Args[i+1])
				if err != nil {
					fv.K = -1
					break
				}
				fv.K = num
			} else {
				fv.K = -1
				break
			}
		case "-n":
			fv.N = true
		case "-r":
			fv.R = true
		case "-u":
			fv.U = true
		}
	}

	if fv.K <= 0 {
		log.Fatal(fmt.Sprintf("sort: field number is zero: invalid field specification '%d'", fv.K))
	}

	return fv
}

func Unique(buffer *StringBuffer, fl *FlagsValue) {
	m := make(map[string]int, len(buffer.Arr))
	str := strings.Builder{}
	str.Grow(buffer.LongestCol)

	uList := make([]stringBufferEl, 0, len(buffer.Arr))
	for i := 0; i < len(buffer.Arr); i++ {
		if buffer.SortColNum > 0 {
			m[string(buffer.Arr[i].Str[buffer.SortColNum])]++
			if m[string(buffer.Arr[i].Str[buffer.SortColNum])] == 1 {
				uList = append(uList, buffer.Arr[i])
			}
		} else {
			for z := 0; z < len(buffer.Arr[i].Str); z++ {
				str.WriteString(string(buffer.Arr[i].Str[z]))
			}
			m[str.String()]++
			if m[str.String()] == 1 {
				uList = append(uList, buffer.Arr[i])
			}
		}

		str.Reset()
	}

	buffer.Arr, uList = uList, buffer.Arr
}

func Reverse(buffer *StringBuffer) {
	for i, z := 0, len(buffer.Arr)-1; i != z; i, z = i+1, z-1 {
		if i >= z {
			break
		}
		buffer.Arr[i], buffer.Arr[z] = buffer.Arr[z], buffer.Arr[i]
	}
}

func ParseFile(path string) (*StringBuffer, error) {
	res := &StringBuffer{
		Arr:        make([]stringBufferEl, 0, 10),
		SortColNum: 0,
		LongestCol: 0,
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	i := 0
	for sc.Scan() {
		tmp := strings.Split(sc.Text(), " ")
		var runeArr [][]rune
		cntr := 0

		for _, v := range tmp {
			runeArr = append(runeArr, []rune(v))
			cntr++
		}

		if res.LongestCol < cntr {
			res.LongestCol = cntr
		}

		res.Arr = append(res.Arr, stringBufferEl{
			StrNum: i,
			Str:    runeArr,
		})
		i++
	}

	return res, nil
}
