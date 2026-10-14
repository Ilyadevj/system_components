package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"system_components/internal/app/ds"

	"github.com/gin-gonic/gin"
)

// HTML-обработчик (ЛР2)
func (h *Handler) GetComponentFeed(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	isNext := c.Query("next") == "true"

	item, err := h.repo.GetFeedComponent(id, isNext)
	if err != nil || item == nil {
		h.ErrorResponse(c, http.StatusNotFound, "Опубликованные компоненты отсутствуют")
		return
	}

	videoURL := item.VideoURL
	if videoURL == "" {
		videoURL = "/static/video/defaultVideo.mp4"
	}

	priorityVal := 0
	if item.Priority != nil {
		priorityVal = *item.Priority
	}

	groupingVal := "—"
	if item.Grouping != nil {
		groupingVal = *item.Grouping
	}

	descVal := ""
	if item.Description != nil {
		descVal = *item.Description
	}

	c.HTML(http.StatusOK, "system_components_feed.html", gin.H{
		"ActiveTab":   "feed",
		"ID":          item.ID,
		"Name":        item.Name,
		"Description": descVal,
		"VideoURL":    videoURL,
		"Priority":    priorityVal,
		"Grouping":    groupingVal,
		"LikesCount":  item.LikesCount,
		"NextURL":     fmt.Sprintf("/system-components/feed?id=%d&next=true", item.ID),
	})
}

// ==================== МЕТОДЫ REST API (ЛР3) ====================

// 2. GET /api/system-components/feed — получение 1 карточки ленты с признаками 0/1 (is_liked и is_creator)
func (h *Handler) GetComponentFeedAPI(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	isNext := c.Query("next") == "true"

	userID := GetCurrentUserID()
	item, err := h.repo.GetFeedComponentAPI(id, isNext, userID)
	if err != nil || item == nil {
		h.errorHandlerAPI(c, http.StatusNotFound, fmt.Errorf("опубликованные записи не найдены"))
		return
	}

	c.JSON(http.StatusOK, item)
}

// 7. POST /api/system-components/:id/like — постановка (1) или отмена (0) лайка от текущего пользователя
func (h *Handler) LikeComponentAPI(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		h.errorHandlerAPI(c, http.StatusBadRequest, fmt.Errorf("некорректный id"))
		return
	}

	var req ds.LikeRequest
	if err := c.ShouldBind(&req); err != nil {
		h.errorHandlerAPI(c, http.StatusBadRequest, err)
		return
	}

	if req.Like != 0 && req.Like != 1 {
		h.errorHandlerAPI(c, http.StatusBadRequest, fmt.Errorf("поле like должно быть 0 или 1"))
		return
	}

	userID := GetCurrentUserID()
	likesCount, err := h.repo.SetComponentLike(uint(id), userID, req.Like)
	if err != nil {
		h.errorHandlerAPI(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"component_id": id,
		"is_liked":     req.Like,
		"likes_count":  likesCount,
	})
}