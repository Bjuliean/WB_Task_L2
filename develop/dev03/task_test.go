package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"testing"
)

const (
	filePath = "txt.txt"
)

func TestSort(t *testing.T) {
	cmdList := [][]string{
		{
			filePath,
			"-k",
			"1",
		},
		{
			filePath,
			"-n",
		},
		{
			filePath,
			"-r",
		},
		{
			filePath,
			"-u",
		},
	}

	for i, tc := range cmdList {
		t.Run(fmt.Sprintf("TEST %d", i), func(t *testing.T) {
			var outb bytes.Buffer
			cmd := exec.Command("sort", tc...)
			cmd.Stdout = &outb
			cmd.Run()

			var outr bytes.Buffer
			res := exec.Command("./task", tc...)
			res.Stdout = &outr
			err := res.Run()
			if err != nil {
				log.Fatal(err)
			}

			if outr.String() != outb.String() {
				t.Errorf("TEST %d - \"%s\" - has: \"\n%s\";\nexpected: \"\n%s\"",
					i, tc, outr.String(), outb.String())
			}
		})
	}
}
