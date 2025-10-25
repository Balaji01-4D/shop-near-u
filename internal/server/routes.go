package server

import (
	"net/http"
	productcatlog "shop-near-u/internal/productCatlog"
	"shop-near-u/internal/shop"
	"shop-near-u/internal/user"
	"shop-near-u/internal/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	user.RegisterRoutes(r, s.db.GetDB())
	shop.RegisterRoutes(r, s.db.GetDB())
	productcatlog.RegisterRoutes(r, s.db.GetDB())

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Hello World", nil)
}

func (s *Server) healthHandler(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Service is healthy", s.db.Health())
}
