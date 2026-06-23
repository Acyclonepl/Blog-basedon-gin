package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Acyclonepl/Blog-basedon-gin/global"
	"github.com/Acyclonepl/Blog-basedon-gin/internal/model"
	"github.com/Acyclonepl/Blog-basedon-gin/internal/routers"
	"github.com/Acyclonepl/Blog-basedon-gin/pkg/logger"
	"github.com/Acyclonepl/Blog-basedon-gin/pkg/setting"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
)

func init() {

	if err := setupSetting(); err != nil {
		log.Fatalf("init.setupSetting err:%v", err)
	}

	if err := setupLogger(); err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}
	if err := setupDBEngine(); err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}

}

// @title 博客系统
// @version 1.0
// @description 基于Gin的博客系统
// @termsOfService https://github.com/Acyclonepl/Blog-basedon-gin
func main() {
	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("服务启动失败：%v", err)
	}

}
func setupSetting() error {
	setting, err := setting.NewSetting()
	if err != nil {
		return err
	}
	global.ServerSetting = &setting.ServerSettingS{}
	global.AppSetting = &setting.AppSettingS{}
	global.DatabaseSetting = &setting.DatabaseSettingS{}
	err = setting.ReadSection("Server", global.ServerSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("App", global.AppSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Database", global.DatabaseSetting)
	if err != nil {
		return err
	}
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	return nil
}
func setupLogger() error {
	fileName := global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  fileName,
		MaxSize:   600,
		MaxAge:    10,
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)
	return nil
}
func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	return err

}
