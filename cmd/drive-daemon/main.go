package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sandeepkv93/googlysync/internal/config"
)

var version = "dev"

func main() {
	configPath := flag.String("config", "", "path to config file (JSON)")
	logLevel := flag.String("log-level", "", "log level")
	socketPath := flag.String("socket", "", "unix socket path")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}

	opts := config.Options{
		ConfigPath: *configPath,
		LogLevel:   *logLevel,
		SocketPath: *socketPath,
	}

	daemon, err := InitializeDaemon(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "init failed: %v\n", err)
		os.Exit(1)
	}
	if daemon.Logger != nil {
		defer daemon.Logger.Sync()
	}
	if daemon.IPC != nil {
		daemon.IPC.WithVersion(version)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := daemon.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "run failed: %v\n", err)
		os.Exit(1)
	}
}
