package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/brunobolting/go-twitch-chat/twitch"
	"github.com/brunobolting/go-twitch-chat/ws"
)

var addr = flag.String("addr", ":8080", "http service address")
var interrupt os.Signal

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "public/index.html")
}

func main() {
	// interrupt := make(chan os.Signal)

	flag.Parse()

	hub := ws.NewHub()
	go hub.Run()

	tw := twitch.NewTwitch()
	go tw.Run("brunobolting", hub)

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
