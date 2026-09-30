package main

import (
	"system_components/internal/app/config"
	"system_components/internal/app/ds"
	"system_components/internal/app/dsn"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	db, err := dsn.Connect(cfg)
	if err != nil {
		logrus.Fatalf("Ошибка подключения к БД: %v", err)
	}

	// 1. Автоматическая миграция схемы для 3 обязательных таблиц
	err = db.AutoMigrate(
		&ds.User{},
		&ds.Component{},
		&ds.ComponentLike{},
	)
	if err != nil {
		logrus.Fatalf("Ошибка миграции БД: %v", err)
	}

	logrus.Info("Таблицы успешно созданы/обновлены!")

	// 2. Сидирование начальных данных (если таблицы пустые)
	seedData(db)
}

func seedData(db *gorm.DB) {
	var count int64
	db.Model(&ds.User{}).Count(&count)
	if count > 0 {
		logrus.Info("Начальные данные уже существуют в базе данных, пропуск сидирования.")
		return
	}

	// Пользователи
	users := []ds.User{
		{ID: 1, Login: "creator_dev", Password: "secret_password1"},
		{ID: 2, Login: "admin_sec", Password: "secret_password2"},
	}
	if err := db.Create(&users).Error; err != nil {
		logrus.Warnf("Ошибка добавления пользователей: %v", err)
	}

	// Описание и поля темы для опубликованных карточек
	desc1 := "Высокопроизводительный входной шлюз для маршрутизации трафика и проверки авторизации."
	group1 := "Бэкенд"
	prio1 := 1

	desc2 := "Кластер кэширования для снижения задержек горячих запросов профилей пользователей."
	group2 := "Кэш"
	prio2 := 2

	desc3 := "Основная реляционная СУБД с шардированием таблиц истории операций."
	group3 := "БД"
	prio3 := 3

	desc4 := "Удаленный устаревший брокер очередей сообщений."
	group4 := "Очереди"
	prio4 := 2

	// Услуги (компоненты): обязательное наличие published, deleted и draft
	components := []ds.Component{
		{
			ID:          1,
			Name:        "Go API Gateway",
			Description: &desc1,
			Status:      ds.StatusPublished,
			ImageURL:    "/static/img/defaultIcon.png",
			VideoURL:    "/static/video/defaultVideo.mp4",
			Priority:    &prio1,
			Grouping:    &group1,
			CreatorID:   1,
		},
		{
			ID:          2,
			Name:        "Кэш-память Redis",
			Description: &desc2,
			Status:      ds.StatusPublished,
			ImageURL:    "/static/img/defaultIcon.png",
			VideoURL:    "/static/video/defaultVideo.mp4",
			Priority:    &prio2,
			Grouping:    &group2,
			CreatorID:   1,
		},
		{
			ID:          3,
			Name:        "Кластер PostgreSQL",
			Description: &desc3,
			Status:      ds.StatusPublished,
			ImageURL:    "/static/img/defaultIcon.png",
			VideoURL:    "/static/video/defaultVideo.mp4",
			Priority:    &prio3,
			Grouping:    &group3,
			CreatorID:   1,
		},
		{
			ID:          4,
			Name:        "Kafka Message Broker",
			Description: &desc4,
			Status:      ds.StatusDeleted,
			ImageURL:    "/static/img/defaultIcon.png",
			VideoURL:    "/static/video/defaultVideo.mp4",
			Priority:    &prio4,
			Grouping:    &group4,
			CreatorID:   1,
		},
		{
			ID:          5,
			Name:        "Elasticsearch Cluster",
			Description: nil,
			Status:      ds.StatusDraft,
			ImageURL:    "/static/img/defaultIcon.png",
			VideoURL:    "/static/video/defaultVideo.mp4",
			Priority:    nil,
			Grouping:    nil,
			CreatorID:   1,
		},
	}

	if err := db.Create(&components).Error; err != nil {
		logrus.Warnf("Ошибка добавления компонентов: %v", err)
	}

	// Лайки (м-м)
	likes := []ds.ComponentLike{
		{UserID: 1, ComponentID: 1},
		{UserID: 2, ComponentID: 1},
		{UserID: 1, ComponentID: 2},
	}
	if err := db.Create(&likes).Error; err != nil {
		logrus.Warnf("Ошибка добавления лайков: %v", err)
	}

	// Синхронизация последовательностей ID в PostgreSQL
	db.Exec("SELECT setval('system_components_id_seq', (SELECT COALESCE(MAX(id), 1) FROM system_components));")
	db.Exec("SELECT setval('users_id_seq', (SELECT COALESCE(MAX(id), 1) FROM users));")

	logrus.Info("Начальные данные успешно добавлены!")
}