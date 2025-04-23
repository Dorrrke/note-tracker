package server

import (
	"fmt"
	"net/http"

	"github.com/Dorrrke/note-tracker/gen/auth"
	"github.com/Dorrrke/note-tracker/internal/note-tracker/config"
	"github.com/Dorrrke/note-tracker/internal/note-tracker/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ServerApi struct {
	server   *http.Server
	auth     auth.AuthServiceClient
	valid    *validator.Validate
	uService *service.UserService
	tService *service.TaskService
}

func New(cfg config.Config, uService *service.UserService, tService *service.TaskService, authClient auth.AuthServiceClient) *ServerApi {
	server := http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	}
	return &ServerApi{
		server:   &server,
		auth:     authClient,
		valid:    validator.New(),
		uService: uService,
		tService: tService,
	}
}

func (s *ServerApi) configRoutes() {
	router := gin.Default()
	router.GET("/tasks", s.getTasks)
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
	s.configRoutes()
	// log := logger.Get()
	// log.Info().Str("server address", s.server.Addr).Msg("server was started")
	return s.server.ListenAndServe()
}
