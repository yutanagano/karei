package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	logFile := openLogFile()
	defer logFile.Close()
	log.SetOutput(logFile)

	tell("info string hello from karei")
	uci(startStdinReader())
	tell("info string goodbye!")
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

func stdoutTell(text ...string) {
	toStdout := ""
	for _, t := range text {
		toStdout += t
	}
	fmt.Println(toStdout)
}

func startStdinReader() chan string {
	line := make(chan string)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			text, err := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			if err != io.EOF && len(text) > 0 {
				line <- text
			}
		}
	}()
	return line
}
