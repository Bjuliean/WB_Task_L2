package main

import (
	"fmt"
	"sort"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type StrKey []rune

func (s *StrKey) Len() int {
	return len(*s)
}

func (s *StrKey) Less(i int, j int) bool {
	return (*s)[i] < (*s)[j]
}

func (s *StrKey) Swap(i int, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}

func (s *StrKey) String() string {
	return string(*s)
}

func main() {
	strList := []string{
		"КАТЯП",
		"пятак",
		"пятка",
		"тяпка",
		"листок",
		"листок",
		"слиток",
		"столик",
		"килост",
		"абоба123",
		"123абоба",
	}

	fmt.Println(Anagram(strList))
}

func Anagram(strList []string) map[string][]string {
	res := make(map[string][]string, len(strList))
	buf := make(map[string]string, len(strList))

	for _, v := range strList {
		v = strings.ToLower(v)

		k := GenerateKey(&buf, v)

		if !ContainsKey(res[k], v) {
			res[k] = append(res[k], v)
		}
	}

	return res
}

func ContainsKey(strArr []string, key string) bool {
	for _, v := range strArr {
		if key == v {
			return true
		}
	}

	return false
}

func GenerateKey(m *map[string]string, str string) string {
	tmp := StrKey(str)
	res := ""
	sort.Sort(&tmp)

	if _, ok := (*m)[string(tmp)]; ok {
		res = (*m)[string(tmp)]
	} else {
		(*m)[string(tmp)] = str
		res = str
	}

	return res
}
