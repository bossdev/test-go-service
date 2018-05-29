package models

import (
	"app/config"
	"time"

	"github.com/go-ozzo/ozzo-validation"
)

// "time"
// "app/config"
// "github.com/jinzhu/gorm"

type Blog struct {
	// gorm.Model
	ID          uint    `gorm:"primary_key;AUTO_INCREMENT;column:id"`
	UserId      int     `gorm:"column:user_id"`
	Slug        string  `gorm:"column:slug"`
	TitleTh     string  `gorm:"column:title_th"`
	TitleEn     string  `gorm:"column:title_en"`
	ContentTh   string  `gorm:"column:content_th"`
	ContentEn   string  `gorm:"column:content_en"`
	Format      string  `gorm:"default:markdown;column:format"`
	Type        string  `gorm:"default:post;column:type"`
	Status      string  `gorm:"default:draft;column:status"`
	CountView   int     `gorm:"default:draft;column:count_view"`
	View        int     `gorm:"column:view"`
	ThumbsL     string  `gorm:"column:thumbs_l"`
	Cover       string  `gorm:"column:cover"`
	ThumbsS     string  `gorm:"column:thumbs_s"`
	TitleCover  string  `gorm:"type:varchar(100);column:title_cover"`
	CreatedAt   string  `gorm:"column:created_at"`
	UpdatedAt   string  `gorm:"column:updated_at"`
	DeletedAt   *string `gorm:"column:deleted_at;"`
	PublishedAt string  `gorm:"column:published_at"`
	// CreatedAt time.Time `json:"created_at"`
	// UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt *time.Time `json:"deleted_at"`
	// PublishedAt time.Time `json:"published_at"`
}

func (b *Blog) TableName() string {
	return "blog_contents"
}

func SetCurrentTime(blog *Blog, lists ...string) {
	createdFormat := "2006-01-02 15:04:05"
	now := time.Now()
	currentDate := now.In(config.Get().LocalTimeZone).Format(createdFormat)
	for _, list := range lists {
		switch list {
		case "CreatedAt":
			blog.CreatedAt = currentDate
		case "UpdatedAt":
			blog.UpdatedAt = currentDate
		case "PublishedAt":
			blog.PublishedAt = currentDate
		case "DeletedAt":
			cur := currentDate
			blog.DeletedAt = &cur
		}
	}
	// t, _ := time.Parse(createdFormat , currentDate) <- set string to time
}

func (b Blog) Validate() error {
	validate := validation.ValidateStruct(&b,
		validation.Field(&b.UserId,
			validation.Required.Error("UserId is required"),
		),
		validation.Field(&b.Slug,
			validation.Required.Error("Slug is required"),
			validation.Length(3, 100).Error("The slug length must be between 3 and 100"),
		),
		validation.Field(&b.TitleTh,
			validation.Required.Error("TitleTh is required"),
		),
		validation.Field(&b.TitleEn,
			validation.Required.Error("TitleEn is required"),
		),
		validation.Field(&b.ContentTh,
			validation.Required.Error("ContentTh is required"),
		),
		validation.Field(&b.ContentEn,
			validation.Required.Error("ContentEn is required"),
		),
		validation.Field(&b.Status,
			validation.Required.Error("Status is required"),
			validation.In("draft", "publish").Error("Status is optional, and should be either 'draft' or 'publish'"),
		),
	)
	return validate
}
