package log

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Sync() error
	Named(name string) Logger
	With(args ...interface{}) Logger
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	DPanic(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	DPanicf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	DPanicw(msg string, keysAndValues ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
}

type logger struct {
	zapLogger *zap.SugaredLogger
}

func (l *logger) Sync() error {
	return l.zapLogger.Sync()
}

func (l *logger) Named(name string) Logger {
	return &logger{l.zapLogger.Named(name)}
}

func (l *logger) With(args ...interface{}) Logger {
	return &logger{l.zapLogger.With(args...)}
}

func (l *logger) Debug(args ...interface{}) {
	l.zapLogger.Debug(args...)
}

func (l *logger) Info(args ...interface{}) {
	l.zapLogger.Info(args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.zapLogger.Info(args...)
}

func (l *logger) Error(args ...interface{}) {
	l.zapLogger.Error(args...)
}

func (l *logger) DPanic(args ...interface{}) {
	l.zapLogger.DPanic(args...)
}

func (l *logger) Panic(args ...interface{}) {
	l.zapLogger.Panic(args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.zapLogger.Fatal(args...)
}

func (l *logger) Debugf(template string, args ...interface{}) {
	l.zapLogger.Debugf(template, args...)
}

func (l *logger) Infof(template string, args ...interface{}) {
	l.zapLogger.Infof(template, args...)
}

func (l *logger) Warnf(template string, args ...interface{}) {
	l.zapLogger.Warnf(template, args...)
}

func (l *logger) Errorf(template string, args ...interface{}) {
	l.zapLogger.Errorf(template, args...)
}

func (l *logger) DPanicf(template string, args ...interface{}) {
	l.zapLogger.DPanicf(template, args...)
}

func (l *logger) Panicf(template string, args ...interface{}) {
	l.zapLogger.Panicf(template, args...)
}

func (l *logger) Fatalf(template string, args ...interface{}) {
	l.zapLogger.Fatalf(template, args...)
}

func (l *logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Debugw(msg, keysAndValues...)
}

func (l *logger) Infow(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Infow(msg, keysAndValues...)
}

func (l *logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Warnw(msg, keysAndValues...)
}

func (l *logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Errorw(msg, keysAndValues...)
}

func (l *logger) DPanicw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.DPanicw(msg, keysAndValues...)
}

func (l *logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Panicw(msg, keysAndValues...)
}

func (l *logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Fatalw(msg, keysAndValues...)
}

func NewCore(jsonEncoder bool, writer io.Writer, debug bool) zapcore.Core {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	if jsonEncoder {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	writeSyncer := zapcore.AddSync(writer)

	level := zapcore.InfoLevel
	if debug {
		level = zapcore.DebugLevel
	}

	return zapcore.NewCore(encoder, writeSyncer, level)
}

func NewLogger(cores ...zapcore.Core) Logger {
	teeCore := zapcore.NewTee(cores...)
	return &logger{
		zapLogger: zap.New(
			teeCore,
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zapcore.FatalLevel+1),
		).Sugar(),
	}
}
