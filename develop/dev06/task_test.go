package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"
)

const(
	fileName = "txt.txt"
)

func TestCut(t *testing.T) {
	testCases := [][]string{
		{
			fileName,
			"-d",
			":",
			"-f",
			"1",
		},
		{
			fileName,
			"-d",
			":",
			"-f",
			"2",
		},
		{
			fileName,
			"-d",
			":",
			"-f",
			"1",
			"-s",
		},
	}

	for i, v := range testCases {
		t.Run(fmt.Sprintf("TEST %d", i), func(t *testing.T) {
			cmdCut := exec.Command("cut", v...)
			cutBuf := bytes.Buffer{}
			cmdCut.Stdout = &cutBuf
			cmdCut.Run()

			cmdMy := exec.Command("./task", v...)
			myBuf := bytes.Buffer{}
			cmdMy.Stdout = &myBuf
			cmdMy.Run()

			if cutBuf.String() != myBuf.String() {
				t.Errorf("has - \n%s\nexpect - \n%s\n", myBuf.String(), cutBuf.String())
			}
		})
	}
}