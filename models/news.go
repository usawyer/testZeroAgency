package models

type News struct {
	Id         int    `gorm:"primaryKey"`
	Title      string `json:"Title"`
	Content    string `json:"Content"`
	Categories []int  `gorm:"-" json:"Categories"`
}

type NewsCategories struct {
	CategoryId int `gorm:"category_id"`
	NewsId     int `gorm:"news_id"`
}
