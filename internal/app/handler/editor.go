package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"system_components/internal/app/ds"

	"github.com/gin-gonic/gin"
)

// HTML-обработчик (ЛР2)
func (h *Handler) GetComponentEditor(c *gin.Context) {
	userID := GetCurrentUserID()
	draft, err := h.repo.GetUserDraft(userID)
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

// HTML-обработчик (ЛР2)
func (h *Handler) CreateDraft(c *gin.Context) {
	userID := GetCurrentUserID()
	existingDraft, _ := h.repo.GetUserDraft(userID)
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

	_, _ = h.repo.CreateDraft(name, imgURL, videoURL, userID)
	c.Redirect(http.StatusFound, "/system-components/editor")
}

// HTML-обработчик (ЛР2)
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

// ==================== МЕТОДЫ REST API (ЛР3) ====================

// 3. GET /api/system-components/draft — получение черновика текущего пользователя (без id в URL)
func (h *Handler) GetDraftAPI(c *gin.Context) {
	userID := GetCurrentUserID()
	draft, err := h.repo.GetUserDraft(userID)
	if err != nil {
		h.errorHandlerAPI(c, http.StatusInternalServerError, err)
		return
	}
	if draft == nil {
		h.errorHandlerAPI(c, http.StatusNotFound, fmt.Errorf("черновик отсутствует"))
		return
	}

	c.JSON(http.StatusOK, draft)
}

// 4. POST /api/system-components — создание черновика (только name + файлы image и video в MinIO)
func (h *Handler) CreateDraftAPI(c *gin.Context) {
	userID := GetCurrentUserID()

	existingDraft, err := h.repo.GetUserDraft(userID)
	if err != nil {
		h.errorHandlerAPI(c, http.StatusInternalServerError, err)
		return
	}
	if existingDraft != nil {
		h.errorHandlerAPI(c, http.StatusConflict, fmt.Errorf("у пользователя уже есть активный черновик с id=%d", existingDraft.ID))
		return
	}

	// Читаем до 20 МБ multipart-формы (для картинки и короткого видео)
	if err := c.Request.ParseMultipartForm(20 << 20); err != nil {
		h.errorHandlerAPI(c, http.StatusBadRequest, err)
		return
	}

	name := c.Request.FormValue("name")
	if name == "" {
		h.errorHandlerAPI(c, http.StatusBadRequest, fmt.Errorf("поле name обязательно"))
		return
	}

	imgHeader, imgErr := c.FormFile("image")
	videoHeader, videoErr := c.FormFile("video")

	allowedImages := []string{"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp"}
	allowedVideos := []string{"video/mp4", "video/webm", "application/octet-stream"}

	if imgErr == nil && imgHeader != nil {
		if code, err := validateMediaUpload(imgHeader, allowedImages); err != nil {
			h.errorHandlerAPI(c, code, err)
			return
		}
	}

	if videoErr == nil && videoHeader != nil {
		if code, err := validateMediaUpload(videoHeader, allowedVideos); err != nil {
			h.errorHandlerAPI(c, code, err)
			return
		}
	}

	defaultImg := "/static/img/defaultIcon.png"
	defaultVideo := "/static/video/defaultVideo.mp4"

	draft, err := h.repo.CreateDraft(name, defaultImg, defaultVideo, userID)
	if err != nil {
		h.errorHandlerAPI(c, http.StatusInternalServerError, err)
		return
	}

	var uploadedImgURL, uploadedVideoURL string

	if imgErr == nil && imgHeader != nil {
		uploadedImgURL, err = h.repo.UploadMediaToMinio(draft.ID, "img", imgHeader, ".png")
		if err != nil {
			h.errorHandlerAPI(c, http.StatusInternalServerError, err)
			return
		}
	}

	if videoErr == nil && videoHeader != nil {
		uploadedVideoURL, err = h.repo.UploadMediaToMinio(draft.ID, "video", videoHeader, ".mp4")
		if err != nil {
			h.errorHandlerAPI(c, http.StatusInternalServerError, err)
			return
		}
	}

	if uploadedImgURL != "" || uploadedVideoURL != "" {
		if err := h.repo.UpdateComponentMedia(draft.ID, uploadedImgURL, uploadedVideoURL); err != nil {
			h.errorHandlerAPI(c, http.StatusInternalServerError, err)
			return
		}
		draft, _ = h.repo.GetComponentByID(draft.ID)
	}

	c.JSON(http.StatusCreated, draft)
}

// 5. PUT /api/system-components/:id — публикация черновика (смена статуса и заполнение полей темы)
func (h *Handler) PublishComponentAPI(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		h.errorHandlerAPI(c, http.StatusBadRequest, fmt.Errorf("некорректный id"))
		return
	}

	var req ds.PublishComponentRequest
	if err := c.ShouldBind(&req); err != nil {
		h.errorHandlerAPI(c, http.StatusBadRequest, err)
		return
	}

	if err := h.repo.PublishComponent(uint(id), req.Description, req.Grouping, req.Priority); err != nil {
		h.errorHandlerAPI(c, http.StatusBadRequest, fmt.Errorf("не удалось опубликовать запись (возможно, она не в статусе черновика)"))
		return
	}

	updated, err := h.repo.GetComponentByID(uint(id))
	if err != nil {
		h.errorHandlerAPI(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, updated)
}