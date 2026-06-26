package model

import (
	"github.com/Acyclonepl/Blog-basedon-gin/pkg/app"
	"gorm.io/gorm" // 仅此处修改
)

type Article struct {
	*Model
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	State         uint8  `json:"state"`
}

type ArticleSwagger struct {
	List  []*Article
	Pager *app.Pager
}

func (a Article) TableName() string {
	return "blog_article"
}

func (a Article) Create(db *gorm.DB) (*Article, error) {
	if err := db.Create(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (a Article) Update(db *gorm.DB, values interface{}) error {
	if err := db.Model(&a).Where("id = ? AND is_del = ?", a.ID, 0).Updates(values).Error; err != nil {
		return err
	}
	return nil
}

func (a Article) Get(db *gorm.DB) (Article, error) {
	var article Article
	db = db.Where("id = ? AND state = ? AND is_del = ?", a.ID, a.State, 0)
	err := db.First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return article, err
	}
	return article, nil
}

func (a Article) Delete(db *gorm.DB) error {
	if err := db.Where("id = ? AND is_del = ?", a.Model.ID, 0).Delete(&a).Error; err != nil {
		return err
	}
	return nil
}

type ArticleRow struct {
	ArticleID     uint32
	TagID         uint32
	TagName       string
	ArticleTitle  string
	ArticleDesc   string
	CoverImageUrl string
	Content       string
}

func (a Article) ListByTagID(db *gorm.DB, tagID uint32, pageOffset, pageSize int) ([]*ArticleRow, error) {
	fields := []string{"ar.id AS article_id", "ar.title AS article_title", "ar.desc AS article_desc", "ar.cover_image_url", "ar.content"}
	fields = append(fields, []string{"t.id AS tag_id", "t.name AS tag_name"}...)

	if pageOffset >= 0 && pageSize > 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}
	rows, err := db.Select(fields).Table(ArticleTag{}.TableName()+" AS at").
		Joins("LEFT JOIN `"+Tag{}.TableName()+"` AS t ON at.tag_id = t.id").
		Joins("LEFT JOIN `"+Article{}.TableName()+"` AS ar ON at.article_id = ar.id").
		Where("at.`tag_id` = ? AND ar.state = ? AND ar.is_del = ?", tagID, a.State, 0).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*ArticleRow
	for rows.Next() {
		r := &ArticleRow{}
		if err := rows.Scan(&r.ArticleID, &r.ArticleTitle, &r.ArticleDesc, &r.CoverImageUrl, &r.Content, &r.TagID, &r.TagName); err != nil {
			return nil, err
		}
		articles = append(articles, r)
	}
	return articles, nil
}

func (a Article) CountByTagID(db *gorm.DB, tagID uint32) (int, error) {
	var count int
	err := db.Table(ArticleTag{}.TableName()+" AS at").
		Joins("LEFT JOIN `"+Tag{}.TableName()+"` AS t ON at.tag_id = t.id").
		Joins("LEFT JOIN `"+Article{}.TableName()+"` AS ar ON at.article_id = ar.id").
		Where("at.`tag_id` = ? AND ar.state = ? AND ar.is_del = ?", tagID, a.State, 0).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CreateArticleRequest 创建文章请求参数
type CreateArticleRequest struct {
	Title         string `json:"title" binding:"required,min=3,max=100" example:"深入理解Go并发"`
	Desc          string `json:"desc" binding:"max=255" example:"本文详细讲解Go语言的并发模型"`
	Content       string `json:"content" binding:"required" example:"Go的并发基于goroutine和channel..."`
	CoverImageUrl string `json:"cover_image_url" binding:"omitempty,url" example:"https://example.com/cover.jpg"`
	State         int    `json:"state" binding:"oneof=0 1" example:"1"` // 代码内可设默认值
	CreatedBy     string `json:"created_by" binding:"required,min=3,max=50" example:"admin"`
}

// UpdateArticleRequest 更新文章请求参数
type UpdateArticleRequest struct {
	Title         *string `json:"title" binding:"omitempty,min=3,max=100" example:"更新后的标题"`
	Desc          *string `json:"desc" binding:"omitempty,max=255" example:"更新后的简介"`
	Content       *string `json:"content" binding:"omitempty" example:"更新后的正文内容"`
	CoverImageUrl *string `json:"cover_image_url" binding:"omitempty,url" example:"https://example.com/new_cover.jpg"`
	State         *int    `json:"state" binding:"omitempty,oneof=0 1" example:"0"`
	ModifiedBy    string  `json:"modified_by" binding:"required,min=3,max=50" example:"editor"`
}
