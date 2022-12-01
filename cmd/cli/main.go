package main

import (
	"context"
	"fmt"
	"github.com/Zzarin/file_parser/internal"
	"github.com/Zzarin/file_parser/internal/local"
	"github.com/Zzarin/file_parser/internal/postgres"
	"github.com/jessevdk/go-flags"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	var cfg Config
	parser := flags.NewParser(&cfg, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		log.Fatal("Error parse env variables", err)
	}

	logger := InitializeLoger()
	defer func() {
		err := logger.Sync()
		// write error: sync /dev/stdout: invalid argument #328
		if err != nil && err.Error() != "sync /dev/stdout: invalid argument" && err.Error() != "sync /dev/stdout: The handle is invalid." { //ignoring error
			logger.Error("erase logs from buffer", zap.Error(err))
		}
	}()

	dbConnection, err := NewDb(cfg.DbDsn, cfg.DbConnMaxLifetime, cfg.DbMaxOpenConns, cfg.DbMaxIdleConns)
	if err != nil {
		logger.Fatal("get DB instance", zap.Error(err))
	}
	defer func() {
		err := dbConnection.Close()
		if err != nil {
			logger.Error("closing dbConnection", zap.Error(err))
		}
	}()
	logger.Info("Successful dbConnection")

	parserStorage := postgres.NewStorage(dbConnection)

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	err = parserStorage.Init(ctx)
	if err != nil {
		logger.Error("init new table", zap.Error(err))
	}
	localFile := local.NewLocalFile(cfg.FileName, cfg.StrReadingFreq)
	appParser := internal.NewParser(parserStorage, logger)

	err = appParser.ParseStruct(ctx, localFile)
	if err != nil {
		logger.Error("parse structure", zap.Error(err))
	}

	go func() {
		<-ctx.Done()
		logger.Info("Shutdown signal received")
	}()
	fmt.Println("app finished")
}
