package main

import (
	"log"
	"os"

	"book-coffee-shop/internal/config"
	"book-coffee-shop/internal/gin_auth"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	dsn := config.DefaultPostgresConfig().DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&gin_auth.Role{}, &gin_auth.User{}); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("Migrations applied successfully")

	seedRoles(db)

	jwtSecret := config.JWTSecret()
	if s := os.Getenv("JWT_SECRET"); s != "" {
		jwtSecret = s
	}
	jwtService := gin_auth.NewJWTService(jwtSecret)
	authHandler := gin_auth.NewAuthHandler(db, jwtService)

	r := gin.Default()

	r.POST("/register/alexandria", authHandler.Register)
	r.POST("/login/alexandria", authHandler.Login)

	protected := r.Group("/api")
	protected.Use(gin_auth.AuthMiddleware(jwtService))
	{
		protected.GET("/profile", getProfile)
		protected.GET("/admin", gin_auth.RequireRole("Admin"), adminPanel)
		protected.GET("/editor", gin_auth.RequireRole("Admin", "Editor"), editorPanel)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Gin server listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func seedRoles(db *gorm.DB) {
	roles := []gin_auth.Role{
		{Name: "Admin"},
		{Name: "Editor"},
		{Name: "User"},
	}
	for _, role := range roles {
		db.Where("name = ?", role.Name).FirstOrCreate(&role)
	}
	log.Println("Default roles seeded")
}

func getProfile(c *gin.Context) {
	userID, _ := c.Get(gin_auth.ContextKeyUserID)
	role, _ := c.Get(gin_auth.ContextKeyRole)
	c.JSON(200, gin.H{
		"user_id": userID,
		"role":    role,
		"message": "authenticated user profile",
	})
}

func adminPanel(c *gin.Context) {
	c.JSON(200, gin.H{"message": "welcome admin, you have full access"})
}

func editorPanel(c *gin.Context) {
	c.JSON(200, gin.H{"message": "welcome editor, you can manage content"})
}
