package main

import (
	"reflect"
	"testing"
)

func TestAnagram(t *testing.T) {
	testCases := struct {
		arr []string
		m   map[string][]string
	}{
		arr: []string{
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
		},

		m: map[string][]string{
			"катяп": {
				"катяп",
				"пятак",
				"пятка",
				"тяпка",
			},

			"листок": {
				"листок",
				"слиток",
				"столик",
				"килост",
			},

			"абоба123": {
				"абоба123",
				"123абоба",
			},
		},
	}

	t.Run("TEST 1", func(t *testing.T) {
		res := Anagram(testCases.arr)

		if !reflect.DeepEqual(res, testCases.m) {
			t.Errorf("\nHAS: %v\nEXPECT: %v\n", res, testCases.m)
		}
	})
}
