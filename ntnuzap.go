// Package ntnuzap is a custom logging configuration for uber zap, adjusted to NTNU's needs
package ntnuzap

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
    "time"
)

/*
 * Copyright (c) 2019 Norwegian University of Science and Technology
 */

// UTCTimeEncoder encodes timestamps to UTC for uber zap
func UTCTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
    enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.000Z0700"))
}

// NTNUEncoderConfig is a custom encoder configuration for uber zap logging
func NTNUEncoderConfig() zapcore.EncoderConfig {
    return zapcore.EncoderConfig{
        // Keys can be anything except the empty string.
        TimeKey:        "ts",
        LevelKey:       "level",
        NameKey:        "logger",
        CallerKey:      "caller",
        MessageKey:     "msg",
        StacktraceKey:  "stacktrace",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    zapcore.CapitalLevelEncoder,
        EncodeTime:     UTCTimeEncoder,
        EncodeDuration: zapcore.StringDurationEncoder,
        EncodeCaller:   zapcore.ShortCallerEncoder,
        EncodeName:     zapcore.FullNameEncoder,
    }
}

// NTNUConfig is a custom configuration for uber zap logging
func NTNUConfig(level zapcore.Level, development bool) zap.Config {
    return zap.Config{
        Level:            zap.NewAtomicLevelAt(level),
        Development:      development,
        Encoding:         "json",
        EncoderConfig:    NTNUEncoderConfig(),
        OutputPaths:      []string{"stderr"},
        ErrorOutputPaths: []string{"stderr"},
    }
}

// NTNUZap builds an custom logger for uber zap that logs to stderr
func NTNUZap(level zapcore.Level, development bool) (*zap.Logger, error) {
    return NTNUConfig(level, development).Build()
}

// NTNULumberjack builds a custom rotating file-logger for uber zap, logfile is destination,
// maxSize is file-size in megabytes, maxBack is max number of backups and maxAge is
// max number of days
func NTNULumberjack(logfile string, maxSize int, maxBack int, maxAge int) (*zap.Logger, error) {
    w := zapcore.AddSync(&lumberjack.Logger{
        Filename:   logfile,
        MaxSize:    maxSize, // megabytes
        MaxBackups: maxBack,
        MaxAge:     maxAge, // days
    })
    core := zapcore.NewCore(
        zapcore.NewJSONEncoder(NTNUEncoderConfig()),
        w,
        zap.InfoLevel,
    )
    logger := zap.New(core)
    return logger.WithOptions(zap.AddCaller()), nil
}
