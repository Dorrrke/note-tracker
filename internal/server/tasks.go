package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Dorrrke/note-tracker/internal/domain/models"
)

func (s *API) getTasks(c *gin.Context) {
	tasks, err := s.tService.GetTasks()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (s *API) createTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindBodyWithJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.tService.CreateTask(task); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (s *API) saveTasks(c *gin.Context) {
	var tasks []models.Task
	if err := c.ShouldBindBodyWithJSON(&tasks); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.tService.SaveTasks(tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, "tasks saved")
}
