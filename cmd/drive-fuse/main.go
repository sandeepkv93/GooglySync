package main

import (
	"flag"
	"fmt"
)

var version = "dev"

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}

	fmt.Println("drive-fuse placeholder: streaming mode not implemented")
}
