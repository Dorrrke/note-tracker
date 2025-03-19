package server

import (
	"fmt"
	"net/http"

	"github.com/Dorrrke/note-tracker/internal/config"
	"github.com/Dorrrke/note-tracker/internal/domain/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Repository interface {
	GetTasks() ([]models.Task, error)
	GetTask(string) (models.Task, error)
	SaveTask(models.Task) error
	UpdateTask(models.Task) error
	DeleteTask(string) error
}

type ServerApi struct {
	server *http.Server
	valid  *validator.Validate
	repo   Repository
}

func New(cfg config.Config, repo Repository) *ServerApi {
	server := http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	}
	return &ServerApi{
		server: &server,
		valid:  validator.New(),
		repo:   repo,
	}
}

func (s *ServerApi) configRoutes() {
	router := gin.Default()
	router.GET("/tasks", s.getTasks)
	router.POST("/tasks", s.createTask)
	task := router.Group("/tasks")
	{
		task.PUT("/:id", func(c *gin.Context) {})
		task.DELETE("/:id", func(c *gin.Context) {})
		task.GET("/:id", func(c *gin.Context) {})
	}
	s.server.Handler = router
}

func (s *ServerApi) Start() error {
	s.configRoutes()
	// log := logger.Get()
	// log.Info().Str("server address", s.server.Addr).Msg("server was started")
	return s.server.ListenAndServe()
}
