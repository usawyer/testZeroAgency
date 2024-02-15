package database

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/usawyer/testZeroAgency/internal/service"
	"github.com/usawyer/testZeroAgency/models"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"time"
)

type PgClient struct {
	db *gorm.DB
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

//func New() {
//	db, err := reform.NewDB()
//}

func New() service.Store {
	connectionParams := map[string]string{
		"host":     getEnv("DB_HOST", "localhost"),
		"user":     getEnv("POSTGRES_USER", "postgres"),
		"password": getEnv("POSTGRES_PASSWORD", "postgres"),
		"dbname":   getEnv("POSTGRES_DB", "test"),
		"port":     getEnv("DB_PORT", "5432"),
		"sslmode":  "disable",
		"TimeZone": "Asia/Novosibirsk",
	}
	//gormLogger := zapgorm2.New(zapLogger)
	var dsn string

	for key, value := range connectionParams {
		dsn += fmt.Sprintf("%s=%s ", key, value)
	}
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * 2)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
		if err != nil {
			log.Warn("Error open", zap.Error(err))
			continue
		}
		db = db.Debug()
		err = db.AutoMigrate(&models.News{})
		if err != nil {
			log.Error(err.Error())
		}
		return &PgClient{db: db}
	}
	log.Fatal("Error open db")
	return nil
}

//func (p *PgClient) NewNews(news models.News) error {
//	p.db.Create(&news)
//	return nil
//}

func (p *PgClient) CreatePost(news models.News) {
	p.db.Create(&news)
}

func (p *PgClient) GetPosts(params models.SearchParams) ([]models.News, error) {
	var news []models.News
	res := p.db.Offset(params.Offset).Limit(params.Limit).Find(&news)
	return news, res.Error
}

func (p *PgClient) EditPost(id int, news models.News) {
	err := p.db.Create(&news)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *PgClient) IfExists(id int) bool {
	return false
}
