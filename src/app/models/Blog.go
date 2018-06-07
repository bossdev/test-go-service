package models

import (
	"app/config"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/go-ozzo/ozzo-validation"
)

// "app/config"
// "github.com/jinzhu/gorm"

type Blog struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	UserId      int           `bson:"user_id" json:"userId" form:"userId"`
	Slug        string        `bson:"slug" json:"slug" form:"slug"`
	TitleTh     string        `bson:"title_th" json:"titleTh" form:"titleTh"`
	TitleEn     string        `bson:"title_en" json:"titleEn" form:"titleEn"`
	ContentTh   string        `bson:"content_th" json:"contentTh" form:"contentTh"`
	ContentEn   string        `bson:"content_en" json:"contentEn" form:"contentEn"`
	Format      string        `bson:"format" json:"format"`
	Type        string        `bson:"type" json:"type"`
	Status      string        `bson:"status" json:"status" form:"status"`
	CountView   int           `bson:"count_view" json:"countView"`
	View        int           `bson:"view" json:"view"`
	ThumbsL     string        `bson:"thumbs_l" json:"thumbsL"`
	Cover       string        `bson:"cover" json:"cover"`
	ThumbsS     string        `bson:"thumbs_s" json:"thumbsS"`
	TitleCover  string        `bson:"title_cover" json:"titleCover"`
	CreatedAt   time.Time     `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time     `bson:"updated_at" json:"updatedAt"`
	DeletedAt   *time.Time    `bson:"deleted_at" json:"deletedAt"`
	PublishedAt *time.Time    `bson:"published_at" json:"publishedAt"`
}

func SetCurrentTime(blog *Blog, lists ...string) {
	createdFormat := "2006-01-02 15:04:05"
	now := time.Now()
	currentDate := now.In(config.Get().LocalTimeZone).Format(createdFormat)
	t, _ := time.Parse(createdFormat, currentDate)
	for _, list := range lists {
		switch list {
		case "CreatedAt":
			blog.CreatedAt = t
		case "UpdatedAt":
			blog.UpdatedAt = t
		case "PublishedAt":
			blog.PublishedAt = &t
		case "DeletedAt":
			// cur := currentDate
			// blog.DeletedAt = &cur
			blog.DeletedAt = &t
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
