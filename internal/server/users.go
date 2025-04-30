package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Dorrrke/note-tracker/internal/domain/models"
)

const cookieTimeout = 3600

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

	c.SetCookie("uid", uid, cookieTimeout, "/users/login", "", false, true)

	c.JSON(http.StatusOK, gin.H{"uid": uid})
}
