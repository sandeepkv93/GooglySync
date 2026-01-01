package main

import (
	"flag"
	"fmt"
	"os"
)

var version = "dev"

func main() {
	configPath := flag.String("config", "", "path to config file")
	logLevel := flag.String("log-level", "info", "log level")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}

	if *configPath != "" {
		fmt.Printf("config override: %q (not yet wired)\n", *configPath)
	}
	if *logLevel != "info" {
		fmt.Printf("log-level override: %q (not yet wired)\n", *logLevel)
	}

	daemon, err := InitializeDaemon()
	if err != nil {
		fmt.Fprintf(os.Stderr, "init failed: %v\n", err)
		os.Exit(1)
	}
	if err := daemon.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "run failed: %v\n", err)
		os.Exit(1)
	}
}
