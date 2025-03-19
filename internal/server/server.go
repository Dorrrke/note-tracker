package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"honnef.co/go/tools/config"
)

type ServerApi struct {
	server *http.Server
	valid  *validator.Validate
	repo   any
}

func New(cfg config.Config, repo any) *ServerApi {
	server := http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	}
	router := gin.Default()
	router.GET("/tasks", func(c *gin.Context) {})
	router.POST("/tasks", func(c *gin.Context) {})
	task := router.Group("/tasks")
	{
		task.PUT("/:id", func(c *gin.Context) {})
		task.DELETE("/:id", func(c *gin.Context) {})
		task.GET("/:id", func(c *gin.Context) {})
	}
	user := router.Group("/user")
	{
		user.POST("/login", func(c *gin.Context) {})
		user.POST("/register", func(c *gin.Context) {})
		user.GET("/profile", func(c *gin.Context) {})
	}

	server.Handler = router
	return &ServerApi{
		server: &server,
		valid:  validator.New(),
		repo:   repo,
	}
}

func (s *ServerApi) Start() error {
	log := logger.Get()
	log.Info().Str("server address", s.server.Addr).Msg("server was started")
	return s.server.ListenAndServe()
}
