package model

type Article struct {
	*Model
	Title          string `json:"title"`
	Desc           string `json:"desc"`
	Content        string `json:"content"`
	ConverImageUrl string `json:"cover_image_url"`
	State          uint8  `json:"state"`
}

func (a Article) TableName() string {
	return "blog_article"
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
