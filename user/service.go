package user

import (
	"context"

	"github.com/TheMickeyMike/grpc-rest-bench/pb"
)

type UserService interface {
	pb.UsersServer
	List(context.Context) ([]*pb.UserAccount, error)
	Get(context.Context, string) (*pb.UserAccount, error)
	Delete(context.Context, string)
}

type Service struct {
	pb.UnimplementedUsersServer
	db *Db
}

func NewService(db *Db) UserService {
	return &Service{db: db}
}

func (s *Service) GetUser(_ context.Context, user *pb.UserRequest) (*pb.UserAccount, error) {
	if user, ok := s.db.Data[user.Id]; ok {
		return user, nil
	}
	return nil, nil
}

func (s *Service) List(context.Context) ([]*pb.UserAccount, error) {
	var users = make([]*pb.UserAccount, len(s.db.Data))
	for _, u := range s.db.Data {
		users = append(users, u)
	}
	return users, nil
}

func (s *Service) Get(_ context.Context, id string) (*pb.UserAccount, error) {
	if user, ok := s.db.Data[id]; ok {
		return user, nil
	}
	return nil, nil
}

func (s *Service) Delete(_ context.Context, id string) {
	delete(s.db.Data, id)
}
