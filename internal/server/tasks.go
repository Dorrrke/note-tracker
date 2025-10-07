package server

import (
	"net/http"

	"github.com/Dorrrke/note-tracker/internal/domain/models"
	"github.com/gin-gonic/gin"
)

// getTasks возвращает список всех задач.
//
// Вызывает сервис для получения всех задач из хранилища.
// Возвращает массив задач или ошибку при отсутствии данных.
//
// Ответы:
//   - 200 OK: список задач
//   - 404 Not Found: ошибка при получении задач
func (s *API) getTasks(c *gin.Context) {
	tasks, err := s.tService.GetTasks()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// createTask создаёт новую задачу.
//
// Принимает данные задачи в формате JSON, генерирует уникальный ID и временные метки.
// Сохраняет задачу через сервис.
//
// Ответы:
//   - 201 Created: задача успешно создана, возвращается её полная структура
//   - 400 Bad Request: ошибка разбора JSON
//   - 409 Conflict: конфликт при сохранении
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

// getTaskByID возвращает задачу по её уникальному идентификатору.
//
// Извлекает ID из URL-параметра ":id" и запрашивает задачу у сервиса.
//
// Ответы:
//   - 200 OK: данные задачи
//   - 404 Not Found: задача не найдена
func (s *API) getTaskByID(c *gin.Context) {
	id := c.Param("id")

	task, err := s.tService.GetTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// deleteTaskByID удаляет задачу по её ID.
//
// Извлекает ID из URL и передаёт его в сервис удаления.
//
// Ответы:
//   - 200 OK: задача успешно удалена
//   - 404 Not Found: задача не найдена
func (s *API) deleteTaskByID(c *gin.Context) {
	id := c.Param("id")

	err := s.tService.DeleteTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// updateTaskByID обновляет существующую задачу по её ID.
//
// Принимает обновлённые данные в теле запроса, устанавливает TID из URL
// и передаёт задачу в сервис обновления.
//
// Ответы:
//   - 200 OK: обновлённая задача
//   - 400 Bad Request: неверный формат тела запроса
//   - 404 Not Found: задача не найдена
func (s *API) updateTaskByID(c *gin.Context) {
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
