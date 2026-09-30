package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FilterItemDTO struct {
	Name      string
	IsChecked bool
}

func (h *Handler) GetComponentsGrid(c *gin.Context) {
	selectedGroups := c.QueryArray("groups")

	items, err := h.repo.GetPublishedComponents(selectedGroups)
	if err != nil {
		h.ErrorResponse(c, http.StatusInternalServerError, "Ошибка получения каталога: "+err.Error())
		return
	}

	availableGroups := h.repo.GetAvailableGroups()
	selectedMap := make(map[string]bool)
	for _, g := range selectedGroups {
		selectedMap[g] = true
	}

	var filterItems []FilterItemDTO
	for _, g := range availableGroups {
		filterItems = append(filterItems, FilterItemDTO{
			Name:      g,
			IsChecked: selectedMap[g],
		})
	}

	c.HTML(http.StatusOK, "system_components_grid.html", gin.H{
		"ActiveTab":   "grid",
		"Components":  items,
		"FilterItems": filterItems,
	})
}

func (h *Handler) DeleteComponent(c *gin.Context) {
	idStr := c.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Redirect(http.StatusFound, "/system-components")
		return
	}

	_ = h.repo.DeleteComponentRawSQL(id)
	c.Redirect(http.StatusFound, "/system-components")
}