package ds

import "time"

type ComponentStatus string

const (
	StatusDraft     ComponentStatus = "draft"
	StatusPublished ComponentStatus = "published"
	StatusDeleted   ComponentStatus = "deleted"
)

type Component struct {
	ID          uint            `gorm:"primaryKey"`
	Name        string          `gorm:"type:varchar(100);not null"`
	Description *string         `gorm:"type:text"`
	Status      ComponentStatus `gorm:"type:varchar(20);not null"`

	ImageURL string `gorm:"type:varchar(255);not null"`
	VideoURL string `gorm:"type:varchar(255);not null"`

	Priority *int    `gorm:"type:int;column:priority"`
	Grouping *string `gorm:"type:varchar(50);column:grouping"`

	CreatedAt time.Time  `gorm:"autoCreateTime;not null"`
	CreatorID uint       `gorm:"not null"`
	FormedAt  *time.Time `gorm:"default:null"`

	Creator User `gorm:"foreignKey:CreatorID;constraint:OnDelete:RESTRICT"`
}

func (Component) TableName() string {
	return "system_components"
}