package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/brunobolting/go-twitch-chat/game"
	"github.com/brunobolting/go-twitch-chat/infra/repo"
	"github.com/brunobolting/go-twitch-chat/twitch"
	"github.com/brunobolting/go-twitch-chat/usecase/question"
	"github.com/brunobolting/go-twitch-chat/ws"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var interrupt os.Signal

func init() {
    err := godotenv.Load(".env")

    if err != nil {
        log.Fatal("Error loading .env file")
    }
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "public")
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("DB_CONN")))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	questionRepo := repo.NewQuestionMongodb(client)
	questionService := question.NewService(questionRepo)

	// log.Println(questionService.CreateQuestion("nascido em 1889, artesao, escultor e roqueiro, ________ procurou refugio na guitarra", []string{"xuxa", "XXX"}))

	// log.Println(questionService.GetQuestion("3a976cf7-e872-41ed-83fa-9806df677605"))
	// log.Println(questionService.GetRandomQuestion())

	hub := ws.NewHub()
	go hub.Run()

	tw := twitch.NewTwitch()
	go tw.Run(os.Getenv("TWITCH_CHANNEL"), hub)

	game := game.NewGame(tw, hub, questionService)
	go game.Run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})

	addr := os.Getenv("API_HOST")+":"+os.Getenv("API_PORT")
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
