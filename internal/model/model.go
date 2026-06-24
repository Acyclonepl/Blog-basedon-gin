package model

import (
	"fmt"

	"github.com/Acyclonepl/Blog-basedon-gin/global"
	"github.com/Acyclonepl/Blog-basedon-gin/pkg/app"
	"github.com/Acyclonepl/Blog-basedon-gin/pkg/setting"
	"gorm.io/driver/postgres" // 导入 PostgreSQL 驱动
	"gorm.io/gorm"
	"gorm.io/gorm/logger" // GORM 的日志包
	"gorm.io/gorm/schema"
)

type Model struct {
	ID         uint32 `gorm:"primary_key" json:"id"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	CreatedOn  uint32 `json:"created_on"`
	ModifiedOn uint32 `json:"modified_on"`
	DeletedOn  uint32 `json:"deleted_on"`
	IsDel      uint8  `json:"is_del"`
}

func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		databaseSetting.Host,
		databaseSetting.Username,
		databaseSetting.Password,
		databaseSetting.DBName,
		databaseSetting.Port,
	)
	logLevel := logger.Silent
	if global.ServerSetting != nil && global.ServerSetting.RunMode == "debug" {
		logLevel = logger.Info
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logLevel), NamingStrategy: schema.NamingStrategy{
		TablePrefix:   databaseSetting.TablePrefix,
		SingularTable: true,
	}})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(databaseSetting.MaxIdleConns)
	sqlDB.SetMaxOpenConns(databaseSetting.MaxOpenConns)
	return db, nil
}

// tag.go
type TagSwagger struct {
	List  []*Tag
	Pager *app.Pager
}

// article.go
type ArticleSwagger struct {
	List  []*Article
	Pager *app.Pager
}
