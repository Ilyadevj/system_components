package ds

import "time"

type ComponentStatus string

const (
	StatusDraft     ComponentStatus = "draft"
	StatusPublished ComponentStatus = "published"
	StatusDeleted   ComponentStatus = "deleted"
)

type Component struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Name        string          `gorm:"type:varchar(100);not null" json:"name"`
	Description *string         `gorm:"type:text" json:"description"`
	Status      ComponentStatus `gorm:"type:varchar(20);not null" json:"-"`

	ImageURL string `gorm:"type:varchar(255);not null" json:"image_url"`
	VideoURL string `gorm:"type:varchar(255);not null" json:"video_url"`

	Priority *int    `gorm:"type:int;column:priority" json:"priority"`
	Grouping *string `gorm:"type:varchar(50);column:grouping" json:"grouping"`

	CreatedAt time.Time  `gorm:"autoCreateTime;not null" json:"-"`
	CreatorID uint       `gorm:"not null" json:"-"`
	FormedAt  *time.Time `gorm:"default:null" json:"-"`

	Creator User `gorm:"foreignKey:CreatorID;constraint:OnDelete:RESTRICT" json:"-"`
}

func (Component) TableName() string {
	return "system_components"
}

type ComponentGridSerializer struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	ImageURL    string  `json:"image_url"`
	VideoURL    string  `json:"video_url"`
	Priority    *int    `json:"priority"`
	Grouping    *string `json:"grouping"`
	LikesCount  int64   `json:"likes_count"`
	IsCreator   int     `json:"is_creator"`
}

type ComponentFeedSerializer struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	ImageURL    string  `json:"image_url"`
	VideoURL    string  `json:"video_url"`
	Priority    *int    `json:"priority"`
	Grouping    *string `json:"grouping"`
	LikesCount  int64   `json:"likes_count"`
	IsLiked     int     `json:"is_liked"`
	IsCreator   int     `json:"is_creator"`
}


type PublishComponentRequest struct {
	Description string `json:"description" form:"description" binding:"required"`
	Grouping    string `json:"grouping" form:"grouping" binding:"required"`
	Priority    int    `json:"priority" form:"priority" binding:"required"`
}