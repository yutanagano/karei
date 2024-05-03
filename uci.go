package main

import (
	"fmt"
	"strings"
)

func uci(fromGUI chan string) {
	tell("info string hello from uci")

	fromEngine, toEngine := startEngine()

	quit, isInfinite := false, false

	var bestMoveCache string
	var tokens []string

	tell("info string listening")

	for !quit {
		select {
		case input := <-fromGUI:
			input = strings.ToLower(input)
			tokens = strings.Split(input, " ")
			command := popFromQueue(&tokens)

			switch command {
			case "uci":
				handleUci()
			case "setoption":
				handleSetOption(&tokens)
			case "isready":
				handleIsReady()
			case "ucinewgame":
				handleNewGame()
			case "position":
				handlePosition(&tokens)
			case "debug":
				handleDebug(&tokens)
			case "register":
				handleRegister(&tokens)
			case "go":
				handleGo(&tokens)
			case "ponderhit":
				handlePonderHit()
			case "stop":
				handleStop(toEngine, &isInfinite, &bestMoveCache)
			case "quit":
				quit = true
				continue
			}
		case bestMove := <-fromEngine:
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

func handleSetOption(tokens *[]string) {
	tell("info string set option ", strings.Join(*tokens, " "))
	tell("info string not implemented yet")
}

func handleIsReady() {
	tell("readyok")
}

func handleNewGame() {
	tell("info string ucinewgame not implemented")
}

func handlePosition(tokens *[]string) {
	// [fen <fenstring> | startpos] (moves <move1> ... <movei>)?
	var positionFen fen

	position_specifier := popFromQueue(tokens)

	switch position_specifier {
	case "fen":
		positionFen = fen{
			boardState:      popFromQueue(tokens),
			activeColour:    popFromQueue(tokens),
			castlingRights:  popFromQueue(tokens),
			enPassantSquare: popFromQueue(tokens),
			halfMoveClock:   popFromQueue(tokens),
			fullMoveNumber:  popFromQueue(tokens),
		}

	case "startpos":
		positionFen = fen{
			boardState:      "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR",
			activeColour:    "w",
			castlingRights:  "KQkq",
			enPassantSquare: "-",
			halfMoveClock:   "0",
			fullMoveNumber:  "1",
		}

	default:
		err := fmt.Errorf("expected position specifier to be 'fen' or 'startpos', got %s", position_specifier)
		tell("info string ", err.Error())
		return
	}

	err := currentPosition.loadFEN(positionFen)
	if err != nil {
		fmt.Printf("error loading FEN: %s", err.Error())
	}

	nextToken := popFromQueue(tokens)
	if nextToken == "moves" {
		tell("info string parsing moves: ", strings.Join(*tokens, " "))
		parseMoves(tokens)
	}
}

func handleDebug(tokens *[]string) {
	// [ on | off ]
	mode := popFromQueue(tokens)
	switch mode {
	case "on":
		debug = true
		tell("info string debug mode on")
	case "off":
		debug = false
		tell("info string debug mode off")
	default:
		tell("info string unrecognised debug mode ", mode)
	}
}

func handleRegister(tokens *[]string) {
	// register [ later | name <x> code <y> ]
	tell("info string register not implemented")
}

func handleGo(tokens *[]string) {
	// go (searchmoves <move1> ... <movei>)? ponder? (wtime <x>)? (btime <x>)? (winc <x>)? (binc <x>)? (movestogo <x>)? (depth <x>)? (nodes <x>)? (mate <x>)? (movetime <x>)? infinite?
	if len(*tokens) > 1 {
		switch (*tokens)[1] {
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
			tell("info string go ", (*tokens)[1], " not implemented")
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

func popFromQueue(queue *[]string) string {
	if len(*queue) == 0 {
		return ""
	}
	result := (*queue)[0]
	*queue = (*queue)[1:]
	return result
}
