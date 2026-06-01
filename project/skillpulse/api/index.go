package handler

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/trainwithshubham/skillpulse/database"
	"github.com/trainwithshubham/skillpulse/handlers"
)

var app *gin.Engine

func init() {
	// Set Gin to release mode when running in production Vercel environments
	gin.SetMode(gin.ReleaseMode)

	app = gin.New()
	app.Use(gin.Recovery())

	// CORS Middleware to facilitate API testing and third-party integrations
	app.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// API endpoints
	api := app.Group("/api")
	{
		api.GET("/skills", handlers.GetSkills)
		api.POST("/skills", handlers.CreateSkill)
		api.GET("/skills/:id", handlers.GetSkill)
		api.DELETE("/skills/:id", handlers.DeleteSkill)
		api.POST("/skills/:id/log", handlers.CreateLog)
		api.GET("/dashboard", handlers.GetDashboard)
	}

	// Health check endpoint
	app.GET("/health", handlers.HealthCheck)
}

// Handler is the entrypoint required by Vercel serverless functions for Go.
func Handler(w http.ResponseWriter, r *http.Request) {
	// Lazy-load database connection. In serverless, functions can stay warm, 
	// so we only connect if the DB connection is not initialized yet.
	if database.DB == nil {
		if err := database.Connect(); err != nil {
			log.Printf("Vercel Serverless Database Connection Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Database Connection Error. Please verify Vercel environment variables. Details: " + err.Error()))
			return
		}
	}

	// Serve the HTTP request using the Gin router
	app.ServeHTTP(w, r)
}
