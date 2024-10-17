package services

import (
	"database/sql"
	"time"
	"trackerApp/internal/models"
	"trackerApp/internal/services/dtos"
)

type TaskService struct {
	db *sql.DB
}

func NewTaskService(db *sql.DB) *TaskService {
	return &TaskService{db: db}
}

const (
	getTasks       = `SELECT id, title, description, is_complete, create_at FROM tasks WHERE user_id = $1`
	getTaskById    = `SELECT id, title, description, is_complete, create_at FROM tasks WHERE id = $1 AND user_id = $2`
	createTask     = `INSERT INTO tasks (title,description,is_complete,create_at, user_id) VALUES($1,$2,$3,$4,$5) RETURNING id`
	updateTaskById = `UPDATE tasks SET title=$1, description=$2, is_complete=$3 WHERE id=$4 AND user_id=$5`
	deleteTaskById = `DELETE FROM tasks WHERE id=$1 AND user_id=$2`
)

func (s *TaskService) Get(userId int) ([]models.Task, error) {
	var tasks []models.Task
	rows, err := s.db.Query(getTasks, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.IsCompleted, &task.CreateAt); err != nil {
			continue
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *TaskService) GetById(taskId, userId int) (*models.Task, error) {
	var task models.Task
	if err := s.db.QueryRow(getTaskById, taskId, userId).Scan(&task.ID, &task.Title, &task.Description, &task.IsCompleted, &task.CreateAt); err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *TaskService) Create(userId int, taskDto dtos.CreateTask) (int, error) {
	task := models.Task{
		Title:       taskDto.Title,
		Description: taskDto.Description,
		IsCompleted: false,
		CreateAt:    time.Now().Format("02/01/2006"),
	}

	var id int
	if err := s.db.QueryRow(createTask, task.Title, task.Description, task.IsCompleted, task.CreateAt, userId).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (s *TaskService) Update(taskId, userId int, updateTask dtos.UpdateTask) error {
	if _, err := s.db.Exec(updateTaskById, updateTask.Title, updateTask.Description, updateTask.IsComplete, taskId, userId); err != nil {
		return err
	}
	return nil
}

func (s *TaskService) Delete(taskId, userId int) error {
	if _, err := s.db.Exec(deleteTaskById, taskId, userId); err != nil {
		return err
	}
	return nil
}
