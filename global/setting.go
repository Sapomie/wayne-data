package global

import (
	"github.com/Sapomie/wayne_data/pkg/logger"
	"github.com/Sapomie/wayne_data/pkg/setting"
)

var (
	DatabaseSetting *setting.DatabaseSettingS
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	Logger          *logger.Logger
)
