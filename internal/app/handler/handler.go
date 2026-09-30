package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"system_components/internal/app/repository"

	"github.com/gin-gonic/gin"
)

type ComponentCardDTO struct {
	ID         int
	Title      string
	ImageURL   string
	Grouping   string
	Priority   int
	LikesCount int
}

type GroupFilterItem struct {
	Name      string
	IsChecked bool
}

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

// GetComponentsGrid — Страница: Плитка (Каталог)
func (h *Handler) GetComponentsGrid(c *gin.Context) {
	selectedGroups := c.QueryArray("groups")

	rawItems := h.repo.GetPublished(selectedGroups)
	availableGroups := h.repo.GetAvailableGroups()

	selectedMap := make(map[string]bool)
	for _, g := range selectedGroups {
		selectedMap[g] = true
	}

	var filterItems []GroupFilterItem
	for _, g := range availableGroups {
		filterItems = append(filterItems, GroupFilterItem{
			Name:      g,
			IsChecked: selectedMap[g],
		})
	}

	var cards []ComponentCardDTO
	for _, item := range rawItems {
		cards = append(cards, ComponentCardDTO{
			ID:         item.ID,
			Title:      item.Title,
			ImageURL:   item.ImageURL,
			Grouping:   item.Grouping,
			Priority:   item.Priority,
			LikesCount: len(item.LikesUserIDs), // Подсчет количества лайков в контроллере
		})
	}

	c.HTML(http.StatusOK, "system_components_grid.html", gin.H{
		"ActiveTab":   "grid",
		"Cards":       cards,
		"FilterItems": filterItems,
	})
}

// GetComponentFeed — Страница: Лента (Vibes)
func (h *Handler) GetComponentFeed(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	isNext := c.Query("next") == "true"

	item := h.repo.GetFeedItem(id, isNext)

	var userIDs []string
	for _, uid := range item.LikesUserIDs {
		userIDs = append(userIDs, strconv.Itoa(uid))
	}

	c.HTML(http.StatusOK, "system_components_feed.html", gin.H{
		"ActiveTab":  "feed",
		"Component":  item,
		"LikesCount": len(item.LikesUserIDs),
		"LikesCSV":   strings.Join(userIDs, ","),
		"NextURL":    fmt.Sprintf("/system-components/feed?id=%d&next=true", item.ID),
	})
}

// GetComponentEditor — Страница: Добавление (Черновик)
func (h *Handler) GetComponentEditor(c *gin.Context) {
	draft := h.repo.GetDraft()
	c.HTML(http.StatusOK, "system_components_editor.html", gin.H{
		"ActiveTab": "editor",
		"Draft":     draft,
	})
}