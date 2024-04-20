package fllog

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var LogLevelMap = map[string]int{
	"debug": 1,
	"info":  2,
	"warn":  3,
	"error": 4,
}

type MyLogUtilsBuilder struct {
	logLevel    string
	logFileName string
	maxSize     int
	maxAge      int
	maxBackups  int
	status      bool
	line        bool
}

func (m *MyLogUtilsBuilder) SetLine(line bool) {
	m.line = line
}

func (m *MyLogUtilsBuilder) GetLine() bool {
	return m.line
}

func (m *MyLogUtilsBuilder) GetLogFileName() string {
	return m.logFileName
}

func (m *MyLogUtilsBuilder) SetConsole(status bool) BuilderInterface {
	m.status = status
	return m
}

func (m *MyLogUtilsBuilder) GetConsole() bool {
	return m.status
}

func NewLogUtilsBuilder(logLevel string, logFileName string, maxSize int, maxAge int, maxBackups int, status bool, line bool) BuilderInterface {
	myLogUtilsBuilder := &MyLogUtilsBuilder{
		logLevel:    logLevel,
		logFileName: logFileName,
		maxSize:     maxSize,
		maxAge:      maxAge,
		maxBackups:  maxBackups,
		status:      status,
		line:        line,
	}
	return myLogUtilsBuilder
}

func (m *MyLogUtilsBuilder) SetLogLevel(logLevel string) BuilderInterface {
	m.logLevel = logLevel
	return m
}

func (m *MyLogUtilsBuilder) GetLogLevel() string {
	//m.logLevel = logLevel
	return m.logLevel
}

func (m *MyLogUtilsBuilder) SetLogFileName(logFileName string) BuilderInterface {
	m.logFileName = logFileName
	return m
}

func (m *MyLogUtilsBuilder) SetMaxSize(MaxSize int) BuilderInterface {
	m.maxSize = MaxSize
	return m
}

func (m *MyLogUtilsBuilder) SetMaxAge(MaxAge int) BuilderInterface {
	m.maxAge = MaxAge
	return m
}

func (m *MyLogUtilsBuilder) SetMaxBackups(MaxBackups int) BuilderInterface {
	m.maxBackups = MaxBackups
	return m
}

type myLogUtils struct {
	builders      BuilderInterface
	sugaredLogger *zap.SugaredLogger
}

var once sync.Once
var instance *myLogUtils

func NewLogUtils() *myLogUtils {
	once.Do(func() {
		instance = &myLogUtils{}
	})
	return instance
}

func (ms *myLogUtils) SetBuilder(builder BuilderInterface) *myLogUtils {
	ms.builders = builder
	return ms
}

func (ms *myLogUtils) getLogsUtils() *zap.SugaredLogger {
	return ms.sugaredLogger
}

func Log() *zap.SugaredLogger {
	utils := NewLogUtils()
	err := utils.Init()
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	return utils.getLogsUtils()
}

func Allow(level string) bool {
	utils := NewLogUtils()
	logLevel := utils.builders.GetLogLevel()
	realLevel:= LogLevelMap[strings.ToLower(logLevel)]
	checkLevel := LogLevelMap[strings.ToLower(level)]
	return realLevel >= checkLevel
}

func (ms *myLogUtils) Init() error {
	var showLine string
	// 日志级别
	//logLevel := "DEBUG"
	atomicLevel := zap.NewAtomicLevel()
	switch ms.builders.GetLogLevel() {
	case "DEBUG":
		atomicLevel.SetLevel(zapcore.DebugLevel)
	case "INFO":
		atomicLevel.SetLevel(zapcore.InfoLevel)
	case "WARN":
		atomicLevel.SetLevel(zapcore.WarnLevel)
	case "ERROR":
		atomicLevel.SetLevel(zapcore.ErrorLevel)
	case "DPANIC":
		atomicLevel.SetLevel(zapcore.DPanicLevel)
	case "PANIC":
		atomicLevel.SetLevel(zapcore.PanicLevel)
	case "FATAL":
		atomicLevel.SetLevel(zapcore.FatalLevel)
	}

	if ms.builders.GetConsole() {
		if ms.builders.GetLine() {
			showLine = "line"
		} else {
			showLine = ""
		}

		// 自定义日志级别颜色
		colors := map[zapcore.Level]color.Attribute{
			zapcore.DebugLevel:  color.FgBlue,
			zapcore.InfoLevel:   color.FgGreen,
			zapcore.WarnLevel:   color.FgYellow,
			zapcore.ErrorLevel:  color.FgRed,
			zapcore.DPanicLevel: color.FgMagenta,
			zapcore.PanicLevel:  color.FgMagenta,
			zapcore.FatalLevel:  color.FgMagenta,
		}

		encoderConfig := zapcore.EncoderConfig{
			TimeKey:          "time",
			LevelKey:         "level",
			MessageKey:       "msg",
			CallerKey:        showLine,
			LineEnding:       zapcore.DefaultLineEnding,
			EncodeLevel:      coloredLevelEncoder(colors),
			EncodeTime:       zapcore.TimeEncoderOfLayout("[2006-01-02 15:04:05]"),
			EncodeDuration:   zapcore.SecondsDurationEncoder,
			EncodeCaller:     zapcore.FullCallerEncoder, //.ShortCallerEncoder,
			EncodeName:       zapcore.FullNameEncoder,
			ConsoleSeparator: "",
		}

		zapCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			atomicLevel,
		)
		logutils := zap.New(zapCore, zap.AddCaller()).Sugar()
		ms.sugaredLogger = logutils
	} else {
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "name",
			CallerKey:      "line",
			MessageKey:     "msg",
			FunctionKey:    "func",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		}
		// 日志轮转
		writer := &lumberjack.Logger{
			// 日志名称
			Filename: ms.builders.GetLogFileName(),
			// 日志大小限制，单位MB
			MaxSize: 50,
			// 历史日志文件保留天数
			MaxAge: 30,
			// 最大保留历史日志数量,其实就是备份数量
			MaxBackups: 10,
			// 本地时区
			LocalTime: true,
			// 历史日志文件压缩标识
			Compress: false,
		}

		zapCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(writer),
			atomicLevel,
		)

		logutils := zap.New(zapCore, zap.AddCaller()).Sugar()
		ms.sugaredLogger = logutils
	}
	return nil
}

// 自定义带颜色的日志级别编码函数
func coloredLevelEncoder(colors map[zapcore.Level]color.Attribute) zapcore.LevelEncoder {
	return func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		c, ok := colors[l]
		if !ok {
			c = color.Reset // 默认为重置颜色
		}

		enc.AppendString(color.New(c).Sprint(l.String()))
	}
}
