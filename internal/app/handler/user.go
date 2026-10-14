package handler

import (
	"fmt"
	"net/http"
	"system_components/internal/app/ds"

	"github.com/gin-gonic/gin"
)

// 8. POST /api/users/register — регистрация нового пользователя
func (h *Handler) RegisterUserAPI(c *gin.Context) {
	var req ds.UserAuthRequest
	if err := c.ShouldBind(&req); err != nil {
		h.errorHandlerAPI(c, http.StatusBadRequest, err)
		return
	}

	user, err := h.repo.RegisterUser(req.Login, req.Password)
	if err != nil {
		h.errorHandlerAPI(c, http.StatusConflict, fmt.Errorf("пользователь с таким логином уже существует"))
		return
	}

	c.JSON(http.StatusCreated, user)
}

// 9. POST /api/users/login — аутентификация (заглушка для 4-й ЛР)
func (h *Handler) LoginUserAPI(c *gin.Context) {
	var req ds.UserAuthRequest
	if err := c.ShouldBind(&req); err != nil {
		h.errorHandlerAPI(c, http.StatusBadRequest, err)
		return
	}

	user, err := h.repo.AuthenticateUser(req.Login, req.Password)
	if err != nil {
		h.errorHandlerAPI(c, http.StatusUnauthorized, fmt.Errorf("неверный логин или пароль"))
		return
	}

	c.JSON(http.StatusOK, user)
}

// 10. POST /api/users/logout — деавторизация (заглушка для 4-й ЛР)
func (h *Handler) LogoutUserAPI(c *gin.Context) {
	userID := GetCurrentUserID()
	c.JSON(http.StatusOK, gin.H{
		"user_id":       userID,
		"authenticated": 0,
	})
}