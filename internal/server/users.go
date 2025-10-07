package server

import (
	"net/http"

	"github.com/Dorrrke/note-tracker/internal/domain/models"
	"github.com/gin-gonic/gin"
)

// cookieMaxAge максимальный срок жизни кук
const cookieMaxAge = 3600

// registerUser обрабатывает запрос на регистрацию нового пользователя.
//
// Ожидает JSON с данными пользователя в теле запроса.
// Возвращает UID созданного пользователя или ошибку при конфликте.
//
// Ответы:
//   - 201 Created: успешная регистрация, в теле — UID
//   - 400 Bad Request: неверный формат JSON
//   - 409 Conflict: пользователь с такими данными уже существует
func (s *API) registerUser(c *gin.Context) {
	var user models.User
	err := c.ShouldBindBodyWithJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, err := s.uService.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"uid": uid})
}

// loginUser обрабатывает запрос на аутентификацию пользователя.
//
// Принимает учётные данные в теле запроса.
// При успешной аутентификации устанавливает cookie с UID и возвращает его.
//
// Ответы:
//   - 200 OK: успешный вход, в теле — UID, установлена cookie
//   - 400 Bad Request: неверный формат запроса
//   - 401 Unauthorized: неверные учётные данные
func (s *API) loginUser(c *gin.Context) {
	var user models.UserRequest
	err := c.ShouldBindBodyWithJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, err := s.uService.LoginUser(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("uid", uid, cookieMaxAge, "/users/login", "", false, true)

	c.JSON(http.StatusOK, gin.H{"uid": uid})
}

// getUsers возвращает список всех пользователей.
//
// Возвращает массив пользователей или ошибку, если данные недоступны.
//
// Ответы:
//   - 200 OK: список пользователей
//   - 404 Not Found: ошибка при получении данных
func (s *API) getUsers(c *gin.Context) {
	users, err := s.uService.GetUsers()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// getUserByID возвращает данные пользователя по его уникальному идентификатору.
//
// Извлекает ID из URL-параметра ":id".
// Возвращает пользователя или ошибку, если он не найден.
//
// Ответы:
//   - 200 OK: данные пользователя
//   - 404 Not Found: пользователь не существует
func (s *API) getUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := s.uService.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// getUserProfile возвращает профиль текущего авторизованного пользователя.
//
// Ответы:
//   - 404 Not Found: функционал временно недоступен
func (s *API) getUserProfile(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "functionality comming soon"})
}

// createUser создаёт нового пользователя (аналог registerUser).
//
// Возвращает UID созданного пользователя.
//
// Ответы:
//   - 201 Created: пользователь создан
//   - 400 Bad Request: ошибка валидации JSON
//   - 409 Conflict: пользователь уже существует
func (s *API) createUser(c *gin.Context) {
	var user models.User
	err := c.ShouldBindBodyWithJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, err := s.uService.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"uid": uid})
}

// deleteUserByID удаляет пользователя по его уникальному идентификатору.
//
// Извлекает ID из URL и передаёт его в сервис удаления.
// Возвращает подтверждение успешного удаления.
//
// Ответы:
//   - 200 OK: пользователь удалён
//   - 404 Not Found: пользователь не найден
func (s *API) deleteUserByID(c *gin.Context) {
	id := c.Param("id")

	err := s.uService.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// updateUserByID обновляет данные пользователя по его ID.
//
// Принимает обновлённые данные в теле запроса и устанавливает UID из URL.
// Передаёт обновлённую структуру в сервис обновления.
//
// Ответы:
//   - 200 OK: обновлённые данные пользователя
//   - 400 Bad Request: неверный формат тела запроса
//   - 404 Not Found: пользователь не найден
func (s *API) updateUserByID(c *gin.Context) {
	id := c.Param("id")

	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updatedUser.UID = id

	err := s.uService.UpdateUser(updatedUser)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}
