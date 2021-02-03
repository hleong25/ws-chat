package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	defer log.Println("ws-server finished")
	done := make(chan bool)

	SetupCloseHandler(done)

	httpServer := New(done)

	go httpServer.Start()

	log.Println("Starting application...")
	<-done
}

// SetupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
func SetupCloseHandler(done chan bool) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Interrupt detected")
		close(done)
	}()
}
