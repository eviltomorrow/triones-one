package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"triones-one/apps/triones-repository/conf"
	"triones-one/lib/buildinfo"

	"github.com/jessevdk/go-flags"
)

var (
	AppName     = "unknown"
	MainVersion = "unknown"
	GitSha      = "unknown"
	BuildTime   = "unknown"
)

func init() {
	buildinfo.AppName = AppName
	buildinfo.MainVersion = MainVersion
	buildinfo.GitSha = GitSha
	buildinfo.BuildTime = BuildTime
}

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	var opts struct {
		ConfigFile string `short:"f" long:"config-file" description:"specifying a config file"`
		Version    bool   `long:"version" description:"show version number"`
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	if opts.Version {
		fmt.Println(buildinfo.Version())
		os.Exit(0)
	}

	initialConfig, err := conf.ReadConfigFile(opts.ConfigFile)
	if err != nil {
		log.Printf("[F] Read config-file failure, nest error: %v", err)
		os.Exit(1)
	}
	currentConfig, err := conf.InitalConfig(initialConfig)
	_ = currentConfig
}
