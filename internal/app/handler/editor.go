package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetComponentEditor(c *gin.Context) {
	draft, err := h.repo.GetUserDraft(CurrentUserID)
	if err != nil {
		h.ErrorResponse(c, http.StatusInternalServerError, "Ошибка черновика: "+err.Error())
		return
	}

	c.HTML(http.StatusOK, "system_components_editor.html", gin.H{
		"ActiveTab": "editor",
		"Draft":     draft,
		"HasDraft":  draft != nil,
	})
}

func (h *Handler) CreateDraft(c *gin.Context) {
	existingDraft, _ := h.repo.GetUserDraft(CurrentUserID)
	if existingDraft != nil {
		c.Redirect(http.StatusFound, "/system-components/editor")
		return
	}

	name := c.PostForm("name")
	imgURL := c.PostForm("image_url")
	videoURL := c.PostForm("video_url")

	if imgURL == "" {
		imgURL = "/static/img/defaultIcon.png"
	}
	if videoURL == "" {
		videoURL = "/static/video/defaultVideo.mp4"
	}

	_, _ = h.repo.CreateDraft(name, imgURL, videoURL, CurrentUserID)
	c.Redirect(http.StatusFound, "/system-components/editor")
}

func (h *Handler) PublishComponent(c *gin.Context) {
	idStr := c.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Redirect(http.StatusFound, "/system-components/editor")
		return
	}

	description := c.PostForm("description")
	grouping := c.PostForm("grouping")
	priority, _ := strconv.Atoi(c.PostForm("priority"))

	_ = h.repo.PublishComponent(uint(id), description, grouping, priority)
	c.Redirect(http.StatusFound, "/system-components")
}