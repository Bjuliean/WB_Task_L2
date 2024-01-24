package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
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
	K *int
	N *bool
	R *bool
	U *bool
	M *bool
	B *bool
	C *bool
	H *bool
}

type stringBufferEl struct {
	StrNum int
	Str    []string
}

type StringBuffer struct {
	Arr        []stringBufferEl
	SortColNum int
}

func (s *StringBuffer) Len() int {
	return len(s.Arr)
}


// При изучении работы оригинального sort, было выяснено, что если
// символы столбца равны, то неявно сравниваются следующие столбцы.
// Поскольку в задании не сказано про то, какие символы могут быть в файле,
// это вынуждает конвертировать все в руны.
func (s *StringBuffer) Less(i int, j int) bool {
	a := []rune(strings.ToLower(s.Arr[i].Str[s.SortColNum]))
	b := []rune(strings.ToLower(s.Arr[j].Str[s.SortColNum]))

	if a[0] == b[0] && len(s.Arr[i].Str) > 1 && len(s.Arr[j].Str) > 1 {
		for z := 1; z < len(s.Arr[i].Str) && z < len(s.Arr[j].Str); z++ {
			fv := []rune(strings.ToLower(s.Arr[i].Str[z]))
			sv := []rune(strings.ToLower(s.Arr[j].Str[z]))

			if fv[0] != sv[0] {
				return fv[0] < sv[0]
			}
		}
	}

	return a[0] < b[0]
}

func (s *StringBuffer) Swap(i int, j int) {
	s.Arr[i], s.Arr[j] = s.Arr[j], s.Arr[i]
}

func main() {
	fv := &FlagsValue{
		K: flag.Int("k", 1, "указание колонки для сортировки"),
		N: flag.Bool("n", false, "сортировать по числовому значению"),
		R: flag.Bool("r", false, "сортировать в обратном порядке"),
		U: flag.Bool("u", false, "не выводить повторяющиеся строки"),
		M: flag.Bool("M", false, "сортировать по названию месяца"),
		B: flag.Bool("b", false, "игнорировать хвостовые пробелы"),
		C: flag.Bool("c", false, "проверять отсортированы ли данные"),
		H: flag.Bool("h", false, "сортировать по числовому значению с учётом суффиксов"),
	}

	flag.Parse()
	if *fv.K <= 0 {
		log.Fatal("sort: field number is zero: invalid field specification ‘0’")
	}

	res, _ := ParseFile("txt.txt")

	SortData(res, fv)

	for _, v := range res.Arr {
		fmt.Println(v.Str)
	}

	a, err := strconv.Atoi("10K")
	if err != nil {
		fmt.Println("aboba")
	}
	_ = a
}

func SortData(buffer *StringBuffer, fl *FlagsValue) {
	buffer.SortColNum = *fl.K - 1

	sort.Sort(buffer)
}

func ParseFile(path string) (*StringBuffer, error) {
	res := &StringBuffer{
		Arr:        make([]stringBufferEl, 0, 10),
		SortColNum: 0,
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	sc := bufio.NewScanner(f)

	i := 0
	for sc.Scan() {
		res.Arr = append(res.Arr, stringBufferEl{
			StrNum: i,
			Str:    strings.Split(sc.Text(), " "),
		})
		i++
	}

	return res, nil
}
