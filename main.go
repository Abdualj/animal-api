package main

import (
	"context"
	_ "embed"
	"log"
	"net/http"
	"time"

	"animal-api/internal/config"
	"animal-api/internal/db"
	"animal-api/internal/handlers"
	"animal-api/internal/middleware"
	"animal-api/internal/models"

	"github.com/gin-gonic/gin"
)

//go:embed docs/openapi.yaml
var openAPISpec []byte

const swaggerUI = `<!doctype html>
<html><head><meta charset="utf-8"><title>Animal API Docs</title>
<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui.css"></head>
<body><div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui-bundle.js"></script>
<script>SwaggerUIBundle({url: "/openapi.yaml", dom_id: "#swagger-ui"});</script>
</body></html>`

func main() {
	cfg := config.Load()

	client, database, err := db.Connect(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}
	defer client.Disconnect(context.Background())

	users := database.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := handlers.EnsureAdmin(ctx, users, cfg.AdminUsername, cfg.AdminPassword); err != nil {
		log.Fatalf("seed admin: %v", err)
	}
	cancel()

	authH := &handlers.AuthHandler{Users: users, JWTSecret: cfg.JWTSecret}
	animalH := &handlers.AnimalHandler{Animals: database.Collection("animals")}
	requireAuth := middleware.Authenticate(cfg.JWTSecret)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	r.GET("/openapi.yaml", func(c *gin.Context) { c.Data(http.StatusOK, "application/yaml", openAPISpec) })
	r.GET("/docs", func(c *gin.Context) { c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerUI)) })

	api := r.Group("/api")
	{
		api.POST("/auth/register", authH.Register)
		api.POST("/auth/login", authH.Login)

		api.GET("/animals", animalH.List)
		api.GET("/animals/:id", animalH.Get)
		api.POST("/animals", requireAuth, animalH.Create)
		api.PUT("/animals/:id", requireAuth, animalH.Update)
		api.DELETE("/animals/:id", requireAuth, middleware.RequireRole(models.RoleAdmin), animalH.Delete)
	}

	log.Printf("listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
