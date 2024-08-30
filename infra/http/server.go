package http

import (
	"context"
	"fmt"
	"liga-app/config"
	"liga-app/infra/http/middlewares"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	htmxRenderer "liga-app/infra/http/htmx/common/renderer"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

var logTag = "[http.server]"

type Server struct {
	logger *zerolog.Logger
	server *http.Server
	db     *gorm.DB
}

func NewServer(config *config.Config, logger *zerolog.Logger, db *gorm.DB) (*Server, error) {
	h, err := setupHandlers(config, db, logger)
	if err != nil {
		return nil, fmt.Errorf("http.server.NewServer(): failed to setup handlers: %w", err)
	}
	m, err := middlewares.SetupMiddlewares(config, db)
	if err != nil {
		return nil, err
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:3000",
		"https://liga.gilatennis.com",
	}
	// To be able to send tokens to the server.
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AllowCredentials = true
	// OPTIONS method for ReactJS
	corsConfig.AddAllowMethods("OPTIONS")

	rs, err := redis.NewStore(10, "tcp", config.Redis.Address, config.Redis.Password, []byte(config.Redis.Secret))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create new redis store: %w", logTag, err)
	}

	router := gin.Default()
	router.Use(cors.New(corsConfig))
	router.HTMLRender = htmxRenderer.Default

	router.Use(sessions.Sessions("liga-gilatennis-cookie-session", rs))

	setupRoutes(router, h, m)

	server := &http.Server{
		Addr:    ":" + config.Http.Port,
		Handler: router,
	}

	return &Server{
		logger: logger,
		server: server,
		db:     db,
	}, nil
}

func (s *Server) Start() {
	go func() {
		// service connections
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1000)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	log.Println("Closing database connection")
	sqlDB, err := s.db.DB()
	if err == nil {
		sqlDB.Close()
	}

	timeOut := 3 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	<-ctx.Done()
	log.Printf("timeout of %d seconds.\n", int(timeOut.Seconds()))
	log.Println("Server exiting")
}
