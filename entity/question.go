package entity

import (
	"strings"
	"time"
)

type Question struct {
	ID ID `bson:"_id"`
	Question string `bson:"question"`
	Answers []string `bson:"answers"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func NewQuestion(question string, answers []string) (*Question, error) {
	if len(answers) == 0 {
		return nil, ErrAnswersCannotBeEmpty
	}

	return &Question{
		ID: NewID(),
		Question: question,
		Answers: answers,
		CreatedAt: time.Now(),
	}, nil
}

func (q *Question) ValidateAnswer(answer string) bool {
	for _, ans := range q.Answers {
		if strings.EqualFold(answer, ans) {
			return true
		}
	}

	return false
}
