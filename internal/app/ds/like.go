package ds

type ComponentLike struct {
	ID          uint `gorm:"primaryKey" json:"id"`
	UserID      uint `gorm:"not null;uniqueIndex:idx_user_component" json:"user_id"`
	ComponentID uint `gorm:"not null;uniqueIndex:idx_user_component" json:"component_id"`

	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:RESTRICT" json:"-"`
	Component Component `gorm:"foreignKey:ComponentID;constraint:OnDelete:RESTRICT" json:"-"`
}

func (ComponentLike) TableName() string {
	return "component_likes"
}

// Запрос на постановку/отмену лайка (0 - отменяет лайк, 1 - ставит лайк)
type LikeRequest struct {
	Like int `json:"like" form:"like"`
}