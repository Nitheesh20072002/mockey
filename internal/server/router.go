package server

import (
	"net/http"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/mockey/internal/server/handlers"
	"github.com/mockey/internal/server/middleware"
)

// SetupRoutes registers routes on the provided Gin engine.
func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		auth := api.Group("/auth")
		// apply a lightweight rate limiter: 5 requests per minute per IP+route
		auth.Use(middleware.NewRateLimiter(5, time.Minute))
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			auth.POST("/otp", func(c *gin.Context) { c.JSON(200, gin.H{"msg": "otp - TODO"}) })
			auth.GET("/me", middleware.JWTAuth(), handlers.Me)
			auth.POST("/refresh", handlers.RefreshToken)
		}

		// Admin, tests, submissions and analytics routes are left as placeholders for now.
		admin := api.Group("/admin")
		{
			admin.POST("/exams", middleware.JWTAuth(), handlers.CreateExam)
			admin.POST("/questions/upload", handlers.UploadQuestions)
			admin.GET("/upload-jobs/:id", func(c *gin.Context) { c.JSON(200, gin.H{"msg": "job status - TODO"}) })
			admin.GET("/questions/:exam_id", func(c *gin.Context) { c.JSON(200, gin.H{"msg": "list questions - TODO"}) })
		}

		// tests := api.Group("/tests")
		// {
		// 	tests.GET("", func(c *gin.Context) { c.JSON(200, gin.H{"msg": "list tests - TODO"}) })
		// 	tests.POST(":test_id/start", func(c *gin.Context) { c.JSON(200, gin.H{"msg": "start test - TODO"}) })
		// }

		// sub := api.Group("/submissions")
		// {
		// 	sub.POST(":submission_id/answer", func(c *gin.Context) { c.JSON(200, gin.H{"msg": "answer - TODO"}) })
		// 	sub.POST(":submission_id/finish", func(c *gin.Context) { c.JSON(200, gin.H{"msg": "finish - TODO"}) })
		// 	sub.GET(":submission_id/result", func(c *gin.Context) { c.JSON(200, gin.H{"msg": "result - TODO"}) })
		// }

		// 	analytics := api.Group("/analytics")
		// 	{
		// 		analytics.GET("/users/:id/stats", func(c *gin.Context) { c.JSON(200, gin.H{"msg": "user stats - TODO"}) })
		// 		analytics.GET("/tests/:test_id/leaderboard", func(c *gin.Context) { c.JSON(200, gin.H{"msg": "leaderboard - TODO"}) })
		// 	}
	}
}
