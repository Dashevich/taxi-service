package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"offering_service/internal/config"
	"offering_service/internal/server"
)

type App struct {
	cfg       *config.Config
	serverApp *server.Server
}

func NewApp() App {
	var filepath string
	flag.StringVar(&filepath, "c", ".config.yaml", "set config path")
	flag.Parse()

	cfg, err := config.NewConfig(filepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	serverApp := server.NewServer(cfg)
	return App{cfg, serverApp}
}

func (app *App) Running() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	defer func() {
		v := recover()
		if v != nil {
			ctx, _ := context.WithTimeout(ctx, 3*time.Second)
			app.serverApp.Stop(ctx)
			fmt.Println(v)
			os.Exit(1)
		}
	}()

	app.serverApp.Run()
	<-ctx.Done()
	ctx, _ = context.WithTimeout(ctx, 3*time.Second)
	app.serverApp.Stop(ctx)
}

func main() {
	app := NewApp()
	app.Running()
}
