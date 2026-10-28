package handler

import (
	"fmt"
	"net/http"
	"system_components/internal/app/ds"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// RegisterUserAPI godoc
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя в БД с хешированием пароля
// @Tags Users
// @Accept json
// @Produce json
// @Param request body ds.UserAuthRequest true "Данные регистрации"
// @Success 201 {object} ds.User
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/users/register [post]
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

// LoginUserAPI godoc
// @Summary Аутентификация пользователя
// @Description Проверяет логин/пароль, возвращает JWT-токен, устанавливает Cookie и записывает сессию в Redis
// @Tags Users
// @Accept json
// @Produce json
// @Param request body ds.UserAuthRequest true "Учетные данные"
// @Success 200 {object} ds.UserLoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/users/login [post]
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

	now := time.Now()
	claims := &ds.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(h.cfg.JWTExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "system_components_backend",
		},
		UserID: user.ID,
		Login:  user.Login,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		h.errorHandlerAPI(c, http.StatusInternalServerError, fmt.Errorf("ошибка генерации токена"))
		return
	}

	// Сохраняем информацию о пользователе в Redis для демонстрации через redis-cli и Lua
	_ = h.redis.SaveActiveUser(c.Request.Context(), user.ID, user.Login, h.cfg.JWTExpiresIn)

	// Устанавливаем куку на случай проверки через вкладку Application -> Cookies
	c.SetCookie("jwt_token", tokenStr, int(h.cfg.JWTExpiresIn.Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, ds.UserLoginResponse{
		AccessToken: tokenStr,
		TokenType:   "Bearer",
		ExpiresIn:   int64(h.cfg.JWTExpiresIn.Seconds()),
		User:        *user,
	})
}

// LogoutUserAPI godoc
// @Summary Деавторизация пользователя (Logout)
// @Description Добавляет текущий JWT-токен в Blacklist Redis и очищает куку
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /api/users/logout [post]
func (h *Handler) LogoutUserAPI(c *gin.Context) {
	userID := GetCurrentUserID(c)
	tokenStr := h.extractToken(c)

	if tokenStr != "" {
		_ = h.redis.WriteJWTToBlacklist(c.Request.Context(), tokenStr, h.cfg.JWTExpiresIn)
	}
	c.SetCookie("jwt_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"user_id":       userID,
		"authenticated": 0,
	})
}