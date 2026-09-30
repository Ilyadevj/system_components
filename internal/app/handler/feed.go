package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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