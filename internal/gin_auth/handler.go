package gin_auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db         *gorm.DB
	jwtService *JWTService
}

func NewAuthHandler(db *gorm.DB, jwtService *JWTService) *AuthHandler {
	return &AuthHandler{db: db, jwtService: jwtService}
}

type RegisterInput struct {
	Name     string `form:"name" binding:"required,min=3,max=100"`
	Email    string `form:"email" binding:"required,email,max=150"`
	Password string `form:"password" binding:"required,min=8,max=72"`
	Role     string `form:"role" binding:"required"`
}

type LoginInput struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation failed",
			"details": parseValidationError(err),
		})
		return
	}

	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	var existing User
	if err := h.db.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}

	var role Role
	if err := h.db.Where("LOWER(name) = ?", strings.ToLower(input.Role)).First(&role).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role not found"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process password"})
		return
	}

	user := User{
		Name:     strings.TrimSpace(input.Name),
		Email:    input.Email,
		Password: string(hashedPassword),
		RoleID:   role.ID,
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	token, err := h.jwtService.GenerateToken(user.ID, role.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  role.Name,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation failed",
			"details": parseValidationError(err),
		})
		return
	}

	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	var user User
	if err := h.db.Preload("Role").Where("email = ?", input.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	token, err := h.jwtService.GenerateToken(user.ID, user.Role.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role.Name,
		},
	})
}
