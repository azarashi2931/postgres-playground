package repositories

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"

	"github.com/koyashiro/postgres-playground/backend/models"
)

type PlaygroundRepository interface {
	GetAll() ([]*models.Playground, error)
	Get(id string) (*models.Playground, error)
	Set(p *models.Playground) error
	Delete(id string) error
}

type PlaygroundRepositoryImpl struct {
	ctx    context.Context
	client *redis.Client
}

func NewPlaygroundRepository() PlaygroundRepository {
	ctx := context.Background()
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &PlaygroundRepositoryImpl{
		ctx:    ctx,
		client: c,
	}
}

func (r *PlaygroundRepositoryImpl) GetAll() ([]*models.Playground, error) {
	ids, err := r.client.Keys(r.ctx, "*").Result()
	if err != nil {
		return nil, err
	}

	ps := make([]*models.Playground, 10)
	for _, id := range ids {
		p, err := r.Get(id)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}

	return ps, nil
}

func (r *PlaygroundRepositoryImpl) Get(id string) (*models.Playground, error) {
	b, err := r.client.Get(r.ctx, id).Bytes()
	if err != nil {
		return nil, err
	}

	var p *models.Playground
	if err = json.Unmarshal(b, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (r *PlaygroundRepositoryImpl) Set(p *models.Playground) error {
	return r.client.Set(r.ctx, p.ID, p, 0).Err()
}

func (r *PlaygroundRepositoryImpl) Delete(id string) error {
	return r.client.Del(r.ctx, id).Err()
}
