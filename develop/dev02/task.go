package main

import (
	"fmt"
	"log"
	"strconv"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	str := `qwe\\5\`
	res, err := Unpack(str)

	fmt.Println(res)
	
	if err != nil {
		log.Fatal(err.Error())
	}
}

func Unpack(str string) (string, error) {
	if len(str) == 0 {
		return "", nil
	}
	
	temp := []rune(str)
	if unicode.IsDigit(temp[0]) || temp[len(temp)-1] == '\\' {
		return "", fmt.Errorf("unpack: incorrect string")
	}

	var res []rune

	for i := 0; i < len(temp); i++ {
		if temp[i] == '\\' {
			res = append(res, temp[i+1])
			i++
			continue
		}

		if !unicode.IsDigit(temp[i]) {
			res = append(res, temp[i])
			continue
		}

		iter, err := strconv.Atoi(string(temp[i]))
		if err != nil {
			return "", err
		}

		for z := 0; z < iter-1; z++ {
			res = append(res, temp[i-1])
		}
	}

	return string(res), nil
}
