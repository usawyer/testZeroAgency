package database

import (
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/usawyer/testZeroAgency/internal/service"
	"github.com/usawyer/testZeroAgency/models"
	"gopkg.in/reform.v1/dialects/postgresql"

	"gopkg.in/reform.v1"
)

type PgClient struct {
	db *reform.DB
}

func New() service.Store {
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New(os.Stderr, "SQL: ", log.Flags())
	reformDB := reform.NewDB(db, postgresql.Dialect, reform.NewPrintfLogger(logger.Printf))

	return &PgClient{db: reformDB}
}

func (p *PgClient) CreatePost(news models.News) error {
	tx, err := p.db.Begin()
	if err != nil {
		return errors.Wrap(err, "error starting transaction")
	}

	if p.ifExists(news.Id) {
		return errors.New("id is already exists")
	}

	defer func() {
		if err != nil {
			if e := tx.Rollback(); e != nil {
				err = errors.Wrap(e, "error making rollback")
			}
		}
	}()

	err = tx.Save(&news)
	if err != nil {
		return errors.Wrap(err, "error creating news")
	}

	for _, c := range news.Categories {
		category := &models.NewsCategories{
			CategoryId: c,
			NewsId:     news.Id,
		}
		err = tx.Save(category)
		if err != nil {
			return errors.Wrap(err, "error creating categories")
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "error committing transaction")
	}

	return nil
}

func (p *PgClient) GetPosts(params models.SearchParams) ([]models.News, error) {
	res, err := p.db.SelectAllFrom(models.NewsTable, `LIMIT $1 OFFSET $2`, params.Limit, params.Offset)
	if err != nil {
		return nil, errors.Wrap(err, "error selecting from news table")
	}

	if len(res) == 0 {
		return nil, errors.New("news table is empty")
	}

	var newsList []models.News

	for _, value := range res {
		news := *value.(*models.News)
		categories, err := p.db.SelectAllFrom(models.NewsCategoriesTable, "WHERE news_id = $1", news.Id)

		if err != nil {
			return nil, errors.Wrap(err, "error selecting from news_categories table")
		}

		for _, c := range categories {
			news.Categories = append(news.Categories, c.(*models.NewsCategories).CategoryId)
		}

		newsList = append(newsList, news)
	}
	return newsList, nil
}

func (p *PgClient) EditPost(id int, news models.News) error {
	tx, err := p.db.Begin()
	if err != nil {
		return errors.Wrap(err, "error starting transaction")
	}

	defer func() {
		if err != nil {
			if e := tx.Rollback(); e != nil {
				err = errors.Wrap(e, "error making rollback")
			}
		}
	}()

	if !p.ifExists(id) {
		return errors.New("news with such id doesn't exist")
	}

	updateFields := make(map[string]interface{})
	if news.Title != "" {
		updateFields["title"] = news.Title
	}
	if news.Content != "" {
		updateFields["content"] = news.Content
	}

	if len(updateFields) == 0 && news.Categories == nil {
		return errors.New("nothing to update")
	} else if len(updateFields) > 0 {
		query := "UPDATE news SET"
		values := make([]interface{}, 0)
		idx := 1
		for key, value := range updateFields {
			query += " " + key + " = $" + strconv.Itoa(idx) + ","
			idx++
			values = append(values, value)
		}
		query = strings.TrimSuffix(query, ",")
		query += " WHERE id = $" + strconv.Itoa(idx)
		values = append(values, id)

		_, err = tx.Exec(query, values...)
		if err != nil {
			return errors.Wrap(err, "error updating news")
		}
	}

	if news.Categories != nil {
		_, err = tx.DeleteFrom(models.NewsCategoriesTable, "WHERE news_id = $1", id)
		if err != nil {
			return errors.Wrap(err, "error deleting old categories")
		}

		_, err = tx.Exec("SELECT setval('news_categories_id_seq', (SELECT MAX(id) FROM news_categories));")
		if err != nil {
			return errors.Wrap(err, "error updating id news_categories")
		}

		for _, c := range news.Categories {
			category := &models.NewsCategories{
				CategoryId: c,
				NewsId:     id,
			}
			err = tx.Save(category)
			if err != nil {
				return errors.Wrap(err, "error updating categories")
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "error committing transaction")
	}

	return nil
}

func (p *PgClient) ifExists(id int) bool {
	var news models.News
	if err := p.db.FindByPrimaryKeyTo(&news, id); err != nil {
		return false
	}
	return true
}
