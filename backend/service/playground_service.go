package service

import (
	"github.com/google/uuid"

	"github.com/koyashiro/postgres-playground/backend/model"
	"github.com/koyashiro/postgres-playground/backend/repository"
)

type PlaygroundService interface {
	GetAll() ([]*model.Playground, error)
	Get(id string) (*model.Playground, error)
	Create(db string) (*model.Playground, error)
	Destroy(id string) error
	Execute(id string, query string) (string, error)
}

type PlaygroundServiceImpl struct {
	pr  repository.PlaygroundRepository
	cr  repository.ContainerRepository
}

func NewPlaygroundService(
	pr repository.PlaygroundRepository,
	cr repository.ContainerRepository,
) PlaygroundService {
	return &PlaygroundServiceImpl{
		pr:  pr,
		cr:  cr,
	}
}

func (s *PlaygroundServiceImpl) GetAll() ([]*model.Playground, error) {
	return s.pr.GetAll()
}

func (s *PlaygroundServiceImpl) Get(id string) (*model.Playground, error) {
	return s.pr.Get(id)
}

func (s *PlaygroundServiceImpl) Create(db string) (*model.Playground, error) {
	c, err := s.cr.Create()
	if err != nil {
		return nil, err
	}

	p := &model.Playground{
		ID:        uuid.NewString(),
		DB:        db,
		Version:   "13.0.0",
		Container: c,
	}

	if err := s.pr.Set(p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *PlaygroundServiceImpl) Destroy(id string) error {
	p, err := s.pr.Get(id)
	if err != nil {
		return err
	}

	if err := s.cr.Delete(p.Container.ID); err != nil {
		return err
	}

	return s.pr.Delete(id)
}

func (s *PlaygroundServiceImpl) Execute(id string, query string) (string, error) {
	// TODO: execute query
	result := "XXXXXXXX"
	return result, nil
}