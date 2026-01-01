package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sandeepkv93/googlysync/internal/config"
	"github.com/sandeepkv93/googlysync/internal/ipc"
	ipcgen "github.com/sandeepkv93/googlysync/internal/ipc/gen"
)

var version = "dev"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "daemon":
		runDaemon(os.Args[2:])
	case "ping":
		runPing(os.Args[2:])
	case "status":
		runStatus(os.Args[2:])
	case "fuse":
		runFuse(os.Args[2:])
	case "version":
		fmt.Println(version)
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Println("Usage: googlysync <command> [options]")
	fmt.Println("Commands:")
	fmt.Println("  daemon   Start the sync daemon")
	fmt.Println("  ping     Ping the daemon and print version")
	fmt.Println("  status   Print daemon sync status")
	fmt.Println("  fuse     Placeholder for streaming mode")
	fmt.Println("  version  Print CLI version")
}

func runDaemon(args []string) {
	fs := flag.NewFlagSet("daemon", flag.ExitOnError)
	configPath := fs.String("config", "", "path to config file (JSON)")
	logLevel := fs.String("log-level", "", "log level")
	socketPath := fs.String("socket", "", "unix socket path")
	_ = fs.Parse(args)

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

func runPing(args []string) {
	fs := flag.NewFlagSet("ping", flag.ExitOnError)
	socketPath := fs.String("socket", "", "unix socket path")
	_ = fs.Parse(args)

	cfg, err := config.NewConfigWithOptions(config.Options{SocketPath: *socketPath})
	if err != nil {
		fmt.Printf("config error: %v\n", err)
		return
	}
	conn, err := ipc.Dial(context.Background(), cfg.SocketPath)
	if err != nil {
		fmt.Printf("dial error: %v\n", err)
		return
	}
	defer conn.Close()

	client := ipcgen.NewDaemonControlClient(conn)
	resp, err := client.Ping(context.Background(), &ipcgen.Empty{})
	if err != nil {
		fmt.Printf("ping error: %v\n", err)
		return
	}
	fmt.Println(resp.Version)
}

func runStatus(args []string) {
	fs := flag.NewFlagSet("status", flag.ExitOnError)
	socketPath := fs.String("socket", "", "unix socket path")
	_ = fs.Parse(args)

	cfg, err := config.NewConfigWithOptions(config.Options{SocketPath: *socketPath})
	if err != nil {
		fmt.Printf("config error: %v\n", err)
		return
	}
	conn, err := ipc.Dial(context.Background(), cfg.SocketPath)
	if err != nil {
		fmt.Printf("dial error: %v\n", err)
		return
	}
	defer conn.Close()

	client := ipcgen.NewSyncStatusClient(conn)
	resp, err := client.GetStatus(context.Background(), &ipcgen.Empty{})
	if err != nil {
		fmt.Printf("status error: %v\n", err)
		return
	}
	fmt.Printf("%s: %s\n", resp.Status.State.String(), resp.Status.Message)
}

func runFuse(args []string) {
	_ = args
	fmt.Println("fuse placeholder: streaming mode not implemented")
}
