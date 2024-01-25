package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

const (
	filePath = "txt.txt"
	pattern = "bla"
)

func TestGrep(t *testing.T) {
	cmdList := [][]string{
		{
			pattern,
			filePath,
			"-A",
			"1",
		},
		{
			pattern,
			filePath,
			"-B",
			"1",
		},
		{
			pattern,
			filePath,
			"-C",
			"1",
		},
		{
			pattern,
			filePath,
			"-c",
		},
		{
			pattern,
			filePath,
			"-i",
		},
		{
			pattern,
			filePath,
			"-v",
		},
		{
			pattern,
			filePath,
			"-F",
		},
		{
			pattern,
			filePath,
			"-n",
		},
	}

	for _, tc := range cmdList {
		t.Run("TEST 1", func(t *testing.T) {
			var outG bytes.Buffer
			cmdGrep := exec.Command("grep", tc...)
			cmdGrep.Stdout = &outG
			cmdGrep.Run()

			var outM bytes.Buffer
			cmdMy := exec.Command("./task", tc...)
			cmdMy.Stdout = &outM
			cmdMy.Run()

			if CleanStr(outG.String()) != outM.String() {
				t.Errorf("\nhas:\n%s\nexpected:\n%s\n", outM.String(), CleanStr(outG.String()))
			}

		})
	}
}

func CleanStr(str string) string {
	//tmp := strings.Split(str, "--")
	tmp := strings.Split(str, "\n")
	res := strings.Builder{}
	res.Grow(len(tmp))

	for i := 0; i < len(tmp); i++ {
		if tmp[i] != "--" {
			if i == len(tmp)-1 {
				res.WriteString(tmp[i])
			} else {
				res.WriteString(tmp[i] + "\n")
			}
			
		}
	}

	return res.String()
}
