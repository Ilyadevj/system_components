package repository

import (
	"crypto/sha1"
	"encoding/hex"
	"system_components/internal/app/ds"
)

func HashPassword(password string) string {
	h := sha1.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func (r *Repository) RegisterUser(login, password string) (*ds.User, error) {
	user := ds.User{
		Login:    login,
		Password: HashPassword(password),
	}
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) AuthenticateUser(login, password string) (*ds.User, error) {
	var user ds.User
	hashed := HashPassword(password)

	// Ищем пользователя, у которого пароль равен либо открытой строке, либо SHA-1 хешу
	err := r.db.Where("login = ? AND (password = ? OR password = ?)", login, password, hashed).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}