package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"strings"
)

// Configuration for logging
type Config struct {
	ConsoleEnabled bool
	ConsoleLevel   string
	ConsoleJson    bool

	FileEnabled bool
	FileLevel   string
	FileJson    bool

	// Directory to log to to when filelogging is enabled
	Directory string
	// Filename is the name of the logfile which will be placed inside the directory
	Filename string
	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int
	// MaxBackups the max number of rolled files to keep
	MaxBackups int
	// MaxAge the max age in days to keep a logfile
	MaxAge int
}

// DefaultZapLogger is the default logger instance that should be used to log
// It's assigned a default value here for tests (which do not call log.Configure())
var DefaultZapLogger *Logger

type Logger struct {
	Unsugared *zap.Logger
	*zap.SugaredLogger
}

func InitializeDefaultZapLogger() {
	DefaultZapLogger = newZapLogger(Config{
		FileEnabled: true,
		FileJson:    true,
		FileLevel:   DefaultLogLvl,
		Directory:   DefaultLogDir,
		Filename:    DefaultLogFileName,
		MaxSize:     DefaultLogFileSizeMB,
		MaxAge:      DefaultLogFileAgeDays,
	})
}

// Debugt Use zap.String(key, value), zap.Int(key, value) to log fields. These fields
// will be marshalled as JSON in the logfile and key value pairs in the console!
func (logger *Logger) Debugt(msg string, fields ...zapcore.Field) {
	logger.Unsugared.Debug(msg, fields...)
}

func (logger *Logger) Debug(msg string, keysAndValues ...interface{}) {
	logger.SugaredLogger.Debugw(msg, keysAndValues...)
}

func (logger *Logger) Debugs(args ...interface{}) {
	logger.SugaredLogger.Debug(args...)
}

// Infot Use zap.String(key, value), zap.Int(key, value) to log fields. These fields
// will be marshalled as JSON in the logfile and key value pairs in the console!
func (logger *Logger) Infot(msg string, fields ...zapcore.Field) {
	logger.Unsugared.Info(msg, fields...)
}

func (logger *Logger) Info(msg string, keysAndValues ...interface{}) {
	logger.SugaredLogger.Infow(msg, keysAndValues...)
}

func (logger *Logger) Infos(args ...interface{}) {
	logger.SugaredLogger.Info(args...)
}

// Warnt Use zap.String(key, value), zap.Int(key, value) to log fields. These fields
// will be marshalled as JSON in the logfile and key value pairs in the console!
func (logger *Logger) Warnt(msg string, fields ...zapcore.Field) {
	logger.Unsugared.Warn(msg, fields...)
}

func (logger *Logger) Warn(msg string, keysAndValues ...interface{}) {
	logger.SugaredLogger.Warnw(msg, keysAndValues...)
}

func (logger *Logger) Warns(args ...interface{}) {
	logger.SugaredLogger.Warn(args...)
}

// Errort Use zap.String(key, value), zap.Int(key, value) to log fields. These fields
// will be marshalled as JSON in the logfile and key value pairs in the console!
func (logger *Logger) Errort(msg string, fields ...zapcore.Field) {
	logger.Unsugared.Error(msg, fields...)
}

func (logger *Logger) Error(msg string, keysAndValues ...interface{}) {
	logger.SugaredLogger.Errorw(msg, keysAndValues...)
}

func (logger *Logger) Errors(args ...interface{}) {
	logger.SugaredLogger.Error(args...)
}

// Panict Use zap.String(key, value), zap.Int(key, value) to log fields. These fields
// will be marshalled as JSON in the logfile and key value pairs in the console!
func (logger *Logger) Panict(msg string, fields ...zapcore.Field) {
	logger.Unsugared.Panic(msg, fields...)
}

func (logger *Logger) Panic(msg string, keysAndValues ...interface{}) {
	logger.SugaredLogger.Panicw(msg, keysAndValues...)
}

func (logger *Logger) Panics(args ...interface{}) {
	logger.SugaredLogger.Panic(args...)
}

// Fatalt Use zap.String(key, value), zap.Int(key, value) to log fields. These fields
// will be marshalled as JSON in the logfile and key value pairs in the console!
func (logger *Logger) Fatalt(msg string, fields ...zapcore.Field) {
	logger.Unsugared.Fatal(msg, fields...)
}

