package engine

var In, Out chan string

func Start() {
	In = make(chan string)
	Out = make(chan string)

	go func() {
		for command := range In {
			switch command {
			case "stop":
			case "quit":
			}
		}
	}()
	return
}
