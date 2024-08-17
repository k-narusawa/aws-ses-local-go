package todo

import (
	"go-app-template/domain"
	"go-app-template/domain/repository"
	"go-app-template/domain/value"
)

type Service struct {
	todoRepo repository.ToDoRepository
	userRepo repository.UserRepository
}

func NewService(
	todoRepo repository.ToDoRepository,
	userRepo repository.UserRepository,
) *Service {
	return &Service{
		todoRepo: todoRepo,
		userRepo: userRepo,
	}
}

type CreateToDoInput struct {
	Title  string       `json:"title"`
	UserID value.UserID `json:"-"`
}

func (s *Service) Create(in CreateToDoInput) (*domain.ToDo, error) {
	user, err := s.userRepo.FindByID(in.UserID)
	if err != nil {
		return nil, err
	}

	todo := domain.NewToDo(in.Title, user.UserID)
	if err := s.todoRepo.Store(*todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *Service) FindByUserId(userId value.UserID) ([]domain.ToDo, error) {
	todos, err := s.todoRepo.FindByUserID(userId)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

type ChangeStatusInput struct {
	UserID value.UserID `json:"-"`
	ToDoID value.ToDoID `json:"-"`
	Done   bool         `json:"done"`
}

func (s *Service) ChangeToDoDone(in ChangeStatusInput) (*domain.ToDo, error) {
	todo, err := s.todoRepo.FindByID(in.ToDoID)
	if err != nil {
		return nil, err
	}

	todo.ChangeDone(value.OfDone(in.Done))
	if err := s.todoRepo.Update(*todo); err != nil {
		return nil, err
	}

	return todo, nil
}
