package ds

type ComponentLike struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint `gorm:"not null;uniqueIndex:idx_user_component"`
	ComponentID uint `gorm:"not null;uniqueIndex:idx_user_component"`

	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:RESTRICT"`
	Component Component `gorm:"foreignKey:ComponentID;constraint:OnDelete:RESTRICT"`
}

func (ComponentLike) TableName() string {
	return "component_likes"
}