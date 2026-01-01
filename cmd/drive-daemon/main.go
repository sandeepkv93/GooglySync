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

	fmt.Printf("drive-daemon starting (config=%q log-level=%q)\n", *configPath, *logLevel)
	os.Exit(0)
}
