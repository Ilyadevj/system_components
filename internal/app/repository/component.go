package repository

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"system_components/internal/app/ds"
	"time"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type ComponentWithLikes struct {
	ds.Component
	LikesCount int64
}

// 1. Получение опубликованных компонентов (для HTML-шаблона ЛР2)
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

// 1.1. Получение опубликованных компонентов с фильтрацией для REST API (с полем 0/1 is_creator)
func (r *Repository) GetPublishedComponentsAPI(selectedGroups []string, currentUserID uint) ([]ds.ComponentGridSerializer, error) {
	var components []ds.Component

	tx := r.db.Where("status = ?", ds.StatusPublished)
	if len(selectedGroups) > 0 {
		tx = tx.Where("grouping IN ?", selectedGroups)
	}

	err := tx.Order("id ASC").Find(&components).Error
	if err != nil {
		return nil, err
	}

	result := make([]ds.ComponentGridSerializer, 0, len(components))
	for _, c := range components {
		var cnt int64
		_ = r.db.Model(&ds.ComponentLike{}).Where("component_id = ?", c.ID).Count(&cnt).Error

		isCreator := 0
		if c.CreatorID == currentUserID {
			isCreator = 1
		}

		result = append(result, ds.ComponentGridSerializer{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
			ImageURL:    c.ImageURL,
			VideoURL:    c.VideoURL,
			Priority:    c.Priority,
			Grouping:    c.Grouping,
			LikesCount:  cnt,
			IsCreator:   isCreator,
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

// 3. Лента для HTML-шаблона ЛР2
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

// 3.1. Лента для REST API (возвращает 1 карточку + 0/1 is_liked + 0/1 is_creator)
func (r *Repository) GetFeedComponentAPI(currentID int, isNext bool, currentUserID uint) (*ds.ComponentFeedSerializer, error) {
	item, err := r.GetFeedComponent(currentID, isNext)
	if err != nil {
		return nil, err
	}

	var likeCnt int64
	_ = r.db.Model(&ds.ComponentLike{}).
		Where("component_id = ? AND user_id = ?", item.ID, currentUserID).
		Count(&likeCnt).Error

	isLiked := 0
	if likeCnt > 0 {
		isLiked = 1
	}

	isCreator := 0
	if item.CreatorID == currentUserID {
		isCreator = 1
	}

	return &ds.ComponentFeedSerializer{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		ImageURL:    item.ImageURL,
		VideoURL:    item.VideoURL,
		Priority:    item.Priority,
		Grouping:    item.Grouping,
		LikesCount:  item.LikesCount,
		IsLiked:     isLiked,
		IsCreator:   isCreator,
	}, nil
}

// 4. Получение черновика текущего пользователя (не более 1 записи, без id в запросе)
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

// 5. Создание черновика
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

// 5.1. Загрузка файла (картинки или видео) в MinIO с генерацией латинского названия
func (r *Repository) UploadMediaToMinio(componentID uint, mediaPrefix string, header *multipart.FileHeader, defaultExt string) (string, error) {
	file, err := header.Open()
	if err != nil {
		return "", fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения файла: %w", err)
	}
	contentType := http.DetectContentType(buffer[:n])

	_, err = file.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("ошибка перемещения по потоку файла: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		ext = defaultExt
	}

	// Латинское имя файла, например: component_6_img.png или component_6_video.mp4
	filename := fmt.Sprintf("component_%d_%s%s", componentID, mediaPrefix, ext)

	_, err = r.minio.PutObject(
		context.Background(),
		r.minioBucketName,
		filename,
		file,
		header.Size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return "", fmt.Errorf("не удалось загрузить объект в MinIO: %w", err)
	}

	fileURL := fmt.Sprintf("http://%s/%s/%s", r.minioEndpoint, r.minioBucketName, filename)
	return fileURL, nil
}

// 5.2. Обновление ссылок на медиа в БД после загрузки в MinIO
func (r *Repository) UpdateComponentMedia(componentID uint, imgURL, videoURL string) error {
	updates := map[string]interface{}{}
	if imgURL != "" {
		updates["image_url"] = imgURL
	}
	if videoURL != "" {
		updates["video_url"] = videoURL
	}
	if len(updates) == 0 {
		return nil
	}
	return r.db.Model(&ds.Component{}).Where("id = ?", componentID).Updates(updates).Error
}

// 6. Публикация черновика (строго из статуса draft в published, вернуть в draft нельзя)
func (r *Repository) PublishComponent(id uint, description, grouping string, priority int) error {
	now := time.Now()
	res := r.db.Model(&ds.Component{}).
		Where("id = ? AND status = ?", id, ds.StatusDraft).
		Updates(map[string]interface{}{
			"description": description,
			"grouping":    grouping,
			"priority":    priority,
			"status":      ds.StatusPublished,
			"formed_at":   &now,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// 7. Мягкое удаление через ORM (только для создателя и только не удаленных записей)
func (r *Repository) SoftDeleteComponentORM(id uint, creatorID uint) error {
	res := r.db.Model(&ds.Component{}).
		Where("id = ? AND creator_id = ? AND status != ?", id, creatorID, ds.StatusDeleted).
		Update("status", ds.StatusDeleted)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// 8. Постановка (1) или отмена (0) лайка от текущего пользователя
func (r *Repository) SetComponentLike(componentID uint, userID uint, likeValue int) (int64, error) {
	var comp ds.Component
	if err := r.db.Where("id = ? AND status = ?", componentID, ds.StatusPublished).First(&comp).Error; err != nil {
		return 0, err
	}

	if likeValue == 1 {
		var existing ds.ComponentLike
		err := r.db.Where("component_id = ? AND user_id = ?", componentID, userID).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newLike := ds.ComponentLike{
				ComponentID: componentID,
				UserID:      userID,
			}
			if err := r.db.Create(&newLike).Error; err != nil {
				return 0, err
			}
		} else if err != nil {
			return 0, err
		}
	} else {
		if err := r.db.Where("component_id = ? AND user_id = ?", componentID, userID).Delete(&ds.ComponentLike{}).Error; err != nil {
			return 0, err
		}
	}

	var cnt int64
	_ = r.db.Model(&ds.ComponentLike{}).Where("component_id = ?", componentID).Count(&cnt).Error
	return cnt, nil
}

// Получение компонента по ID
func (r *Repository) GetComponentByID(id uint) (*ds.Component, error) {
	var comp ds.Component
	if err := r.db.Where("id = ? AND status != ?", id, ds.StatusDeleted).First(&comp).Error; err != nil {
		return nil, err
	}
	return &comp, nil
}