package main

import (
	"context"
)

type Service interface {
	List(context.Context) ([]*UserAccount, error)
	Get(context.Context, string) (*UserAccount, error)
	Delete(context.Context, string)
}

type UserService struct {
	db *Db
}

func NewUserService(db *Db) Service {
	return &UserService{db}
}

func (s *UserService) List(context.Context) ([]*UserAccount, error) {
	var users = make([]*UserAccount, len(s.db.data))
	for _, u := range s.db.data {
		users = append(users, u)
	}
	return users, nil
}

func (s *UserService) Get(_ context.Context, id string) (*UserAccount, error) {
	if user, ok := s.db.data[id]; ok {
		return user, nil
	}
	return nil, nil
}

func (s *UserService) Delete(_ context.Context, id string) {
	delete(s.db.data, id)
}
