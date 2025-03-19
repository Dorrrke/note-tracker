package server

import (
	"net/http"
	"time"

	"github.com/Dorrrke/note-tracker/internal/domain/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *ServerApi) getTasks(c *gin.Context) {
	tasks, err := s.repo.GetTasks()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (s *ServerApi) createTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindBodyWithJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.TID = uuid.New().String()

	now := time.Now()

	task.Stsatus = "pending"
	task.CreatedAt = now
	task.UpdatedAt = now

	if err := s.repo.SaveTask(task); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}
