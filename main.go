package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	tell("info string hello from karei")
	uci(startStdinReader())
	tell("info string goodbye!")
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
