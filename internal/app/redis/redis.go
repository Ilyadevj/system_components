package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const servicePrefix = "system_components."
const jwtBlacklistPrefix = "jwt_blacklist."
const activeSessionPrefix = "active_user."

type Client struct {
	client *redis.Client
}

func NewClient(ctx context.Context, addr, password string) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("ошибка подключения к Redis: %w", err)
	}

	return &Client{client: rdb}, nil
}

func getBlacklistKey(token string) string {
	return servicePrefix + jwtBlacklistPrefix + token
}

func getSessionKey(userID uint) string {
	return fmt.Sprintf("%s%s%d", servicePrefix, activeSessionPrefix, userID)
}

// Сохранение активной сессии пользователя в Redis (чтобы показать в cmd через KEYS * и Lua скрипт)
func (c *Client) SaveActiveUser(ctx context.Context, userID uint, login string, ttl time.Duration) error {
	key := getSessionKey(userID)
	return c.client.HSet(ctx, key, "user_id", userID, "login", login).Err()
}

// Запись JWT-токена в Blacklist при Logout
func (c *Client) WriteJWTToBlacklist(ctx context.Context, token string, ttl time.Duration) error {
	return c.client.Set(ctx, getBlacklistKey(token), "revoked", ttl).Err()
}

// Проверка наличия JWT-токена в Blacklist
func (c *Client) IsJWTBlacklisted(ctx context.Context, token string) (bool, error) {
	err := c.client.Get(ctx, getBlacklistKey(token)).Err()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}