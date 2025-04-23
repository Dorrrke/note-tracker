package server

import (
	"context"
	"net/http"

	"github.com/Dorrrke/note-tracker/gen/auth"
	"github.com/Dorrrke/note-tracker/internal/note-tracker/domain/models"
	"github.com/gin-gonic/gin"
)

func (s *ServerApi) registerUser(c *gin.Context) {
	var user models.User
	err := c.ShouldBindBodyWithJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := s.auth.Register(context.Background(),
		&auth.User{
			Name:     user.Name,
			Login:    user.Login,
			Password: user.Password,
		})
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
	c.Header("Authorization", "Bearer "+resp.Token)
}

func (s *ServerApi) loginUser(c *gin.Context) {
	var user models.UserRequest
	err := c.ShouldBindBodyWithJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := s.auth.Login(context.Background(),
		&auth.UserCredentials{
			Login:    user.Login,
			Password: user.Password,
		})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("Authorization", "Bearer "+resp.Token, 3600, "/users/login", "", false, true)

	c.Status(http.StatusOK)
}
