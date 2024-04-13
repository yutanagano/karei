package main

import (
	"fmt"
	"strings"
)

func uci(fromGUI chan string) {
	fmt.Println("info string Hello from uci")

	fromEngine, toEngine := startEngine()

	quit := false
	var command string
	var command_parts []string
	var bestMove string

	for !quit {
		select {
		case command = <-fromGUI:
			command_parts = strings.Split(command, " ")
			switch command_parts[0] {
			case "uci":
				handleUci()
			case "isready":
				handleIsReady()
			case "setoption":
				handleSetOption(command_parts)
			case "stop":
				handleStop(toEngine)
			case "quit":
				quit = true
				continue
			}
		case bestMove = <-fromEngine:
			handleBestMove(bestMove)
			continue
		}
	}
}

func handleUci() {
	tell("id name Karei")
	tell("id author Yuta Nagano")
	tell("option name Hash type spin default 32 min 1 max 1024")
	tell("option name Threads type spin default 1 min 1 max 16")
	tell("uciok")
}

func handleIsReady() {
	tell("readyok")
}

func handleSetOption(option []string) {
	tell("info string set option ", strings.Join(option, " "))
	tell("info string not implemented yet")
}

func handleStop(toEngine chan string) {
	toEngine <- "stop"
}

func handleBestMove(bestMove string) {
	tell(bestMove)
}
