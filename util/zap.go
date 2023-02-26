package util

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
)

var (
	//logger zap 默认的 logger 不支持格式化输出，要打印指定值要用 zap.String、zap.Int 等封装，代码就显得非常冗长
	logger *zap.Logger
	//sugar 但它提供了 Sugar（语法糖的糖），只要一点点额外的性能损失（但是仍比大部分库快），可以比较简单地格式化输出。
	sugar *zap.SugaredLogger
	once  sync.Once
)

func init() {
	once.Do(func() {
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   "logs/ota.log",
			MaxSize:    10,    // 单个日志文件最大大小，单位 MB
			MaxBackups: 10,    // 保留的历史日志文件个数
			MaxAge:     30,    // 保留历史日志文件的最大天数
			Compress:   false, // 是否压缩历史日志文件
		})

		// 设置日志级别
		atomicLevel := zap.NewAtomicLevel()
		atomicLevel.SetLevel(zapcore.DebugLevel)
		// 配置日志输出格式
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

		// 创建核心日志记录器
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), fileWriter),
			atomicLevel,
		)
		logger = zap.New(core)
		defer logger.Sync()
		sugar = logger.Sugar()
	})
}

func GetLogger() *zap.Logger {
	return logger
}

func GetSugar() *zap.SugaredLogger {
	return sugar
}
