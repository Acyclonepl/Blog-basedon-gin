package global

import (
	"github.com/Acyclonepl/Blog-basedon-gin/pkg/logger"
	"github.com/Acyclonepl/Blog-basedon-gin/pkg/setting"
)

var (
	ServerSetting   *setting.ServerSetting
	AppSetting      *setting.AppSetting
	EmailSetting    *setting.EmailSetting
	JWTSetting      *setting.JWTSetting
	DatabaseSetting *setting.DatabaseSetting
	Logger          *logger.Logger
)
