package log

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/topfreegames/pitaya/v2/logger/interfaces"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

const ctxPitayaLoggerKey = "zapPitayaLogger"

type PitayaLogger struct {
	*zap.Logger
}

func NewPitayaLog(conf *viper.Viper) *PitayaLogger {
	// log address "out.log" User-defined
	lp := conf.GetString("log.log_file_name")
	lv := conf.GetString("log.log_level")
	var level zapcore.Level
	//debug<info<warn<error<fatal<panic
	switch lv {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	hook := lumberjack.Logger{
		Filename:   lp,                             // Log file path
		MaxSize:    conf.GetInt("log.max_size"),    // Maximum size unit for each log file: M
		MaxBackups: conf.GetInt("log.max_backups"), // The maximum number of backups that can be saved for log files
		MaxAge:     conf.GetInt("log.max_age"),     // Maximum number of days the file can be saved
		Compress:   conf.GetBool("log.compress"),   // Compression or not
	}

	var encoder zapcore.Encoder
	if conf.GetString("log.encoding") == "console" {
		encoder = zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "Logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
			EncodeTime:     timeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		})
	} else {
		encoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		})
	}
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // Print to console and file
		level,
	)
	if conf.GetString("env") != "prod" {
		return &PitayaLogger{zap.New(core, zap.Development(), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel), zap.AddCallerSkip(1))}
	}
	return &PitayaLogger{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel), zap.AddCallerSkip(1))}
}

// WithValue Adds a field to the specified context
func (l *PitayaLogger) WithValue(ctx context.Context, fields ...zapcore.Field) context.Context {
	if c, ok := ctx.(*gin.Context); ok {
		ctx = c.Request.Context()
		c.Request = c.Request.WithContext(context.WithValue(ctx, ctxPitayaLoggerKey, l.WithContext(ctx).With(fields...)))
		return c
	}
	return context.WithValue(ctx, ctxPitayaLoggerKey, l.WithContext(ctx).With(fields...))
}

// WithContext Returns a zap instance from the specified context
func (l *PitayaLogger) WithContext(ctx context.Context) *PitayaLogger {
	if c, ok := ctx.(*gin.Context); ok {
		ctx = c.Request.Context()
	}
	zl := ctx.Value(ctxPitayaLoggerKey)
	ctxLogger, ok := zl.(*zap.Logger)
	if ok {
		return &PitayaLogger{Logger: ctxLogger}
	}
	return l
}

func (l *PitayaLogger) Fatal(format ...interface{}) {
	l.Sugar().Fatal(format...)
}

func (l *PitayaLogger) Fatalf(format string, args ...interface{}) {
	l.Sugar().Fatalf(format, args...)
}

func (l *PitayaLogger) Fatalln(args ...interface{}) {
	l.Sugar().Fatalln(args...)
}

func (l *PitayaLogger) Debug(args ...interface{}) {
	l.Sugar().Debug(args...)
}

func (l *PitayaLogger) Debugf(format string, args ...interface{}) {
	l.Sugar().Debugf(format, args...)
}

func (l *PitayaLogger) Debugln(args ...interface{}) {
	l.Sugar().Debugln(args...)
}

func (l *PitayaLogger) Error(args ...interface{}) {
	l.Sugar().Error(args...)
}

func (l *PitayaLogger) Errorf(format string, args ...interface{}) {
	l.Sugar().Errorf(format, args...)
}

func (l *PitayaLogger) Errorln(args ...interface{}) {
	l.Sugar().Errorln(args...)
}

func (l *PitayaLogger) Info(args ...interface{}) {
	l.Sugar().Info(args...)
}

func (l *PitayaLogger) Infof(format string, args ...interface{}) {
	l.Sugar().Infof(format, args...)
}

func (l *PitayaLogger) Infoln(args ...interface{}) {
	l.Sugar().Infoln(args...)
}

func (l *PitayaLogger) Warn(args ...interface{}) {
	l.Sugar().Warn(args...)
}

func (l *PitayaLogger) Warnf(format string, args ...interface{}) {
	l.Sugar().Warnf(format, args...)
}

func (l *PitayaLogger) Warnln(args ...interface{}) {
	l.Sugar().Warnln(args...)
}

func (l *PitayaLogger) Panic(args ...interface{}) {
	l.Sugar().Panic(args...)
}

func (l *PitayaLogger) Panicf(format string, args ...interface{}) {
	l.Sugar().Panicf(format, args...)
}

func (l *PitayaLogger) Panicln(args ...interface{}) {
	l.Sugar().Panicln(args...)
}

func (l *PitayaLogger) WithFields(fields map[string]interface{}) interfaces.Logger {
	return &PitayaLogger{l.Logger.With(zap.Any("fields", fields))}
}

func (l *PitayaLogger) WithField(key string, value interface{}) interfaces.Logger {
	return &PitayaLogger{l.With(zap.Any(key, value))}
}

func (l *PitayaLogger) WithError(err error) interfaces.Logger {
	return &PitayaLogger{l.With(zap.Error(err))}
}

func (l *PitayaLogger) GetInternalLogger() any {
	return l.Sugar()
}
