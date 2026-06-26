package service

import (
	"context"

	"github.com/Acyclonepl/Blog-basedon-gin/global"
	"github.com/Acyclonepl/Blog-basedon-gin/internal/dao"
)

type Service struct {
	ctx context.Context
	dao *dao.Dao
}

func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	svc.dao = dao.New(global.DBEngine.WithContext(ctx)) // 直接使用 v2 的 WithContext
	return svc
}
