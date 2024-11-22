package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"

	"github.com/topi314/shelly-alert-killswitch/kswitch"
)

var (
	Version = "dev"
	Commit  = "unknown"
)

func main() {
	cfgPath := flag.String("config", "config.toml", "Path to config file")
	flag.Parse()

	cfg, err := kswitch.LoadConfig(*cfgPath)
	if err != nil {
		slog.Error("Failed to load config", slog.Any("err", err))
		return
	}

	setupLogger(cfg.Log)

	slog.Info("Starting Shelly Alert Killswitch...", slog.String("version", Version), slog.String("commit", Commit), slog.String("config", *cfgPath), slog.Any("cfg", cfg))

	server := kswitch.New(cfg, Version)
	server.SetupRoutes()

	go server.Start()

	defer server.Stop()

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)

	slog.Info("Started Shelly Alert Killswitch", slog.String("addr", cfg.Server.ListenAddr))

	<-s
}

func setupLogger(cfg kswitch.LogConfig) {
	opts := &slog.HandlerOptions{
		AddSource: cfg.AddSource,
		Level:     cfg.Level,
	}
	var handler slog.Handler
	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}
	slog.SetDefault(slog.New(handler))
}
