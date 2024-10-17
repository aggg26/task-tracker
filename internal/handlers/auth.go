package handlers

import (
	"errors"
	"net/http"
	"trackerApp/internal/services/dtos"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware godoc
// @Summary Authorization middleware
// @Description Middleware для проверки заголовка Authorization и валидации JWT токена.
// @Tags middleware
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{"user_id": int}
// @Failure 401 {object} gin.H{"error": string}
// @Router /auth [get]
func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Отсутствует заголовок Authorization"})
			c.Abort()
			return
		}

		tokenString := authHeader[len("Bearer "):]

		userId, err := h.services.IAuthService.ParseJwt(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("user_id", userId)
		c.Next()
	}
}

func (h *Handler) GetUserId(c *gin.Context) (int, error) {
	id, ok := c.Get("user_id")
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id is invalid type")
	}
	return idInt, nil
}

// GetUserId godoc
// @Summary Retrieve user ID from context
// @Description Получает ID пользователя из контекста, который был установлен в middleware авторизации.
// @Tags helpers
// @Produce json
// @Success 200 {object} int "User ID"
// @Failure 400 {object} gin.H{"error": string}
func (h *Handler) SignIn(c *gin.Context) {
	var request dtos.UserForm
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid body request")
		return
	}
	accessToken, err := h.services.IAuthService.GenerateJwt(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": accessToken})
}

// SignUp godoc
// @Summary User registration
// @Description Register a new user with the provided user information.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.UserForm true "User registration form"
// @Success 200 {object} gin.H{"message": string} "Registration success message"
// @Failure 400 {string} string "Invalid body request"
// @Failure 500 {object} gin.H{"error": string}
// @Router /signup [post]
func (h *Handler) SignUp(c *gin.Context) {
	var request dtos.UserForm
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid body request")
		return
	}
	if _, err := h.services.IAuthService.AddUser(request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "success registration"})
}

// Logout godoc
// @Summary User logout
// @Description Logs out the user by clearing the Authorization header.
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{"message": string} "Logout success message"
// @Router /logout [post]
func (h *Handler) Logout(c *gin.Context) {
	c.Header("Authorization", "")
	c.JSON(200, gin.H{"message": "Goodbye!"})
}
