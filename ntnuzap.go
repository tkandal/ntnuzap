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
//
// utc set to true if the timestamps should be UTC time
func NTNUEncoderConfig(utc bool) zapcore.EncoderConfig {
    enc := zapcore.EncoderConfig{
        // Keys can be anything except the empty string.
        TimeKey:        "time",
        LevelKey:       "level",
        NameKey:        "logger",
        CallerKey:      "module",
        MessageKey:     "msg",
        StacktraceKey:  "stack",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    zapcore.CapitalLevelEncoder,
        EncodeTime:     zapcore.ISO8601TimeEncoder,
        EncodeDuration: zapcore.StringDurationEncoder,
        EncodeCaller:   zapcore.ShortCallerEncoder,
        EncodeName:     zapcore.FullNameEncoder,
    }
    if utc {
        enc.EncodeTime = UTCTimeEncoder
    }
    return enc
}

// NTNUConfig is a custom configuration for uber zap logging.
// level is logging level
// developement ...
// utc set to true if the timestamps should be UTC time
func NTNUConfig(level zapcore.Level, development bool, utc bool) zap.Config {
    return zap.Config{
        Level:            zap.NewAtomicLevelAt(level),
        Development:      development,
        Encoding:         "json",
        EncoderConfig:    NTNUEncoderConfig(utc),
        OutputPaths:      []string{"stderr"},
        ErrorOutputPaths: []string{"stderr"},
    }
}

// NTNUZap builds an custom logger for uber zap that logs to stderr.
//
// level is logging level
// developement ...
// utc set to true if the timestamps should be UTC time
func NTNUZap(level zapcore.Level, development bool, utc bool) (*zap.Logger, error) {
    return NTNUConfig(level, development, utc).Build()
}

// NTNULumberjack builds a custom rotating file-logger for uber zap.
//
// logfile is file-destination
// maxSize is max file-size in megabytes
// maxBack is max number of backup-files
// maxAge is max number of days
// utc set to true if the timestamps should be UTC time
func NTNULumberjack(logfile string, maxSize int, maxBack int, maxAge int, utc bool) (*zap.Logger, error) {
    w := zapcore.AddSync(&lumberjack.Logger{
        Filename:   logfile,
        MaxSize:    maxSize, // megabytes
        MaxBackups: maxBack,
        MaxAge:     maxAge, // days
    })
    core := zapcore.NewCore(
        zapcore.NewJSONEncoder(NTNUEncoderConfig(utc)),
        w,
        zap.InfoLevel,
    )
    logger := zap.New(core)
    return logger.WithOptions(zap.AddCaller()), nil
}
