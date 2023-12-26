package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"trip_service/internal/app"

	"trip_service/internal/config"
)

func ParseConfig() *config.Config {
	var filepath string
	flag.StringVar(&filepath, "c", ".config.yaml", "set config path")
	flag.Parse()

	cfg, err := config.NewConfig(filepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return cfg
}

func main() {
	cfg := ParseConfig()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	tripApp := app.NewApp(cfg, ctx)
	tripApp.Run(ctx)
	<-ctx.Done()
	ctx, _ = context.WithTimeout(ctx, 3*time.Second)
}
