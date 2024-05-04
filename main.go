package main

import (
	"bufio"
	"fmt"
	"github.com/yutanagano/karei/internal/engine"
	"github.com/yutanagano/karei/internal/uci"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	logFile := openLogFile()
	defer logFile.Close()
	log.SetOutput(logFile)

	fromStdIn := startStdInReader()
	toStdOut := startStdOutWriter()
	uci.ConnectClient(fromStdIn, toStdOut)

	engine.Start()
	uci.ConnectEngine(engine.Out, engine.In)

	uci.Start()
}

func openLogFile() *os.File {
	homeDirectory := os.ExpandEnv("$HOME")

	err := os.Mkdir(homeDirectory+"/.cache", 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	err = os.Mkdir(homeDirectory+"/.cache/karei", 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	f, err := os.Create(homeDirectory + "/.cache/karei/karei.log")
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func startStdInReader() chan string {
	c := make(chan string)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			text, err := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			if err != io.EOF && len(text) > 0 {
				c <- text
			}
		}
	}()
	return c
}

func startStdOutWriter() chan string {
	c := make(chan string)
	go func() {
		for line := range c {
			fmt.Println(line)
		}
	}()
	return c
}
