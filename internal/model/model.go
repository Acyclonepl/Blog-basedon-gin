package model

import (
	"fmt"
	"reflect"
	"time"

	"github.com/Acyclonepl/Blog-basedon-gin/global"
	"github.com/Acyclonepl/Blog-basedon-gin/pkg/setting"
	otgorm "github.com/eddycjy/opentracing-gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	STATE_OPEN  = 1
	STATE_CLOSE = 0
)

type Model struct {
	ID         uint32 `gorm:"primaryKey" json:"id"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	CreatedOn  uint32 `json:"created_on"`
	ModifiedOn uint32 `json:"modified_on"`
	DeletedOn  uint32 `json:"deleted_on"`
	IsDel      uint8  `json:"is_del"`
}

// NewDBEngine 初始化数据库引擎（PostgreSQL + GORM v2）
func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		databaseSetting.Host,
		databaseSetting.Username,
		databaseSetting.Password,
		databaseSetting.DBName,
		databaseSetting.Port,
	)

	// 日志级别
	logLevel := logger.Silent
	if global.ServerSetting != nil && global.ServerSetting.RunMode == "debug" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   databaseSetting.TablePrefix,
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}

	// 注册回调（均采用 Before 钩子，避免与内置逻辑冲突）
	db.Callback().Create().Before("gorm:create").Register("my:update_timestamps_on_create", updateTimeStampForCreateCallback)
	db.Callback().Update().Before("gorm:update").Register("my:update_timestamp_on_update", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Before("gorm:delete").Register("my:soft_delete", deleteCallback)

	// 连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(databaseSetting.MaxIdleConns)
	sqlDB.SetMaxOpenConns(databaseSetting.MaxOpenConns)
	otgorm.AddGormCallbacks(db)
	return db, nil
}

// updateTimeStampForCreateCallback 创建前自动填充 CreatedOn 和 ModifiedOn（若为零值）
func updateTimeStampForCreateCallback(db *gorm.DB) {
	if db.Error != nil {
		return
	}

	stmt := db.Statement
	if stmt.Schema == nil {
		return
	}

	now := time.Now().Unix()

	// 获取模型实例（兼容指针、结构体和切片批量创建）
	rv := stmt.ReflectValue
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface:
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		// 批量创建（切片）直接返回，由内置逻辑处理
		return
	}

	// 处理 CreatedOn
	if field := stmt.Schema.LookUpField("CreatedOn"); field != nil {
		fieldVal, zero := field.ValueOf(rv)
		if zero {
			if err := field.Set(rv, now); err != nil {
				db.AddError(err)
			}
		}
		_ = fieldVal // 忽略未使用的值
	}

	// 处理 ModifiedOn
	if field := stmt.Schema.LookUpField("ModifiedOn"); field != nil {
		_, zero := field.ValueOf(rv)
		if zero {
			if err := field.Set(rv, now); err != nil {
				db.AddError(err)
			}
		}
	}
}

// updateTimeStampForUpdateCallback 更新前强制刷新 ModifiedOn
func updateTimeStampForUpdateCallback(db *gorm.DB) {
	if db.Error != nil {
		return
	}

	stmt := db.Statement
	if stmt.Schema == nil {
		return
	}

	// 若显式使用了 UpdateColumn，不进行自动时间更新（与原逻辑保持一致）
	if _, ok := stmt.Get("gorm:update_column"); ok {
		return
	}

	// 查找并设置 ModifiedOn 字段
	if field := stmt.Schema.LookUpField("ModifiedOn"); field != nil {
		rv := stmt.ReflectValue
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		if rv.Kind() == reflect.Struct {
			if err := field.Set(rv, time.Now().Unix()); err != nil {
				db.AddError(err)
			}
		}
	}
}

// deleteCallback 自定义软删除：将 DELETE 转换为 UPDATE，设置 DeletedOn 和 IsDel
func deleteCallback(db *gorm.DB) {
	if db.Error != nil {
		return
	}

	stmt := db.Statement
	if stmt.Schema == nil {
		return
	}

	// 若使用了 Unscoped() 则执行标准硬删除
	if stmt.Unscoped {
		return
	}

	// 检查模型是否包含软删除字段
	deletedOnField := stmt.Schema.LookUpField("DeletedOn")
	isDelField := stmt.Schema.LookUpField("IsDel")
	if deletedOnField == nil || isDelField == nil {
		return // 模型无需软删除，走默认删除
	}

	now := time.Now().Unix()

	// 将后续生成的 SQL 从 DELETE 转换为 UPDATE
	db.Statement.AddClause(clause.Update{})
	db.Statement.AddClause(clause.Set{
		clause.Assignment{
			Column: clause.Column{Name: deletedOnField.DBName},
			Value:  now,
		},
		clause.Assignment{
			Column: clause.Column{Name: isDelField.DBName},
			Value:  1,
		},
	})
	db.Statement.BuildClauses = []string{"UPDATE", "SET", "WHERE"}

	// 同时更新内存中的对象，保持一致性
	if stmt.ReflectValue.IsValid() {
		rv := stmt.ReflectValue
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		if rv.Kind() == reflect.Struct {
			_ = deletedOnField.Set(rv, now)
			_ = isDelField.Set(rv, 1)
		}
	}
}

// addExtraSpaceIfExist 辅助函数（原代码保留，实际重构后未使用）
func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
