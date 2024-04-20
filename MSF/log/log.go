package fllog

import (
	"fmt"

	"github.com/xiaolongdeng1990/forlife/MSF/config"
)

type LogCfg struct {
	LogConf struct {
		Name       string `default:"../log/fllog.log"`
		Level      string `default:"INFO"`
		MaxSize    int    `default:"1073741824"`
		MaxAge     int    `default:"30"` //默认最多保存10个日志文件
		MaxBackups int    `default:"10"` // 最大保存日志数量
	}
}

func Init(cfg string) error {
	logCfg := LogCfg{}
	if err := config.ParseConfigWithPath(&logCfg, cfg); err != nil {
		fmt.Printf("load logcfg failed. err:%+v cfg:%s", err, cfg)
		return err
	}
	builder := NewLogUtilsBuilder(
		logCfg.LogConf.Level,
		logCfg.LogConf.Name,
		logCfg.LogConf.MaxSize,
		logCfg.LogConf.MaxAge,
		logCfg.LogConf.MaxBackups,
		true,
		true,
	)
	logUtils := NewLogUtils().SetBuilder(builder)
	err := logUtils.Init()
	if err != nil {
		return err
	}
	return nil
}

func Debug(f string, p ...interface{}) {
	if Allow("DEBUG") {
		msg := fmt.Sprintf(f, p...)
		Log().Debug(msg)
	}
}

func Info(f string, p ...interface{}) {
	if Allow("INFO") {
		msg := fmt.Sprintf(f, p...)
		Log().Info(msg)
	}
}

func Warn(f string, p ...interface{}) {
	if Allow("WARN") {
		msg := fmt.Sprintf(f, p...)
		Log().Warn(msg)
	}
}

func Error(f string, p ...interface{}) {
	if Allow("ERROR") {
		msg := fmt.Sprintf(f, p...)
		Log().Error(msg)
	}
}
