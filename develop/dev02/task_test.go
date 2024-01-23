package main

import (
	"fmt"
	"testing"
)

func TestUnpack(t *testing.T) {
	testCases := []struct {
		Test   string
		Expect string
	}{
		{"qwe", "qwe"},
		{"a4bc2d5e", "aaaabccddddde"},
		{"abcd", "abcd"},
		{"45", ""},
		{"", ""},
		{`qwe\4\5`, "qwe45"},
		{`qwe\45`, "qwe44444"},
		{`qwe\\5`, `qwe\\\\\`},
		{`qwe\\5\`, ""},
		{"\\23\\34", "2223333"},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TEST %d", i), func(t *testing.T) {
			res, _ := Unpack(tc.Test)

			if res != tc.Expect {
				t.Errorf("TEST %d - \"%s\" - has: \"%s\"; expected: \"%s\"",
					i, tc.Test, res, tc.Expect)
			}
		})
	}
}
