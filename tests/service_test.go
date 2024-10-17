package tests

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"testing"
	"trackerApp/internal/models"
	"trackerApp/internal/services"
	"trackerApp/internal/services/dtos"
)

const (
	createTables = `
	CREATE TABLE users(
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL UNIQUE
	);
	CREATE TABLE tasks (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT NOT NULL, 
		is_complete BOOLEAN DEFAULT FALSE,
		create_at VARCHAR(255) NOT NULL,
		user_id INT REFERENCES users(id) ON DELETE CASCADE
	);
	INSERT INTO users (username, password) VALUES ('test','test');
	INSERT INTO tasks (title, description, is_complete, create_at, user_id) VALUES ('task 1','description 1', false, '25/03/2004', 1);
	INSERT INTO tasks (title, description, is_complete, create_at, user_id) VALUES ('task 2','description 2', false, '25/03/2004', 1);
	INSERT INTO tasks (title, description, is_complete, create_at, user_id) VALUES ('task 3','description 3', false, '25/03/2004', 1);
`
	dropTables = `
	DROP TABLE tasks IF EXISTS; 
	DROP TABLE users IF EXISTS;`
)

var db *sql.DB

func setupDb() {
	var err error
	connStr := "user=postgres dbname=tracker_testdb password=qwerty host=localhost port=5432 sslmode=disable"
	db, err = sql.Open("pgx", connStr)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(createTables)
	if err != nil {
		panic(err)
	}
}

func teardown() {
	db.Exec(`DROP TABLE tasks; DROP TABLE users;`)
	db.Close()
}

func TestServices(t *testing.T) {
	setupDb()
	defer teardown()
	service := services.NewService(db)

	t.Run("AuthService", func(t *testing.T) {
		t.Run("AddUser", func(t *testing.T) {
			// arrange
			expectedId := 2
			form := dtos.UserForm{"user", "user"}
			// act
			id, err := service.IAuthService.AddUser(form)
			// assert
			assert.NoError(t, err)
			assert.Equal(t, expectedId, id)
		})
		t.Run("GenerateJwt", func(t *testing.T) {
			// arrange
			form := dtos.UserForm{"user", "user"}

			// act
			actual, err := service.IAuthService.GenerateJwt(form)
			// assert
			assert.NotEmpty(t, actual)
			assert.NoError(t, err)
		})
	})

	t.Run("TaskService", func(t *testing.T) {
		t.Run("GetTasks", func(t *testing.T) {
			// arrange
			userId := 1
			expected := []models.Task{
				{1, "task 1", "description 1", false, "25/03/2004"},
				{2, "task 2", "description 2", false, "25/03/2004"},
				{3, "task 3", "description 3", false, "25/03/2004"},
			}
			// act
			actual, err := service.ITaskService.Get(userId)
			// assert
			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		})
		t.Run("GetTaskById", func(t *testing.T) {
			// arrange
			userId, taskId := 1, 1
			expected := models.Task{1, "task 1", "description 1", false, "25/03/2004"}
			// act
			actual, err := service.ITaskService.GetById(taskId, userId)
			actTask := models.Task{1, actual.Title, actual.Description, actual.IsCompleted, actual.CreateAt}
			// assert
			assert.NoError(t, err)
			assert.Equal(t, expected, actTask)
		})
		t.Run("CreateTask", func(t *testing.T) {
			// arrange
			userId := 1
			expectedId := 4
			task := dtos.CreateTask{"test 4", "description 4"}
			// act
			actualId, err := service.ITaskService.Create(userId, task)
			// assert
			assert.NoError(t, err)
			assert.Equal(t, expectedId, actualId)
		})
		t.Run("UpdateTask", func(t *testing.T) {
			// arrange
			taskId := 4
			userId := 1
			update := dtos.UpdateTask{"test 4", "description 4", true}
			// act
			err := service.ITaskService.Update(taskId, userId, update)
			// assert
			assert.NoError(t, err)
			assert.NotEmpty(t, update)
		})
		t.Run("DeleteTask", func(t *testing.T) {
			// arrange
			userId := 1
			taskId := 4
			// act
			err := service.ITaskService.Delete(taskId, userId)
			// assert
			assert.NoError(t, err)
		})
	})
}
