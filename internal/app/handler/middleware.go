package handler

import (
	"fmt"
	"net/http"
	"strings"
	"system_components/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const jwtPrefix = "Bearer "
const ctxUserIDKey = "current_user_id"
const ctxTokenKey = "current_jwt_token"

func (h *Handler) extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, jwtPrefix) {
		return strings.TrimPrefix(authHeader, jwtPrefix)
	}
	if cookieToken, err := c.Cookie("jwt_token"); err == nil && cookieToken != "" {
		return cookieToken
	}
	return ""
}

func (h *Handler) parseAndValidateToken(c *gin.Context, tokenStr string) (*ds.JWTClaims, error) {
	isBlacklisted, err := h.redis.IsJWTBlacklisted(c.Request.Context(), tokenStr)
	if err != nil {
		return nil, err
	}
	if isBlacklisted {
		return nil, fmt.Errorf("токен отозван (в blacklist)")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &ds.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(h.cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("невалидный токен")
	}

	claims, ok := token.Claims.(*ds.JWTClaims)
	if !ok {
		return nil, fmt.Errorf("ошибка чтения claims")
	}
	return claims, nil
}

// OptionalAuthMiddleware — для публичных GET-методов (работает и для гостя, и для создателя)
func (h *Handler) OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := h.extractToken(c)
		if tokenStr != "" {
			if claims, err := h.parseAndValidateToken(c, tokenStr); err == nil {
				c.Set(ctxUserIDKey, claims.UserID)
				c.Set(ctxTokenKey, tokenStr)
			}
		}
		c.Next()
	}
}

// RequireAuthMiddleware — строго для методов создателя (возвращает 401/403 для гостя)
func (h *Handler) RequireAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := h.extractToken(c)
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "требуется авторизация (отсутствует токен)"})
			return
		}

		claims, err := h.parseAndValidateToken(c, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		c.Set(ctxUserIDKey, claims.UserID)
		c.Set(ctxTokenKey, tokenStr)
		c.Next()
	}
}