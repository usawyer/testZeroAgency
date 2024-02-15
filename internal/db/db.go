package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/usawyer/testZeroAgency/internal/service"
	"github.com/usawyer/testZeroAgency/models"
	"gopkg.in/reform.v1/dialects/postgresql"
	"os"
	"time"

	"log"

	"gopkg.in/reform.v1"
)

type PgClient struct {
	db *reform.DB
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// сложить в отдельный файл?

func initDB() (*sql.DB, error) {
	connectionParams := map[string]string{
		"host":     getEnv("DB_HOST", "localhost"),
		"user":     getEnv("POSTGRES_USER", "postgres"),
		"password": getEnv("POSTGRES_PASSWORD", "postgres"),
		"dbname":   getEnv("POSTGRES_DB", "test"),
		"port":     getEnv("DB_PORT", "5432"),
		"sslmode":  "disable",
		"TimeZone": "Asia/Novosibirsk",
	}

	var dsn string
	for key, value := range connectionParams {
		dsn += fmt.Sprintf("%s=%s ", key, value)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Minute * 5)

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
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
	defer func() {
		if err != nil {
			tx.Rollback()
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

	var newslist []models.News

	for _, value := range res {
		news := *value.(*models.News)
		categories, err := p.db.SelectAllFrom(models.NewsCategoriesTable, "WHERE news_id = $1", news.Id)

		if err != nil {
			return nil, errors.Wrap(err, "error selecting from news_categories table")
		}

		for _, c := range categories {
			news.Categories = append(news.Categories, c.(*models.NewsCategories).CategoryId)
		}

		newslist = append(newslist, news)
	}
	return newslist, nil
}

func (p *PgClient) EditPost(id int, news models.News) {
	//err := p.db.Create(&news)
	//if err != nil {
	//	fmt.Println(err)
	//}
	log.Println("EDITED POST")
}

func (p *PgClient) IfExists(id int) bool {
	var news models.News
	if err := p.db.FindByPrimaryKeyTo(&news, id); err != nil {
		return false
	}
	return true
}

//	// Пример обновления пользователя
//	userByID.Age = 31
//	if err := db2.Save(userByID); err != nil {
//		log.Fatal(err)
//	}
