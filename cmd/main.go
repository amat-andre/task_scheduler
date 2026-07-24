package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"task_scheduler/internal/auth"
	"task_scheduler/internal/config"
	"task_scheduler/internal/db"
	"task_scheduler/internal/server"

	"github.com/joho/godotenv"
)

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Printf("Ошибка загрузки .env файла: %v", err)
    }
}

func main() {
	logger := log.New(log.Writer(), "", log.LstdFlags|log.Lshortfile)

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("failed to load config: %v", err)
	}

	auth.Init(cfg.JWT.Secret)

	err = db.Init(cfg.DB)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	server := server.New(cfg.Server, logger)

	go func(){
		err = server.Run()
		if err != nil{
			logger.Fatalf("failed to run server: %v", err)
		}
	}()
	
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Printf("failed to shutdown server: %v", err)
	}

	logger.Println("application stopped")
}
