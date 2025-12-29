package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"koalbot_api/internal/handler"
	"koalbot_api/internal/middleware"
)

type Options struct {
	MaxBodyBytes         int64
	RateLimitRPS         float64
	RateLimitBurst       int
	LoginRateLimitRPS    float64
	LoginRateLimitBurst  int
	TrustedProxies       []string
	JWTSecret            string
	CORSAllowedOrigins   []string
	CORSAllowCredentials bool
}

func New(auth *handler.AuthHandler, v1Login *handler.V1LoginHandler, users *handler.UserHandler, masterPengguna *handler.MasterPenggunaHandler, stream *handler.StreamHandler, health *handler.HealthHandler, opts Options) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	engine.Use(cors.New(buildCORSConfig(opts)))
	engine.Use(middleware.BodyLimit(opts.MaxBodyBytes))
	engine.Use(middleware.RateLimit(opts.RateLimitRPS, opts.RateLimitBurst, 10*time.Minute))

	if err := engine.SetTrustedProxies(opts.TrustedProxies); err != nil {
		panic(err)
	}

	engine.GET("/healthz", health.Health)
	engine.GET("/status/stream", stream.StatusStream)
	engine.POST("/login", middleware.RateLimit(opts.LoginRateLimitRPS, opts.LoginRateLimitBurst, 10*time.Minute), auth.Login)
	engine.POST("/v1/login", middleware.RateLimit(opts.LoginRateLimitRPS, opts.LoginRateLimitBurst, 10*time.Minute), v1Login.Login)

	usersGroup := engine.Group("/users")
	usersGroup.Use(middleware.AuthMiddleware(opts.JWTSecret), middleware.RequireAdmin())
	usersGroup.POST("", users.Register)
	usersGroup.GET("", users.List)
	usersGroup.PUT("/:uid", users.Update)
	usersGroup.DELETE("/:uid", users.Delete)

	masterPenggunaGroup := engine.Group("/master-pengguna")
	masterPenggunaGroup.Use(middleware.AuthMiddleware(opts.JWTSecret), middleware.RequireAdmin())
	masterPenggunaGroup.POST("", masterPengguna.Create)
	masterPenggunaGroup.GET("", masterPengguna.List)
	masterPenggunaGroup.GET("/:id", masterPengguna.Get)
	masterPenggunaGroup.PUT("/:id", masterPengguna.Update)
	masterPenggunaGroup.DELETE("/:id", masterPengguna.Delete)

	dashboardGroup := engine.Group("/dashboard")
	dashboardGroup.Use(middleware.AuthMiddleware(opts.JWTSecret), middleware.RequireAdmin())
	dashboardGroup.GET("/summary", masterPengguna.Summary)

	return engine
}

func buildCORSConfig(opts Options) cors.Config {
	if len(opts.CORSAllowedOrigins) == 0 {
		return cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Content-Type", "Authorization", "Device-Id", "Device-Type", "Authorization-Token"},
			ExposeHeaders:    []string{"Authorization"},
			AllowCredentials: false,
			MaxAge:           12 * time.Hour,
		}
	}

	return cors.Config{
		AllowOrigins:     opts.CORSAllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Device-Id", "Device-Type", "Authorization-Token"},
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: opts.CORSAllowCredentials,
		MaxAge:           12 * time.Hour,
	}
}
