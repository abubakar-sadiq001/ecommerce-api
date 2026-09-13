package server

import (
	"net/http"

	"github.com/abubakar-sadiq001/learning-go-shop/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Server struct {
	config *config.Config
	db     *gorm.DB
	logger zerolog.Logger
}

func New(cfg *config.Config, db *gorm.DB, logger *zerolog.Logger) *Server {
	return &Server{
		config: cfg,
		db:     db,
		logger: *logger,
	}
}

func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.New()

	// middlewares
	// router.Use(gin.Logger())
	// router.Use(gin.Recovery())
	// router.Use(s.corsMiddleware())
	router.Use(gin.Logger(), gin.Recovery(), s.corsMiddleware())

	// routes
	router.GET("/health", s.healthCheck)

	api := router.Group("/api/v1")
	{ // nolint:gocritic
		auth := api.Group("/auth")
		{
			auth.POST("/register", s.register)
			auth.POST("/login", s.login)
			auth.POST("/refresh", s.refreshToken)
			auth.POST("/logout", s.logout)
		}

		/* creates a sibling group off api at the same URL prefix (/api/v1/) but with authMiddleware tacked on. From that point forward, any subgroup or route registered under protected inherits the auth middleware, while anything registered directly under api does not. */
		protected := api.Group("/")       // copies [Logger, Recovery, cors]
		protected.Use(s.authMiddleware()) // appends to protected's slice ONLY
		{
			users := protected.Group("/users")
			{
				userRoutes := users
				userRoutes.GET("/profile", s.getProfile)
				userRoutes.GET("/update_profile", s.updateProfile)
			}
		}
	}

	return router
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Method", "GET, POST, PUT DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, ")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}

}
