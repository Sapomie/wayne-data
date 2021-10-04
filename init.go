package main

import (
	"flag"
	"github.com/Sapomie/wayne_data/global"
	"github.com/Sapomie/wayne_data/internal/model"
	"github.com/Sapomie/wayne_data/pkg/setting"
	"time"
)

var (
	port    string
	runMode string
	config  string
)

func BeforeStarting() error {

	SetupFlag()

	//加载配置文件
	err := SetupSetting()
	if err != nil {
		return err
	}
	//启动数据库
	err = SetupDBEngine()
	if err != nil {
		return err
	}

	return nil
}

func SetupFlag() {
	flag.StringVar(&config, "config", "configs/", "配置文件路径")
	flag.StringVar(&runMode, "mode", "", "启动模式")
	flag.StringVar(&port, "port", "", "启动端口")
}

func SetupSetting() error {
	s, err := setting.NewSetting(config)
	if err != nil {
		return err
	}
	err = s.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}

	global.AppSetting.DefaultContextTimeout *= time.Second
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	if port != "" {
		global.ServerSetting.HttpPort = port
	}
	if runMode != "" {
		global.ServerSetting.RunMode = runMode
	}

	return nil
}

func SetupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}
