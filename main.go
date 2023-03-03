package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/brunobolting/go-twitch-chat/infra/repo"
	"github.com/brunobolting/go-twitch-chat/usecase/question"
	"github.com/brunobolting/go-twitch-chat/ws"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var interrupt os.Signal

func serveHome(w http.ResponseWriter, r *http.Request) {
    p := "./public" + r.URL.Path
    if p == "./" {
        p = "./public/index.html"
    }

	http.ServeFile(w, r, p)
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


	// if os.Getenv("RUN_FIXTURE") == "true" {
	// 	err = fixture.Run(questionService)
	// 	if err != nil {
	// 		log.Fatal("Error to apply fixtures: ", err)
	// 	}
	// }

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(w, r, questionService)
	})

	addr := os.Getenv("API_HOST")+":"+os.Getenv("API_PORT")
	log.Println("listening:", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
