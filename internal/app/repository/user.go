package repository

import "system_components/internal/app/ds"

func (r *Repository) RegisterUser(login, password string) (*ds.User, error) {
	user := ds.User{
		Login:    login,
		Password: password,
	}
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) AuthenticateUser(login, password string) (*ds.User, error) {
	var user ds.User
	err := r.db.Where("login = ? AND password = ?", login, password).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}