package game

import (
	"encoding/json"
	"log"

	"github.com/brunobolting/go-twitch-chat/entity"
	"github.com/brunobolting/go-twitch-chat/twitch"
	"github.com/brunobolting/go-twitch-chat/usecase/question"
	"github.com/brunobolting/go-twitch-chat/ws"
)

type Game struct {
	Twitch *twitch.Twitch
	Hub *ws.Hub
	QuestionService *question.Service
	Question entity.Question
	PreviusQuestions []string
	Client *ws.Client
	RoundEnded bool
}

func NewGame(tw *twitch.Twitch, hub *ws.Hub, qs *question.Service) *Game {
	return &Game{
		Twitch: tw,
		Hub: hub,
		QuestionService: qs,
		RoundEnded: false,
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
		case received := <-g.Hub.Receive:
			if received.Message.Command == "START_ROUND" {
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
				g.Client = received.Client
				g.Hub.SendPrivate <- ws.Private{Client: received.Client, Message: json}
			}

			if received.Message.Command == "ROUND_END" {
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
				g.Hub.SendPrivate <- ws.Private{Client: g.Client, Message: json}
			}
		}
	}
}
