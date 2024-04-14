package main

import (
	"fmt"
	"strings"
)

func uci(fromGUI chan string) {
	fmt.Println("info string Hello from uci")

	fromEngine, toEngine := startEngine()

	quit, isInfinite := false, false

	var command, bestMove, bestMoveCache string
	var command_parts []string

	for !quit {
		select {
		case command = <-fromGUI:
			command = strings.ToLower(command)
			command_parts = strings.Split(command, " ")
			switch command_parts[0] {
			case "uci":
				handleUci()
			case "setoption":
				handleSetOption(command_parts)
			case "isready":
				handleIsReady()
			case "ucinewgame":
				handleNewGame()
			case "position":
				handlePosition(command_parts)
			case "debug":
				handleDebug(command_parts)
			case "register":
				handleRegister(command_parts)
			case "go":
				handleGo(command_parts)
			case "ponderhit":
				handlePonderHit()
			case "stop":
				handleStop(toEngine, &isInfinite, &bestMoveCache)
			case "quit":
				quit = true
				continue
			}
		case bestMove = <-fromEngine:
			handleBestMove(bestMove, isInfinite, &bestMoveCache)
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

func handleSetOption(command_parts []string) {
	tell("info string set option ", strings.Join(command_parts, " "))
	tell("info string not implemented yet")
}

func handleIsReady() {
	tell("readyok")
}

func handleNewGame() {
	tell("info string ucinewgame not implemented")
}

func handlePosition(command_parts []string) {
	// position [fen <fenstring> | startpos] (moves <move1> ... <movei>)?
	if len(command_parts) > 1 {
		switch command_parts[1] {
		case "startpos":
			tell("info string position startpos not implemented")
		case "fen":
			tell("info string position fen not implemented")
		default:
			tell("info string position ", command_parts[1], " not implemented")
		}
	} else {
		tell("info string position not implemented")
	}
}

func handleDebug(command_parts []string) {
	// debug [ on | off ]
	tell("info string debug not implemented")
}

func handleRegister(command_parts []string) {
	// register [ later | name <x> code <y> ]
	tell("info string register not implemented")
}

func handleGo(command_parts []string) {
	// go (searchmoves <move1> ... <movei>)? ponder? (wtime <x>)? (btime <x>)? (winc <x>)? (binc <x>)? (movestogo <x>)? (depth <x>)? (nodes <x>)? (mate <x>)? (movetime <x>)? infinite?
	if len(command_parts) > 1 {
		switch command_parts[1] {
		case "searchmoves":
			tell("info string go searchmoves not implemented")
		case "ponder":
			tell("info string go ponder not implemented")
		case "wtime":
			tell("info string go wtime not implemented")
		case "btime":
			tell("info string go btime not implemented")
		case "winc":
			tell("info string go winc not implemented")
		case "binc":
			tell("info string go binc not implemented")
		case "movestogo":
			tell("info string go movestogo not implemented")
		case "depth":
			tell("info string go depth not implemented")
		case "nodes":
			tell("info string go nodes not implemented")
		case "mate":
			tell("info string go mate not implemented")
		case "movetime":
			tell("info string go movetime not implemented")
		case "infinite":
			tell("info string go infinite not implemented")
		default:
			tell("info string go ", command_parts[1], " not implemented")
		}
	} else {
		tell("info string go not implemented")
	}
}

func handlePonderHit() {
	tell("info string ponderhit not implemented")
}

func handleStop(toEngine chan string, isInfinite *bool, bestMoveCache *string) {
	if *isInfinite {
		if *bestMoveCache != "" {
			tell(*bestMoveCache)
			*bestMoveCache = ""
		}

		toEngine <- "stop"
		*isInfinite = false
	}
}

func handleBestMove(bestMove string, isInfinite bool, bestMoveCache *string) {
	if isInfinite {
		*bestMoveCache = bestMove
		return
	}
	tell(bestMove)
}
