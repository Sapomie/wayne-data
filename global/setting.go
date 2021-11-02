package global

import (
	"github.com/Sapomie/wayne-data/pkg/logger"
	"github.com/Sapomie/wayne-data/pkg/setting"
)

var (
	DatabaseSetting *setting.DatabaseSettingS
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	Logger          *logger.Logger
)
