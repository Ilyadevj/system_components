package repository

import (
	"errors"
	"system_components/internal/app/ds"
	"time"

	"gorm.io/gorm"
)

type ComponentWithLikes struct {
	ds.Component
	LikesCount int64
}

// 1. Получение опубликованных компонентов с фильтрацией по чекбоксам (ORM)
func (r *Repository) GetPublishedComponents(selectedGroups []string) ([]ComponentWithLikes, error) {
	var components []ds.Component

	tx := r.db.Where("status = ?", ds.StatusPublished)
	if len(selectedGroups) > 0 {
		tx = tx.Where("grouping IN ?", selectedGroups)
	}

	err := tx.Order("id ASC").Find(&components).Error
	if err != nil {
		return nil, err
	}

	var result []ComponentWithLikes
	for _, c := range components {
		var cnt int64
		_ = r.db.Model(&ds.ComponentLike{}).Where("component_id = ?", c.ID).Count(&cnt).Error
		result = append(result, ComponentWithLikes{
			Component:  c,
			LikesCount: cnt,
		})
	}

	return result, nil
}

// 2. Получение списка доступных слоев для чекбоксов
func (r *Repository) GetAvailableGroups() []string {
	var groups []string
	r.db.Model(&ds.Component{}).
		Where("status = ? AND grouping IS NOT NULL AND grouping != ''", ds.StatusPublished).
		Distinct("grouping").
		Pluck("grouping", &groups)
	return groups
}

func (r *Repository) GetFeedComponent(currentID int, isNext bool) (*ComponentWithLikes, error) {
	var component ds.Component
	var err error

	if isNext && currentID > 0 {
		err = r.db.Where("status = ? AND id > ?", ds.StatusPublished, currentID).
			Order("id ASC").
			Limit(1).
			First(&component).Error
	} else if currentID > 0 && !isNext {
		err = r.db.Where("status = ? AND id = ?", ds.StatusPublished, currentID).
			First(&component).Error
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || currentID == 0 {
		err = r.db.Where("status = ?", ds.StatusPublished).
			Order("id ASC").
			Limit(1).
			First(&component).Error
	}

	if err != nil {
		return nil, err
	}

	var cnt int64
	_ = r.db.Model(&ds.ComponentLike{}).Where("component_id = ?", component.ID).Count(&cnt).Error

	return &ComponentWithLikes{
		Component:  component,
		LikesCount: cnt,
	}, nil
}

func (r *Repository) GetUserDraft(creatorID uint) (*ds.Component, error) {
	var draft ds.Component
	err := r.db.Where("creator_id = ? AND status = ?", creatorID, ds.StatusDraft).First(&draft).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &draft, nil
}

func (r *Repository) CreateDraft(name, imgURL, videoURL string, creatorID uint) (*ds.Component, error) {
	draft := ds.Component{
		Name:      name,
		ImageURL:  imgURL,
		VideoURL:  videoURL,
		Status:    ds.StatusDraft,
		CreatorID: creatorID,
		CreatedAt: time.Now(),
	}
	err := r.db.Create(&draft).Error
	return &draft, err
}

func (r *Repository) PublishComponent(id uint, description, grouping string, priority int) error {
	now := time.Now()
	return r.db.Model(&ds.Component{}).
		Where("id = ? AND status = ?", id, ds.StatusDraft).
		Updates(map[string]interface{}{
			"description": description,
			"grouping":    grouping,
			"priority":    priority,
			"status":      ds.StatusPublished,
			"formed_at":   &now,
		}).Error
}