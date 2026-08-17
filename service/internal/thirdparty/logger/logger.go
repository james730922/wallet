package logger

import (
	"encoding/json"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/james730922/wallet/service/internal/utils/conf"
)

var (
	apLog     ILogger
	sysLog    ILogger
	accessLog ILogger
)

func TestMock() {
	conf.Mock()
	apLog = newLogger(getZapLevel("debug"))
	sysLog = newLogger(getZapLevel("debug"))
	accessLog = newLogger(getZapLevel("info"))
}

func Start() {
	apLog = newLogger(getZapLevel(conf.Log().GetApLog()))
	sysLog = newLogger(getZapLevel(conf.Log().GetSysLog()))
	accessLog = newLogger(getZapLevel(conf.Log().GetAccessLog()))
}

func ApLog() ILogger {
	return apLog.Named("ApLog")
}

func SysLog() ILogger {
	return sysLog.Named("SysLog")
}

func AccessLog() ILogger {
	return accessLog.Named("AccessLog")
}

func newLogger(level zapcore.Level) ILogger {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, level)
	logger := zap.New(core, zap.AddCaller())
	return logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	if conf.Env() == conf.EnvTypeLocal {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}
	// package/file:line
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	return zapcore.AddSync(os.Stdout)
}

func getZapLevel(l string) zapcore.Level {
	switch l {
	case zapcore.DebugLevel.String(): // "debug"
		return zapcore.DebugLevel
	case zapcore.InfoLevel.String(): // "info"
		return zapcore.InfoLevel
	case zapcore.WarnLevel.String(): // "warn"
		return zapcore.WarnLevel
	case zapcore.ErrorLevel.String(): // "error"
		return zapcore.ErrorLevel
	case zapcore.DPanicLevel.String(): // "dpanic"
		return zapcore.DPanicLevel
	case zapcore.PanicLevel.String(): // "panic"
		return zapcore.PanicLevel
	case zapcore.FatalLevel.String(): // "fatal"
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}

type ILogger interface {
	Named(name string) *zap.SugaredLogger

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugw(msg string, fields ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infow(msg string, fields ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Warnw(msg string, fields ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorw(msg string, fields ...interface{})

	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Panicw(msg string, fields ...interface{})
}

type Logger struct {
	logger *zap.SugaredLogger
}

func (lg *Logger) Named(name string) *zap.SugaredLogger {
	lg.logger.Named(name)
	return lg.logger
}

func (lg *Logger) Debug(args ...interface{}) {
	lg.logger.Debug(args...)
}

func (lg *Logger) Debugf(format string, args ...interface{}) {
	lg.logger.Debugf(format, args...)
}

func (lg *Logger) Debugw(msg string, fields ...interface{}) {
	lg.logger.Debugw(msg, fields...)
}

func (lg *Logger) Info(args ...interface{}) {
	lg.logger.Info(args...)
}

func (lg *Logger) Infof(format string, args ...interface{}) {
	lg.logger.Infof(format, args...)
}

func (lg *Logger) Infow(msg string, fields ...interface{}) {
	lg.logger.Infow(msg, fields...)
}

func (lg *Logger) Warn(args ...interface{}) {
	lg.logger.Warn(args...)
}

func (lg *Logger) Warnf(format string, args ...interface{}) {
	lg.logger.Warnf(format, args...)
}

func (lg *Logger) Warnw(msg string, fields ...interface{}) {
	lg.logger.Warnw(msg, fields...)
}

func (lg *Logger) Error(args ...interface{}) {
	lg.logger.Error(args...)
}

func (lg *Logger) Errorf(format string, args ...interface{}) {
	lg.logger.Errorf(format, args...)
}

func (lg *Logger) Errorw(msg string, fields ...interface{}) {
	lg.logger.Errorw(msg, fields...)
}

func (lg *Logger) Panic(args ...interface{}) {
	lg.logger.Panic(args...)
}

func (lg *Logger) Panicf(format string, args ...interface{}) {
	lg.logger.Panicf(format, args...)
}

func (lg *Logger) Panicw(msg string, fields ...interface{}) {
	lg.logger.Panicw(msg, fields...)
}

func Json(obj interface{}) string {
	buf, err := json.Marshal(obj)
	if err != nil {
		return err.Error()
	}

	return string(buf)
}