func (logger *Logger) Fatal(msg string, keysAndValues ...interface{}) {
	logger.SugaredLogger.Fatalw(msg, keysAndValues...)
}

func (logger *Logger) Fatals(args ...interface{}) {
	logger.SugaredLogger.Fatal(args...)
}

func Debugt(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.Debugt(msg, fields...)
}

func Debugf(template string, args ...interface{}) {
	DefaultZapLogger.Debugf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	DefaultZapLogger.Debugw(msg, keysAndValues...)
}

func Debug(msg string, keysAndValues ...interface{}) {
	DefaultZapLogger.Debug(msg, keysAndValues...)
}

func Debugs(args ...interface{}) {
	DefaultZapLogger.Debugs(args...)
}

func Infot(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.Infot(msg, fields...)
}

func Infof(template string, args ...interface{}) {
	DefaultZapLogger.Infof(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	DefaultZapLogger.Infow(msg, keysAndValues...)
}

func Info(msg string, keysAndValues ...interface{}) {
	DefaultZapLogger.Info(msg, keysAndValues...)
}

func Infos(args ...interface{}) {
	DefaultZapLogger.Infos(args...)
}

func Warnt(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.Warnt(msg, fields...)
}

func Warnf(template string, args ...interface{}) {
	DefaultZapLogger.Warnf(template, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	DefaultZapLogger.Warnw(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...interface{}) {
	DefaultZapLogger.Warn(msg, keysAndValues...)
}

func Warns(args ...interface{}) {
	DefaultZapLogger.Warns(args...)
}

func Errort(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.Errort(msg, fields...)
}

func Errorf(template string, args ...interface{}) {
	DefaultZapLogger.Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	DefaultZapLogger.Errorw(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	DefaultZapLogger.Error(msg, keysAndValues...)
}

func Errors(args ...interface{}) {
	DefaultZapLogger.Errors(args...)
}

func Panict(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.Panict(msg, fields...)
}

func Panicf(template string, args ...interface{}) {
	DefaultZapLogger.Panicf(template, args...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	DefaultZapLogger.Panicw(msg, keysAndValues...)
}

func Panic(msg string, keysAndValues ...interface{}) {
	DefaultZapLogger.Panic(msg, keysAndValues...)
}

func Panics(args ...interface{}) {
	DefaultZapLogger.Panics(args...)
}

func Fatalt(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.Fatalt(msg, fields...)
}

func Fatalf(template string, args ...interface{}) {
	DefaultZapLogger.Fatalf(template, args...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	DefaultZapLogger.Fatalw(msg, keysAndValues...)
}

func Fatal(msg string, keysAndValues ...interface{}) {
	DefaultZapLogger.Fatal(msg, keysAndValues...)
}

func Fatals(args ...interface{}) {
	DefaultZapLogger.Fatals(args...)
}

func newRollingFile(config Config) zapcore.WriteSyncer {
	if err := os.MkdirAll(config.Directory, 0744); err != nil {
		Error("can't create log directory", zap.Error(err), zap.String("path", config.Directory))
		return nil
	}

	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
		MaxBackups: config.MaxBackups, // files
	})
}

func newZapLogger(config Config) *Logger {
	var consoleLevel zapcore.Level
	consoleLevel.Set(strings.ToLower(config.ConsoleLevel))

	var fileLevel zapcore.Level
	fileLevel.Set(strings.ToLower(config.FileLevel))

	consoleEncCfg := zapcore.EncoderConfig{
		TimeKey:        "@timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
	}
	jsonEncCfg := zapcore.EncoderConfig{
		TimeKey:        "@timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
	}

	consoleLevelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= consoleLevel
	})
	fileLevelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= fileLevel
	})

	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncCfg)
	fileEncoder := zapcore.NewJSONEncoder(jsonEncCfg)

	var cores []zapcore.Core

	if config.ConsoleEnabled {
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stderr), consoleLevelEnabler))
	}
	if config.FileEnabled {
		cores = append(cores, zapcore.NewCore(fileEncoder, newRollingFile(config), fileLevelEnabler))
	}
	core := zapcore.NewTee(cores...)

	unsugared := zap.New(core)
	return &Logger{
		Unsugared:     unsugared,
		SugaredLogger: unsugared.Sugar(),
	}
}
