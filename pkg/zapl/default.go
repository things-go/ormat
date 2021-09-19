package zapl

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

var defaultSugar = zap.NewNop().Sugar()

func ReplaceGlobals(s *zap.SugaredLogger) {
	defaultSugar = s
}

func Desugar() *zap.Logger {
	return defaultSugar.Desugar()
}

func With(args ...interface{}) *zap.SugaredLogger {
	return defaultSugar.With(args...)
}

func Named(name string) *zap.SugaredLogger {
	return defaultSugar.Named(name)
}

func Sync() error {
	return defaultSugar.Sync()
}

func Debug(args ...interface{}) {
	defaultSugar.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	defaultSugar.Debugf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	defaultSugar.Debugw(msg, keysAndValues...)
}

func Info(args ...interface{}) {
	defaultSugar.Info(args...)
}

func Infof(template string, args ...interface{}) {
	defaultSugar.Infof(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	defaultSugar.Infow(msg, keysAndValues...)
}

func Warn(args ...interface{}) {
	defaultSugar.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	defaultSugar.Warnf(template, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	defaultSugar.Warnw(msg, keysAndValues...)
}

func Error(args ...interface{}) {
	defaultSugar.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	defaultSugar.Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	defaultSugar.Errorw(msg, keysAndValues...)
}

func Fatal(args ...interface{}) {
	defaultSugar.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	defaultSugar.Fatalf(template, args...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	defaultSugar.Fatalw(msg, keysAndValues...)
}

func DPanic(args ...interface{}) {
	defaultSugar.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	defaultSugar.DPanicf(template, args...)
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	defaultSugar.DPanicw(msg, keysAndValues...)
}

func Panic(args ...interface{}) {
	defaultSugar.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	defaultSugar.Panicf(template, args...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	defaultSugar.Panicw(msg, keysAndValues...)
}

func JSON(v ...interface{}) {
	for _, vv := range v {
		b, _ := json.MarshalIndent(vv, "", "  ")
		fmt.Println(string(b))
	}
}
