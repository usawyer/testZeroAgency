package service

import (
	"github.com/usawyer/testZeroAgency/models"
)

type Service struct {
	DB Store
}

func New(db Store) *Service {
	return &Service{DB: db}
}

func (s *Service) CreatePost(news models.News) error {
	err := s.DB.CreatePost(news)
	//fmt.Println(err)
	return err
}

func (s *Service) GetPosts(params models.SearchParams) ([]models.News, error) {
	res, err := s.DB.GetPosts(params)
	return res, err
}

func (s *Service) EditPost(id int, news models.News) error {
	s.DB.EditPost(id, news)
	//обработка ошибок
	return nil
}

func (s *Service) IfExists(id int) bool {
	return s.DB.IfExists(id)
}
