package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type HTTPServer struct {
	done     chan bool
	upgrader websocket.Upgrader
}

func New(done chan bool) HTTPServer {
	return HTTPServer{
		done:     done,
		upgrader: websocket.Upgrader{},
	}
}

func (httpServer HTTPServer) Start() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("hello world")

		msg := fmt.Sprintf("hello world %s\n", time.Now().Local())

		w.Write([]byte(msg))
	})

	http.HandleFunc("/quit", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("quit")
		close(httpServer.done)
	})

	http.HandleFunc("/time", httpServer.currentTime)

	log.Printf("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}

func (httpServer HTTPServer) currentTime(w http.ResponseWriter, r *http.Request) {

	ws, err := httpServer.upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	defer ws.Close()

	doneWs := make(chan bool)

	go func() {
		for {
			msgType, msgBytes, err := ws.ReadMessage()

			if err != nil {
				log.Println("failed to read message:", err)
			}

			if msgType == websocket.CloseMessage {
				close(doneWs)
				return
			}

			msg := strings.TrimSpace(string(msgBytes[:]))

			log.Printf("type:%d msg:%s", msgType, msg)

			if msg == "quit" {
				close(doneWs)
				return
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				msg := fmt.Sprintf("ws %s\n", time.Now().Local())
				ws.WriteMessage(websocket.TextMessage, []byte(msg))
			}
		}
	}()

	<-doneWs
}
