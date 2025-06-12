package server

import (
	"net/http"

	"github.com/Dorrrke/note-tracker/internal/domain/models"
	"github.com/gin-gonic/gin"
)

func (s *ServerApi) getTasks(c *gin.Context) {
	tasks, err := s.tService.GetTasks()
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

	if err := s.tService.CreateTask(task); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (s *ServerApi) getTaskByID(c *gin.Context) {
	id := c.Param("id")

	task, err := s.tService.GetTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (s *ServerApi) deleteTaskByID(c *gin.Context) {
	id := c.Param("id")

	err := s.tService.DeleteTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func (s *ServerApi) updateTaskByID(c *gin.Context) {
	id := c.Param("id")

	var updatedTask models.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updatedTask.TID = id

	err := s.tService.UpdateTask(updatedTask)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}
