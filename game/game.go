package game

import (
	"encoding/json"
	"log"

	"github.com/brunobolting/go-twitch-chat/entity"
	"github.com/brunobolting/go-twitch-chat/twitch"
	"github.com/brunobolting/go-twitch-chat/usecase/question"
)

type Command struct {
	Command string `json:"command"`
}

type Game struct {
	Twitch *twitch.Twitch
	QuestionService *question.Service
	Question entity.Question
	PreviusQuestions []string
	RoundEnded bool
	Command chan *Command
	SendToClient chan []byte
	Close chan struct{}
}

func NewGame(tw *twitch.Twitch, qs *question.Service) *Game {
	return &Game{
		Twitch: tw,
		QuestionService: qs,
		RoundEnded: false,
		Command: make(chan *Command),
		SendToClient: make(chan []byte),
		Close: make(chan struct{}),
	}
}

func (g *Game) Start() error {
	q, err := g.QuestionService.GetRandomQuestion(g.PreviusQuestions)
	if err != nil {
		if err.Error() != "EOF" {
			return err
		}
		g.PreviusQuestions = []string{g.Question.ID}
		q, err = g.QuestionService.GetRandomQuestion(g.PreviusQuestions)
	}

	g.Question = *q
	g.PreviusQuestions = append(g.PreviusQuestions, g.Question.ID)
	g.RoundEnded = false

	return nil
}

func (g *Game) Run() {
	for {
		select {
		case received := <-g.Command:
			if received.Command == "START_ROUND" {
				g.RoundEnded = false
				err := g.Start()
				if err != nil {
					log.Fatalln(err)
				}
				data := map[string]string{"event": "NEW_ROUND", "question": g.Question.Question}
				json, err := json.Marshal(data)
				if err != nil {
					log.Fatalln(err)
				}
				g.SendToClient <- json
			}

			if received.Command == "ROUND_END" {
				g.RoundEnded = true
			}
		case message := <-g.Twitch.Answer:
			if g.Question.ValidateAnswer(message.Message) && g.RoundEnded == false {
				data := map[string]string{"event": "WINNER", "user": message.Author, "answer": message.Message}
				json, err := json.Marshal(data)
				if err != nil {
					log.Fatalln(err)
				}
				g.RoundEnded = true
				g.SendToClient <- json
			}
		case <-g.Close:
			return
		}
	}
}
