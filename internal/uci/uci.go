package uci

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/yutanagano/karei/internal/chess"
	"github.com/yutanagano/karei/internal/util"
)

var fromClient, toClient, fromEngine, toEngine chan string
var clientConnected, engineConnected = false, false

var debug, isInfinite = false, false
var currentPosition chess.Position
var bestMove string

func ConnectClient(toUCI, fromUCI chan string) {
	fromClient = toUCI
	toClient = fromUCI
	clientConnected = true
}

func ConnectEngine(toUCI, fromUCI chan string) {
	fromEngine = toUCI
	toEngine = fromUCI
	engineConnected = true
}

func Start() {
	if !(clientConnected && engineConnected) {
		err := errors.New("cannot start UCI without client/engine connection")
		fmt.Println(err.Error())
		log.Fatalln(err)
	}

	toClient <- "info string hello from karei"

Repl:
	for {
		select {
		case input := <-fromClient:
			tokens := util.Queue[string](strings.Split(input, " "))
			command := tokens.Pop()

			switch command {
			case "uci":
				handleUci()
			case "setoption":
				handleSetOption(tokens)
			case "isready":
				handleIsReady()
			case "ucinewgame":
				handleNewGame()
			case "position":
				handlePosition(tokens)
			case "debug":
				handleDebug(tokens)
			case "register":
				handleRegister(tokens)
			case "go":
				handleGo(tokens)
			case "ponderhit":
				handlePonderHit()
			case "stop":
				handleStop()
			case "quit":
				break Repl
			}
		case bestMove := <-fromEngine:
			handleBestMove(bestMove)
		}
	}
}

func handleUci() {
	toClient <- "id name Karei"
	toClient <- "id author Yuta Nagano"
	toClient <- "option name Hash type spin default 32 min 1 max 1024"
	toClient <- "option name Threads type spin default 1 min 1 max 16"
	toClient <- "uciok"
}

func handleSetOption(tokens util.Queue[string]) {
	toClient <- "info string set option " + strings.Join(tokens, " ")
	toClient <- "info string not implemented yet"
}

func handleIsReady() {
	toClient <- "readyok"
}

func handleNewGame() {
	currentPosition.LoadFEN(chess.GetStartingFEN())
}

func handlePosition(tokens util.Queue[string]) {
	// [fen <fenstring> | startpos] (moves <move1> ... <movei>)?
	var positionFen chess.FEN

	position_specifier := tokens.Pop()

	switch position_specifier {
	case "fen":
		positionFen = chess.FEN{
			BoardState:      tokens.Pop(),
			ActiveColour:    tokens.Pop(),
			CastlingRights:  tokens.Pop(),
			EnPassantSquare: tokens.Pop(),
			HalfMoveClock:   tokens.Pop(),
			FullMoveNumber:  tokens.Pop(),
		}

	case "startpos":
		positionFen = chess.GetStartingFEN()
	default:
		err := fmt.Errorf("expected position specifier to be 'fen' or 'startpos', got %s", position_specifier)
		toClient <- "info string " + err.Error()
		return
	}

	err := currentPosition.LoadFEN(positionFen)
	if err != nil {
		fmt.Printf("error loading FEN: %s", err.Error())
	}

	nextToken := tokens.Pop()
	if nextToken == "moves" {
		toClient <- "info string parsing moves: " + strings.Join(tokens, " ")
		// TODO parse moves
	}
}

func handleDebug(tokens util.Queue[string]) {
	// [ on | off ]
	mode := tokens.Pop()
	switch mode {
	case "on":
		debug = true
		toClient <- "info string debug mode on"
	case "off":
		debug = false
		toClient <- "info string debug mode off"
	default:
		toClient <- "info string unrecognised debug mode " + mode
	}
}

func handleRegister(tokens util.Queue[string]) {
	// register [ later | name <x> code <y> ]
	toClient <- "info string register not implemented"
}

func handleGo(tokens util.Queue[string]) {
	// go (searchmoves <move1> ... <movei>)? ponder? (wtime <x>)? (btime <x>)? (winc <x>)? (binc <x>)? (movestogo <x>)? (depth <x>)? (nodes <x>)? (mate <x>)? (movetime <x>)? infinite?
	if len(tokens) > 1 {
		switch (tokens)[1] {
		case "searchmoves":
			toClient <- "info string go searchmoves not implemented"
		case "ponder":
			toClient <- "info string go ponder not implemented"
		case "wtime":
			toClient <- "info string go wtime not implemented"
		case "btime":
			toClient <- "info string go btime not implemented"
		case "winc":
			toClient <- "info string go winc not implemented"
		case "binc":
			toClient <- "info string go binc not implemented"
		case "movestogo":
			toClient <- "info string go movestogo not implemented"
		case "depth":
			toClient <- "info string go depth not implemented"
		case "nodes":
			toClient <- "info string go nodes not implemented"
		case "mate":
			toClient <- "info string go mate not implemented"
		case "movetime":
			toClient <- "info string go movetime not implemented"
		case "infinite":
			toClient <- "info string go infinite not implemented"
		default:
			toClient <- "info string go " + tokens[1] + " not implemented"
		}
	} else {
		toClient <- "info string go not implemented"
	}
}

func handlePonderHit() {
	toClient <- "info string ponderhit not implemented"
}

func handleStop() {
	if isInfinite {
		if bestMove != "" {
			toClient <- bestMove
			bestMove = ""
		}

		toEngine <- "stop"
		isInfinite = false
	}
}

func handleBestMove(move string) {
	if isInfinite {
		bestMove = move
		return
	}
	toClient <- move
}
