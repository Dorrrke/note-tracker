package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Dorrrke/note-tracker/internal/config"
	"github.com/Dorrrke/note-tracker/internal/service"
	"github.com/Dorrrke/note-tracker/pkg/logger"
	"github.com/golang-jwt/jwt/v4"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var jwtKey = []byte("super_secret_key_1285@ggwg")

type ServerApi struct {
	server   *http.Server
	valid    *validator.Validate
	uService *service.UserService
	tService *service.TaskService
}

func New(cfg config.Config, uService *service.UserService, tService *service.TaskService) *ServerApi {
	server := http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	}
	return &ServerApi{
		server:   &server,
		valid:    validator.New(),
		uService: uService,
		tService: tService,
	}
}

func (s *ServerApi) configRoutes() {
	router := gin.Default()
	router.GET("/tasks", s.JWTMiddleware(), s.getTasks)
	router.POST("/tasks", s.createTask)
	task := router.Group("/tasks")
	{
		task.POST("/save-tasks", s.saveTasks)
		task.PUT("/:id", func(c *gin.Context) {})
		task.DELETE("/:id", func(c *gin.Context) {})
		task.GET("/:id", func(c *gin.Context) {})
	}
	users := router.Group("/users")
	{
		users.POST("/register", s.registerUser)
		users.POST("/login", s.loginUser)
	}
	s.server.Handler = router
}

func (s *ServerApi) Start() error {
	log := logger.Get()
	s.configRoutes()
	log.Info().Str("server address", s.server.Addr).Msg("server was started")
	return s.server.ListenAndServe()
}

func (s *ServerApi) JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.Get()
		token := c.GetHeader("Authorization")
		if token == `` {
			log.Error().Msg("token not found")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		uid, err := validateJwtToken(token)
		if err != nil {
			log.Error().Err(err).Msg("failed to validate token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, "invalid token")
			return
		}
		log.Debug().Str("uid", uid).Msg("user was authorized")
		c.Set("uid", uid)
		c.Next()
	}
}

func genJwtToken(uid string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   uid,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 3)),
	})
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		return ``, err
	}
	return tokenStr, nil
}

func validateJwtToken(tokenStr string) (string, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return ``, err
	}

	if !token.Valid {
		return ``, fmt.Errorf("invalid token")
	}

	return claims.Subject, nil
}
