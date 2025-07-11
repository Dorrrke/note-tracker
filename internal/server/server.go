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

const readHeaderTimeout = 5
const writeTimeout = 10
const idleTimeout = 30

type API struct {
	server   *http.Server
	valid    *validator.Validate
	uService *service.UserService
	tService *service.TaskService
}

func New(cfg config.Config, uService *service.UserService, tService *service.TaskService) *API {
	server := http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		ReadHeaderTimeout: readHeaderTimeout * time.Second,
		WriteTimeout:      writeTimeout * time.Second,
		IdleTimeout:       idleTimeout * time.Second,
	}
	return &API{
		server:   &server,
		valid:    validator.New(),
		uService: uService,
		tService: tService,
	}
}

func (s *API) configRoutes() {
	// Task routers
	router := gin.Default()
	router.GET("/tasks", s.getTasks)
	router.POST("/tasks", s.createTask)
	task := router.Group("/tasks")
	{
		task.PUT("/:id", s.updateTaskByID)
		task.DELETE("/:id", s.deleteTaskByID)
		task.GET("/:id", s.getTaskByID)
	}

	// User routers
	router.GET("/users", s.getUsers)
	router.POST("/users", s.createUser)
	users := router.Group("/users")
	{
		users.PUT("/:id", s.updateUserByID)
		users.DELETE("/:id", s.deleteUserByID)
		users.GET("/:id", s.getUserByID)

		users.GET("/profile", s.getUserProfile)

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
