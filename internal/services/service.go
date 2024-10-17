package services

import (
	"database/sql"
	"trackerApp/internal/models"
	"trackerApp/internal/services/dtos"
)

// будет использоваться pgx драйвер

type ITaskService interface {
	Get(userId int) ([]models.Task, error)
	GetById(taskId, userId int) (*models.Task, error)
	Create(userId int, taskDto dtos.CreateTask) (int, error)
	Update(taskId, userId int, updateTask dtos.UpdateTask) error
	Delete(taskId, userId int) error
}

type IAuthService interface {
	AddUser(form dtos.UserForm) (int, error)
	GenerateJwt(form dtos.UserForm) (string, error)
	ParseJwt(accessToken string) (int, error)
}

type Service struct {
	ITaskService
	IAuthService
}

func NewService(db *sql.DB) *Service {
	return &Service{
		ITaskService: NewTaskService(db),
		IAuthService: NewAuthService(db),
	}
}
