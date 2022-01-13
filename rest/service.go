package main

import (
	"context"

	"github.com/TheMickeyMike/grpc-rest-bench/warehouse"
)

type Service interface {
	List(context.Context) ([]*warehouse.UserAccount, error)
	Get(context.Context, string) (*warehouse.UserAccount, error)
	Delete(context.Context, string)
}

type UserService struct {
	db *warehouse.Db
}

func NewUserService(db *warehouse.Db) Service {
	return &UserService{db}
}

func (s *UserService) List(context.Context) ([]*warehouse.UserAccount, error) {
	var users = make([]*warehouse.UserAccount, len(s.db.Data))
	for _, u := range s.db.Data {
		users = append(users, u)
	}
	return users, nil
}

func (s *UserService) Get(_ context.Context, id string) (*warehouse.UserAccount, error) {
	if user, ok := s.db.Data[id]; ok {
		return user, nil
	}
	return nil, nil
}

func (s *UserService) Delete(_ context.Context, id string) {
	delete(s.db.Data, id)
}
