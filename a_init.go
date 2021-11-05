package main

import (
	"flag"
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/pkg/logger"
	"github.com/Sapomie/wayne-data/pkg/setting"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"time"
)

var (
	port    string
	runMode string
	config  string
)

func BeforeStarting() error {
	setupFlag()
	//加载配置文件
	err := setupSetting()
	if err != nil {
		return err
	}
	//启动数据库
	err = setupDBEngine()
	if err != nil {
		return err
	}
	//日志
	err = setupLogger()
	if err != nil {
		return err
	}
	//数据库 data
	return nil
}

func setupFlag() {
	flag.StringVar(&config, "config", "configs/", "配置文件路径")
	flag.StringVar(&runMode, "mode", "", "启动模式")
	flag.StringVar(&port, "port", "", "启动端口")
}

func setupSetting() error {
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

func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}

func setupLogger() error {
	fileName := global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt
	global.Logger = logger.NewLogger(
		&lumberjack.Logger{
			Filename:  fileName,
			MaxSize:   500,
			MaxAge:    10,
			LocalTime: true,
		},
		"",
		log.LstdFlags,
	).WithCaller(1)
	return nil
}
