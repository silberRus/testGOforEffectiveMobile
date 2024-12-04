package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/testTask/docs"
	"github.com/testTask/internal/app"
	"github.com/testTask/internal/config"
	"go.uber.org/zap"
)

// @title Music Library API
// @version 1.0
// @description API for managing music library
// @host localhost:8080
// @BasePath /api/v1
func main() {

	// Загружаем переменные из .env файла
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Отображаем переменные окружения на всякий случай для админа
	fmt.Println("DB_USER:", os.Getenv("DB_USER"))
	fmt.Println("DB_PASSWORD:", os.Getenv("DB_PASSWORD"))
	fmt.Println("DB_HOST:", os.Getenv("DB_HOST"))
	fmt.Println("DB_PORT:", os.Getenv("DB_PORT"))
	fmt.Println("DB_NAME:", os.Getenv("DB_NAME"))

	// Инициализируем логгер
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	// у логгеров тоже бывают проблемы
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to sync logger: %v\n", err)
		}
	}(logger)

	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Создаем приложение
	application := app.New(cfg, logger)

	// Инициализируем компоненты
	if err := application.Initialize(); err != nil {
		logger.Fatal("Failed to initialize application", zap.Error(err))
	}

	// Запускаем приложение в горутине
	go func() {
		if err := application.Run(); err != nil {
			logger.Fatal("Failed to run application", zap.Error(err))
		}
	}()

	// Ожидаем сигнал для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown за 30 секунд более надежнее
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Останавливаем приложение
	if err := application.Shutdown(ctx); err != nil {
		logger.Fatal("Failed to shutdown application", zap.Error(err))
	}

	logger.Info("Application stopped")
}
