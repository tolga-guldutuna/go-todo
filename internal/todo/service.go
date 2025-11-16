package todo

import (
	"context"
	"errors"
)

type Service interface {
	GetTodos(ctx context.Context) ([]Todo, error)
	GetTodo(ctx context.Context, id int) (*Todo, error)
	AddTodo(ctx context.Context, title, description string) (*Todo, error)
	UpdateTodo(ctx context.Context, id int, title, description string, done bool) (*Todo, error)
	DeleteTodo(ctx context.Context, id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetTodos(ctx context.Context) ([]Todo, error) {
	return s.repo.GetAll(ctx)
}

func (s *service) GetTodo(ctx context.Context, id int) (*Todo, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) AddTodo(ctx context.Context, title, description string) (*Todo, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}
	return s.repo.Create(ctx, title, description)
}

func (s *service) UpdateTodo(ctx context.Context, id int, title, description string, done bool) (*Todo, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if title != "" {
		t.Title = title
	}
	if description != "" {
		t.Description = description
	}
	t.IsDone = done

	return s.repo.Update(ctx, *t)
}

func (s *service) DeleteTodo(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
