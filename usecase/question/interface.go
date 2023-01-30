package question

import "github.com/brunobolting/go-twitch-chat/entity"

type Repository interface {
	Get(id entity.ID) (*entity.Question, error)
	GetRandom(whereNotIn []string) (*entity.Question, error)
	Create(e *entity.Question) (entity.ID, error)
	Update(e *entity.Question) error
	Delete(id entity.ID) error
}

type UseCase interface {
	GetQuestion(id entity.ID) (*entity.Question, error)
	GetRandomQuestion(whereNotIn []string) (*entity.Question, error)
	CreateQuestion(question string, answers []string) (entity.ID, error)
	UpdateQuestion(e *entity.Question) error
	DeleteQuestion(id entity.ID) error
}
