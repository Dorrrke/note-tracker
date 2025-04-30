package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Dorrrke/note-tracker/internal/config"
	"github.com/Dorrrke/note-tracker/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const rdTimeout = 15 * time.Second

type API struct {
	server   *http.Server
	valid    *validator.Validate
	uService *service.UserService
	tService *service.TaskService
}

func New(cfg config.Config, uService *service.UserService, tService *service.TaskService) *API {
	server := http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		ReadHeaderTimeout: rdTimeout,
	}
	return &API{
		server:   &server,
		valid:    validator.New(),
		uService: uService,
		tService: tService,
	}
}

func (s *API) configRoutes() {
	router := gin.Default()
	router.GET("/tasks", s.getTasks)
	router.POST("/tasks", s.createTask)
	task := router.Group("/tasks")
	{
		task.POST("/save-tasks", s.saveTasks)
		task.PUT("/:id", func(_ *gin.Context) {})
		task.DELETE("/:id", func(_ *gin.Context) {})
		task.GET("/:id", func(_ *gin.Context) {})
	}
	users := router.Group("/users")
	{
		users.POST("/register", s.registerUser)
		users.POST("/login", s.loginUser)
	}
	s.server.Handler = router
}

func (s *API) Start() error {
	s.configRoutes()
	// log := logger.Get()
	// log.Info().Str("server address", s.server.Addr).Msg("server was started")
	return s.server.ListenAndServe()
}
