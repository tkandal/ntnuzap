// Package ntnuzap is a custom logging configuration for uber zap, adjusted to NTNU's needs
package ntnuzap

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
    "os"
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
func NTNUConfig(files []string, level zapcore.Level, development bool, utc bool) zap.Config {
    cfg := zap.Config{
        Level:            zap.NewAtomicLevelAt(level),
        Development:      development,
        Encoding:         "json",
        EncoderConfig:    NTNUEncoderConfig(utc),
        OutputPaths:      []string{"stderr"},
        ErrorOutputPaths: []string{"stderr"},
    }
    if len(files) > 0 {
        cfg.OutputPaths = files
    }
    return cfg
}

// NTNUZap builds an custom logger for uber zap that logs to stderr.
//
// level is logging level
// developement ...
// utc set to true if the timestamps should be UTC time
func NTNUZap(files []string, level zapcore.Level, development bool, utc bool) (*zap.Logger, error) {
    return NTNUConfig(files, level, development, utc).Build()
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

    // First, define our level-handling logic.
    highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
        return lvl >= zapcore.ErrorLevel
    })

    infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
        return lvl >= zapcore.InfoLevel
    })

    lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
        return lvl < zapcore.ErrorLevel
    })

    // High-priority output should also go to standard error, and low-priority
    // output should also go to standard out.
    consoleDebugging := zapcore.Lock(os.Stdout)
    consoleErrors := zapcore.Lock(os.Stderr)

    // Encode as JSON to all endpoints
    consoleStdErr := zapcore.NewCore(zapcore.NewJSONEncoder(NTNUEncoderConfig(utc)), consoleErrors, highPriority)
    consoleStdout := zapcore.NewCore(zapcore.NewJSONEncoder(NTNUEncoderConfig(utc)), consoleDebugging, lowPriority)
    rolling := zapcore.NewCore(zapcore.NewJSONEncoder(NTNUEncoderConfig(utc)), w, infoLevel)

    // Join the outputs, encoders, and level-handling functions into
    // zapcore.Cores, then tee the four cores together.
    core := zapcore.NewTee(consoleStdErr, consoleStdout, rolling)

    logger := zap.New(core)

    return logger.WithOptions(zap.AddCaller()), nil
}
