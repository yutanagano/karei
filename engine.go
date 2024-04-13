package main

import "fmt"

func startEngine() (fromEngine, toEngine chan string) {
	fmt.Println("info string Hello from engine")

	fromEngine = make(chan string)
	toEngine = make(chan string)

	go func() {
		for command := range toEngine {
			switch command {
			case "stop":
			case "quit":
			}
		}
	}()
	return
}
