package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

type FlagsData struct {
	Host    *string
	Port    *int
	Timeout *time.Duration
}

func main() {
	flagsData, err := ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("host: %s | port: %d | timeout %v\n", *flagsData.Host,
		*flagsData.Port, *flagsData.Timeout)

	err = Connection(flagsData)
	if err != nil {
		log.Fatal(err)
	}
}

func Connection(flagsData *FlagsData) error {
	con, err := net.DialTimeout("tcp",
		fmt.Sprintf("%s:%d", *flagsData.Host, *flagsData.Port),
		*flagsData.Timeout)
	if err != nil {
		return err
	}
	defer con.Close()

	go func() {
		io.Copy(con, os.Stdin)

		fmt.Println("SHUTDOWN")
		os.Exit(0)

	}()

	go func() {
		io.Copy(os.Stdout, con)

		fmt.Println("Server closed connection")
		os.Exit(0)

	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)

	<-c
	fmt.Println("SHUTDOWN")

	return nil
}

func ParseFlags() (*FlagsData, error) {
	res := &FlagsData{
		Host:    flag.String("host", "", "host"),
		Port:    flag.Int("port", 80, "port"),
		Timeout: flag.Duration("timeout", time.Second*10, "connection timeout"),
	}

	flag.Parse()

	return res, nil
}
