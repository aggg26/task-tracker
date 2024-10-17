package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"trackerApp/internal/services/dtos"

	"github.com/gin-gonic/gin"
)

// AllTasks godoc
// @Summary Get all tasks for a user
// @Description Retrieves all tasks associated with the user ID obtained from the context
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}{"tasks": []Task}
// @Failure 500 {object} gin.H{"error": string}
// @Router /tasks [get]
func (h *Handler) AllTasks(c *gin.Context) {
	userId, err := h.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tasks, err := h.services.ITaskService.Get(userId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{"tasks": tasks})
}

// TaskById godoc
// @Summary Get task by ID for a user
// @Description Retrieves a specific task associated with the given task ID and user ID obtained from the context
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} map[string]interface{}{"task": Task}
// @Failure 400 {object} gin.H{"error": string}
// @Failure 500 {object} gin.H{"error": string}
// @Router /tasks/{id} [get]
func (h *Handler) TaskById(c *gin.Context) {
	userId, err := h.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	taskId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	task, err := h.services.ITaskService.GetById(taskId, userId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{"task": task})
}

// PostTask godoc
// @Summary Create a new task for a user
// @Description Creates a new task associated with the user ID obtained from the context
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body dtos.CreateTask true "Create task request"
// @Success 200 {string} string "Task created message"
// @Failure 400 {object} gin.H{"error": string}
// @Failure 500 {object} gin.H{"error": string}
// @Router /tasks [post]
func (h *Handler) PostTask(c *gin.Context) {
	userId, err := h.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var request dtos.CreateTask
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	taskId, err := h.services.ITaskService.Create(userId, request)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, fmt.Sprintf("Task %d was created", taskId))
}

// PutTask godoc
// @Summary Update an existing task for a user
// @Description Updates a task associated with the given task ID and user ID obtained from the context
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param request body dtos.UpdateTask true "Update task request"
// @Success 200 {string} string "Task updated message"
// @Failure 400 {object} gin.H{"error": string}
// @Failure 500 {object} gin.H{"error": string}
// @Router /tasks/{id} [put]
func (h *Handler) PutTask(c *gin.Context) {
	userId, err := h.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var request dtos.UpdateTask
	taskId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	if err := h.services.ITaskService.Update(taskId, userId, request); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, "Task was updated")
}

// DeleteTask godoc
// @Summary Delete a task for a user
// @Description Deletes a specific task associated with the given task ID and user ID obtained from the context
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {string} string "Task deleted message"
// @Failure 400 {object} gin.H{"error": string}
// @Failure 500 {object} gin.H{"error": string}
// @Router /tasks/{id} [delete]
func (h *Handler) DeleteTask(c *gin.Context) {
	userId, err := h.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	taskId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	if err := h.services.ITaskService.Delete(taskId, userId); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, "Task was deleted")
}
