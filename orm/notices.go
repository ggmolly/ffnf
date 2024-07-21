package orm

import (
	"fmt"
	"time"
)

type Notice struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	ButtonTitle string    `gorm:"not_null" json:"btn_title"`
	Title       string    `gorm:"not_null" json:"title"`
	Content     string    `gorm:"not_null" json:"content"`
	CreatedAt   time.Time `gorm:"autoCreateTime:true" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime:true" json:"updated_at"`
}

type GameNotice struct {
	ID          int    `json:"id"`
	Version     string `json:"version"`
	ButtonTitle string `json:"btn_title"`
	Title       string `json:"title"`
	BannerURL   string `json:"title_image"`
	Date        string `json:"time_desc"`
	Content     string `json:"content"`
	TagType     int    `json:"tag_type"`
	Icon        int    `json:"icon"`
}

func (n *Notice) ToGameNotice() GameNotice {
	return GameNotice{
		ID:          n.ID,
		Version:     fmt.Sprint(n.UpdatedAt.UnixMicro()),
		ButtonTitle: n.ButtonTitle,
		Title:       n.Title,
		BannerURL:   "https://ffnf.mana.rip/static/banner.png",
		Date:        n.CreatedAt.Format("02/01/2006"),
		Content:     n.Content,
		TagType:     2, // system category
		Icon:        1, // gear icon
	}
}
