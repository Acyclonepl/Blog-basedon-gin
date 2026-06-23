package global

import (
	"github.com/Acyclonepl/Blog-basedon-gin/pkg/logger"
	"github.com/Acyclonepl/Blog-basedon-gin/pkg/setting"
)

var (
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS
	Logger          *logger.Logger
)
