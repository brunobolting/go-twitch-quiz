package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/brunobolting/go-twitch-chat/infra/repo"
	"github.com/brunobolting/go-twitch-chat/usecase/question"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Question struct {
	Question string `json:"question"`
	Answers []string `json:"answers"`
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("DB_CONN")))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	questionRepo := repo.NewQuestionMongodb(client)
	questionService := question.NewService(questionRepo)

	jsonFile, err := os.Open("./fixture/questions.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalln(err)
	}

	var questions []Question

	err = json.Unmarshal(byteValue, &questions)
	if err != nil {
		log.Fatalln(err)
	}

	for _, q := range questions {
		_, err := questionService.CreateQuestion(q.Question, q.Answers)
		if err != nil {
			log.Fatalln(err)
		}
	}

	log.Println("fixture run successfully")
}
