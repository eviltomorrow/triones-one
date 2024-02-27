package conf

import (
	"fmt"
	"os"
	"path/filepath"

	"triones-one/lib/config"
	"triones-one/lib/system"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
)

type Notify interface {
	Reload(*Config) error
}

type Config struct {
	Repository Repository `mapstructure:"repository" json:"repository"`
	SQLite     SQLite     `mapstructure:"sqlite" json:"sqlite"`

	Server config.Server `mapstructure:"server" json:"server"`
	Log    config.Log    `mapstructure:"log" json:"log"`
	Etcd   config.Etcd   `mapstructure:"etcd" json:"etcd"`
}

type Repository struct {
	ImagesPath string `mapstructure:"images-path" json:"images-path"`
}

type SQLite struct {
	DBPath string `mapstructure:"db-path" json:"db-path"`
}

func (c *Config) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(c)
	return string(buf)
}

func ReadConfigFile(path string) (*Config, error) {
	findConfigFile := func() (string, error) {
		for _, p := range []string{
			path,
			filepath.Join(system.Runtime.RootDir, "/etc/config.toml"),
			filepath.Join("/etc/config.toml"),
		} {
			fi, err := os.Stat(path)
			if os.IsNotExist(err) {
				continue
			}
			if err != nil {
				return "", err
			}

			if fi.IsDir() {
				continue
			}

			p, err = filepath.Abs(p)
			if err != nil {
				return "", err
			}
			return p, nil
		}
		return "", fmt.Errorf("not found config-file")
	}

	var err error
	path, err = findConfigFile()
	if err != nil {
		return nil, err
	}
	return loadConfig(path)
}

func loadConfig(path string) (*Config, error) {
	config := &Config{
		Repository: Repository{
			ImagesPath: filepath.Join(system.Runtime.RootDir, "/data/images"),
		},
		SQLite: SQLite{
			DBPath: filepath.Join(system.Runtime.RootDir, "/data/db"),
		},
		Etcd: config.Etcd{
			Endpoints: []string{"127.0.0.1:2379"},
		},
		Log: config.Log{
			DisableTimestamp: false,
			Level:            "info",
			MaxSize:          100,
			MaxDays:          90,
			MaxBackups:       180,
			Compress:         true,
			FileDir:          filepath.Join(system.Runtime.RootDir, "/var/log"),
		},
		Server: config.Server{
			GrpcAddressList: []string{"0.0.0.0", "::"},
			GrpcPort:        50001,
			TlsCaFile:       filepath.Join(system.Runtime.RootDir, "/etc/certs/ca.crt"),
			TlsServerCert:   filepath.Join(system.Runtime.RootDir, "/var/certs/server.crt"),
			TlsServerPem:    filepath.Join(system.Runtime.RootDir, "/var/certs/server.pem"),
		},
	}

	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("toml")

	var err error
	if err = v.ReadInConfig(); err != nil {
		return nil, err
	}
	if err = v.UnmarshalExact(config); err != nil {
		return nil, err
	}
	return config, nil
}

func InitalConfig(currentConfig *Config) (*Config, error) {
	return nil, nil
}
