package model

type Tag struct {
	*Model
	Name  string `json:"name"`
	State uint8  `json:"state"`
}

func (t Tag) TableName() string {
	return "blog_tag"
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
