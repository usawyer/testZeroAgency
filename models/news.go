package models

//type News struct {
//	Id         int    `gorm:"primaryKey" `
//	Title      string `json:"Title" validate:"required"`
//	Content    string `json:"Content" validate:"required"`
//	Categories []int  `gorm:"-" json:"Categories" `
//}
//
//type NewsCategories struct {
//	CategoryId int `gorm:"category_id"`
//	NewsId     int `gorm:"news_id"`
//}

//go:generate reform

//reform:news
type News struct {
	Id         int    `json:"Id" reform:"id,pk"`
	Title      string `json:"Title" reform:"title"`
	Content    string `json:"Content" reform:"content"`
	Categories []int  `json:"Categories" reform:"-"`
}

//reform:news_categories
type NewsCategories struct {
	Id         int `reform:"id,pk"`
	NewsId     int `reform:"news_id"`
	CategoryId int `reform:"category_id"`
}
