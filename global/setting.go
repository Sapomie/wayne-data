package global

import (
	"github.com/Sapomie/wayne-data/pkg/loggerV0"
	"github.com/Sapomie/wayne-data/pkg/setting"
)

var (
	DatabaseSetting *setting.DatabaseSettingS
	RedisSetting    *setting.RedisSettingS
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	Logger          *loggerV0.Logger
)
