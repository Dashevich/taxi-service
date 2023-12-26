package main

import (
	"client/app"
	"client/internal/config"
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"
)

func handleRecover(ctx context.Context) {
	if r := recover(); r != nil {
		log.Println("recovered panic from %+v, stack: %+v", r, string(debug.Stack()))
	}
}

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "cfg", "./client/cmd/.config.yaml", "set cfg path")
	flag.Parse()

	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		fmt.Println(fmt.Errorf("fatal: init config %w", err))
		log.Fatal()
	}

	a := app.NewApp(cfg) // server

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	go func() {
		defer handleRecover(ctx)
	}()

	a.Run()

	<-ctx.Done()
	ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	a.Stop(ctx)
}
