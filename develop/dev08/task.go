package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать свой собственный UNIX-шелл-утилиту с поддержкой ряда простейших команд:

- cd <args> - смена директории (в качестве аргумента могут быть то-то и то)
- pwd - показать путь до текущего каталога
- echo <args> - вывод аргумента в STDOUT
- kill <args> - "убить" процесс, переданный в качесте аргумента (пример: такой-то пример)
- ps - выводит общую информацию по запущенным процессам в формате *такой-то формат*
Так же требуется поддерживать функционал fork/exec-команд

*Шелл — это обычная консольная программа, которая будучи запущенной, в интерактивном сеансе выводит некое приглашение
в STDOUT и ожидает ввода пользователя через STDIN. Дождавшись ввода, обрабатывает команду согласно своей логике
и при необходимости выводит результат на экран. Интерактивный сеанс поддерживается до тех пор, пока не будет введена команда выхода (например \quit).
*/

type ShellData struct {
	Username           *user.User
	Hostname           string
	CurrentDir         string
	HostUserInfoColor  *color.Color
	DirectoryInfoColor *color.Color
}

func NewShellData() (*ShellData, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	curDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	huColor := color.New(color.FgMagenta).Add(color.Bold)
	diColor := color.New(color.FgGreen).Add(color.Bold)

	return &ShellData{
			Username:           user,
			Hostname:           host,
			CurrentDir:         curDir,
			HostUserInfoColor:  huColor,
			DirectoryInfoColor: diColor,
		},
		nil
}

func (s *ShellData) Info() {
	s.HostUserInfoColor.Printf("%s@%s:", s.Username.Username, s.Hostname)
	s.DirectoryInfoColor.Printf("~%s$ ", strings.ReplaceAll(s.CurrentDir, s.Username.HomeDir, ""))
}

func main() {
	shellData, err := NewShellData()
	if err != nil {
		log.Fatal(err)
	}

	sc := bufio.NewScanner(os.Stdin)

	shellData.Info()

	for sc.Scan() {
		args := ScanArgs(sc.Text())

		switch args[0] {
		case "exit", "quit":
			os.Exit(0)
		case "":
		case "pwd":
			PWDCommand(shellData)
		case "cd":
			err = CDCommand(shellData, args[1])
			if err != nil {
				fmt.Fprintf(os.Stderr, err.Error()+"\n")
			}
		case "ls":
			err := LSCommand()
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error()+"\n")
			}
		case "echo":
			EchoCommand(args[1:]...)
		case "ps":
			err := ExecCommand(args)
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error()+"\n")
			}
		case "kill":
			err := KillCommand(args[1:])
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error()+"\n")
			}
		case "fork", "exec":
			err := ExecCommand(args)
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error()+"\n")
			}
		default:
			fmt.Printf("myshell: %s: command not found\n", args[0])
		}

		shellData.Info()
	}

}

func KillCommand(args []string) error {
	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	err = proc.Kill()
	if err != nil {
		return err
	}

	return nil
}

func ExecCommand(args []string) error {
	buf := bytes.Buffer{}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Print(buf.String())

	return nil
}

func CDCommand(shellData *ShellData, path string) error {
	err := os.Chdir(path)
	if err != nil {
		return err
	}

	shellData.CurrentDir, err = os.Getwd()
	if err != nil {
		return err
	}

	return nil
}

func LSCommand() error {
	dir, err := os.ReadDir(".")
	if err != nil {
		return err
	}

	c := color.New(color.FgBlue).Add(color.Bold)

	for i := 0; i < len(dir); i++ {

		if dir[i].IsDir() {
			c.Printf("%s", dir[i].Name())
		} else {
			fmt.Printf("%s", dir[i].Name())
		}

		if i < len(dir)-1 {
			fmt.Printf(" ")
		}
	}
	fmt.Println()

	return nil
}

func EchoCommand(strs ...string) {
	fmt.Println(strings.Join(strs, " "))
}

func PWDCommand(shellData *ShellData) {
	fmt.Printf("%s\n", shellData.CurrentDir)
}

func ScanArgs(str string) []string {
	return strings.Split(str, " ")
}
