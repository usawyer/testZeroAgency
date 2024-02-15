package service

import "github.com/usawyer/testZeroAgency/models"

type Store interface {
	CreatePost(news models.News) error
	GetPosts(params models.SearchParams) ([]models.News, error)
	EditPost(id int, news models.News)
	IfExists(id int) bool
}
