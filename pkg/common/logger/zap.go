package logger

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/novel/pkg/common/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	levelMap = map[string]zapcore.Level{
		"debug":  zapcore.DebugLevel,
		"info":   zapcore.InfoLevel,
		"warn":   zapcore.WarnLevel,
		"error":  zapcore.ErrorLevel,
		"dpanic": zapcore.DPanicLevel,
		"panic":  zapcore.PanicLevel,
		"fatal":  zapcore.FatalLevel,
	}
)

var (
	once   sync.Once
	logger *zap.Logger
)

func Get() *zap.Logger {
	once.Do(func() {
		conf := config.Cfg
		logger = NewZapLogger(conf.LogFile, conf.LogLevel, conf.LogMaxSize, conf.LogBackups, conf.LogMaxAge)
	})
	return logger
}

type ctxKey struct{}

func FromCtx(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return l
	} else if l := logger; l != nil {
		return l
	}

	return zap.NewNop()
}

func WithCtx(ctx context.Context, l *zap.Logger) context.Context {
	if lp, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		if lp == l {
			return ctx
		}
	}

	return context.WithValue(ctx, ctxKey{}, l)
}

func NewZapLogger(path string, logLevel string, maxSize int, maxBackups int, maxAge int) *zap.Logger {
	l := levelMap[strings.ToLower(logLevel)]
	level := zap.NewAtomicLevelAt(l)

	encoder := getEncoder()

	var core zapcore.Core
	logWs := zapcore.AddSync(log.Writer())
	if path != "" {
		file := zapcore.AddSync(&lumberjack.Logger{
			Filename:   path,
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			MaxAge:     maxAge,
		})
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, logWs, level),
			zapcore.NewCore(encoder, file, level),
		)
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, logWs, level),
		)
	}
	return zap.New(core, zap.AddCaller())
}

func getEncoder() zapcore.Encoder {
	conf := zap.NewProductionEncoderConfig()
	conf.TimeKey = "timestamp"
	conf.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewJSONEncoder(conf)
}
