package server

import (
	"net/http"

	"github.com/Dorrrke/note-tracker/internal/domain/models"
	"github.com/gin-gonic/gin"
)

const cookieMaxAge = 3600

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

func (s *API) getUsers(c *gin.Context) {
	users, err := s.uService.GetUsers()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (s *API) getUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := s.uService.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (s *API) getUserProfile(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "functionality comming soon"})
}

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

func (s *API) deleteUserByID(c *gin.Context) {
	id := c.Param("id")

	err := s.uService.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

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
