package main

import (
	"flag"
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/service/a_procession"
	"github.com/Sapomie/wayne-data/internal/service/b_raw_event"
	"github.com/Sapomie/wayne-data/pkg/log"
	"github.com/Sapomie/wayne-data/pkg/loggerV0"
	"github.com/Sapomie/wayne-data/pkg/setting"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	log2 "log"
	"os"
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
	//日志
	err = setupLogger()
	if err != nil {
		return err
	}
	//日志V2
	err = setupLoggerV2()
	if err != nil {
		return err
	}
	//启动数据库
	err = setupDBEngine()
	if err != nil {
		return err
	}
	//启动redis
	err = setupRedisEngine()
	if err != nil {
		return err
	}

	//task默认值
	err = b_raw_event.NewRawEventService(nil, global.DBEngine, global.CacheEngine).ReadDefaultTaskValue()
	if err != nil {
		return err
	}
	//更新field 全局变量
	_, err = a_procession.NewProcessionService(nil, global.DBEngine, global.CacheEngine).ProcessAll()
	if err != nil {
		return err
	}

	return nil
}

func setupLoggerV2() error {
	//
	var w io.Writer = os.Stdout
	if err := os.MkdirAll(global.AppSetting.LogV2SavePath, 0755); err != nil {
		log.Fatal("initLogger", "mkdir log path failed")
	}
	if fp, err := os.OpenFile(global.AppSetting.LogV2SavePath+"/"+global.AppSetting.LogV2FileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755); err != nil {
		log.Fatal("initLogger", "log path not exists")
	} else {
		w = io.Writer(fp)
	}
	log.SetDefaultFileLogger(log.NewFileLogger(w, "info", 0))

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
	err = s.ReadSection("Redis", &global.RedisSetting)
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

func setupRedisEngine() error {
	var err error
	global.CacheEngine, err = model.NewCacheEngine(global.RedisSetting)
	if err != nil {
		return err
	}

	return model.NewCache(global.CacheEngine).FlushDb()
}

func setupLogger() error {
	fileName := global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt
	global.Logger = loggerV0.NewLogger(
		&lumberjack.Logger{
			Filename:  fileName,
			MaxSize:   500,
			MaxAge:    10,
			LocalTime: true,
		},
		"",
		log2.LstdFlags,
	).WithCaller(1)
	return nil
}
