package main

func startEngine() (fromEngine, toEngine chan string) {
	tell("info string hello from engine")

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
