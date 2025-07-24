package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/inkxk/bad-bot/app"
	"github.com/inkxk/bad-bot/app/linebot"
	"github.com/inkxk/bad-bot/config"
	zap "github.com/inkxk/bad-bot/logger"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	lineSdk "github.com/line/line-bot-sdk-go/v8/linebot"
)

func main() {
	// config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Panicf("Error config: %v", err)
	}

	// logger
	logger, undo := zap.NewZap(cfg.LogLevel)
	defer undo()

	// app
	r := app.NewRouter()
	r.Use(echoMiddleware.Recover())
	r.Use(echoMiddleware.RequestID())
	r.Use(zap.LoggerToContextMiddleware(logger))
	r.Use(zap.ZapLoggerMiddleware(logger))

	// create instance line bot
	botClient, err := lineSdk.New(cfg.Line.LineChannelSecret, cfg.Line.LineChanneAccessToken)
	if err != nil {
		logger.Sugar().Fatalf("Failed to create LINE bot: %v", err)
	}

	// handler
	lineHandler := linebot.NewHandler(botClient, logger)

	// router
	r.HealthCheck()
	r.POST("/callback", lineHandler.Callback)

	// server
	srv := http.Server{
		Addr:              ":" + cfg.HTTPServer.Port,
		Handler:           r,
		ReadHeaderTimeout: cfg.HTTPServer.ReadHeaderTimeout,
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		logger.Sugar().Infof("shutting down in %s...", cfg.GracefulTimeout)
		ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefulTimeout)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Sugar().Info("HTTP server Shutdown: " + err.Error())
		}
		close(idleConnsClosed)
	}()

	logger.Info(":" + cfg.HTTPServer.Port + " is serve")

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error("HTTP server ListenAndServe: " + err.Error())
		return
	}

	<-idleConnsClosed
	logger.Info("Server shutdown gracefully")
}
