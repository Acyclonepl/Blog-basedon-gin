package model

import (
	"github.com/Acyclonepl/Blog-basedon-gin/pkg/app"
	"gorm.io/gorm"
)

type Tag struct {
	*Model
	Name  string `json:"name"`
	State uint8  `json:"state"`
}

func (t Tag) TableName() string {
	return "blog_tag"
}

type TagSwagger struct {
	List  []*Tag
	Pager *app.Pager
}

func (t Tag) Count(db *gorm.DB) (int, error) {
	var count int64
	if t.Name != "" {
		db = db.Where("name=?", t.Name)
	}
	db = db.Where("state=?", t.State)
	if err := db.Model(&t).Where("is_del=?", 0).Count(&count).Error; err != nil {
		return 0, err
	}
	result := int(count)
	return result, nil
}
func (t Tag) List(db *gorm.DB, pageOffset, pageSize int) ([]*Tag, error) {
	var tags []*Tag
	var err error
	if pageOffset >= 0 && pageSize > 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}
	if t.Name != "" {
		db = db.Where("name=?", t.Name)
	}
	db = db.Where("state=?", t.State)
	if err = db.Where("is_del=?", 0).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}
func (t Tag) ListByIDs(db *gorm.DB, ids []uint32) ([]*Tag, error) {
	var tags []*Tag
	db = db.Where("state = ? AND is_del = ?", t.State, 0)
	err := db.Where("id IN (?)", ids).Find(&tags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return tags, nil
}
func (t Tag) Get(db *gorm.DB) (Tag, error) {
	var tag Tag
	err := db.Where("id = ? AND is_del = ? AND state = ?", t.ID, 0, t.State).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return tag, err
	}

	return tag, nil
}
func (t Tag) Create(db *gorm.DB) error {
	return db.Create(&t).Error
}
func (t Tag) Update(db *gorm.DB, values interface{}) error {
	db = db.Model(&Tag{}).Where("id=? AND is_del=?", t.ID, 0)
	return db.Model(&t).Where("id =? AND is_del=?", t.ID, 0).Updates(values).Error
}
func (t Tag) Delete(db *gorm.DB) error {
	return db.Where("id=? AND is_del=?", t.Model.ID, 0).Delete(&t).Error
}

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name      string `json:"name" binding:"required,min=3,max=100" example:"Go"`
	State     int    `json:"state" binding:"oneof=0 1" example:"1"`
	CreatedBy string `json:"created_by" binding:"required,min=3,max=100" example:"admin"`
}

// UpdateTagRequest 更新标签请求
type UpdateTagRequest struct {
	Name       string `json:"name" binding:"min=3,max=100" example:"Golang"`
	State      int    `json:"state" binding:"oneof=0 1" example:"1"`
	ModifiedBy string `json:"modified_by" binding:"required,min=3,max=100" example:"admin"`
}
