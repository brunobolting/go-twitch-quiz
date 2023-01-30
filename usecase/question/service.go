package question

import "github.com/brunobolting/go-twitch-chat/entity"

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) GetQuestion(id entity.ID) (*entity.Question, error) {
	e, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (s *Service) GetRandomQuestion(whereNotIn []string) (*entity.Question, error) {
	e, err := s.repo.GetRandom(whereNotIn)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (s *Service) CreateQuestion(question string, answers []string) (entity.ID, error) {
	e, err := entity.NewQuestion(question, answers)
	if err != nil {
		return e.ID, err
	}

	return s.repo.Create(e)
}

func (s *Service) UpdateQuestion(e *entity.Question) error {
	err := s.repo.Update(e)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteQuestion(id entity.ID) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
