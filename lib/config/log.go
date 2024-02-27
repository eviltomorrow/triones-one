package config

import (
	"triones-one/lib/zlog"

	jsoniter "github.com/json-iterator/go"
)

type Log struct {
	DisableTimestamp bool   `mapstructure:"disable-timestamp" json:"disable-timestamp"`
	Level            string `mapstructure:"level" json:"level"`
	Format           string `mapstructure:"format" json:"format"`
	FileDir          string `mapstructure:"file-dir" json:"file-dir"`
	MaxSize          int    `mapstructure:"maxsize" json:"maxsize"`
	MaxDays          int    `mapstructure:"max-days" json:"max-days"`
	MaxBackups       int    `mapstructure:"max-backups" json:"max-backups"`
	Compress         bool   `mapstructure:"compress" json:"compress"`
}

func (c *Log) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(c)
	return string(buf)
}

func SetupLogger(l Log) ([]func() error, error) {
	global, prop, err := zlog.InitLogger(&zlog.Config{
		Level:            l.Level,
		Format:           l.Format,
		DisableTimestamp: l.DisableTimestamp,
		File: zlog.FileLogConfig{
			Filename:   l.FileDir,
			MaxSize:    l.MaxSize,
			MaxDays:    l.MaxDays,
			MaxBackups: l.MaxBackups,
			Compress:   l.Compress,
		},
		DisableStacktrace:   true,
		DisableErrorVerbose: true,
	})
	if err != nil {
		return nil, err
	}
	zlog.ReplaceGlobals(global, prop)

	return []func() error{
		func() error { return global.Sync() },
	}, nil
}
