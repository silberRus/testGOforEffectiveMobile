package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/testTask/internal/config"
	"github.com/testTask/internal/handlers"
	"github.com/testTask/internal/middleware"
	"github.com/testTask/internal/repository"
	"github.com/testTask/internal/service"
	"go.uber.org/zap"
)

// App структура приложения
type App struct {
	config     *config.Config
	logger     *zap.Logger
	db         *sql.DB
	httpServer *http.Server
}

// New конструктор нового экземпляра приложения
func New(cfg *config.Config, logger *zap.Logger) *App {
	return &App{
		config: cfg,
		logger: logger,
	}
}

// Initialize инициализирует компоненты приложения
func (a *App) Initialize() error {
	// Создаем базу данных, если она не существует
	if err := a.createDatabaseIfNotExists(); err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	if err := a.initDatabase(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	if err := a.runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err := a.initHTTPServer(); err != nil {
		return fmt.Errorf("failed to initialize HTTP server: %w", err)
	}

	return nil
}

// initHTTPServer инициализирует HTTP сервер
func (a *App) initHTTPServer() error {

	// Инициализируем репозиторий, сервис и обработчики
	repo := repository.NewPostgresSongRepository(a.db)
	svc := service.NewSongService(repo, a.logger)
	handler := handlers.NewSongHandler(svc, a.logger)

	// Создаем роутер и регистрируем маршруты
	r := mux.NewRouter()

	// Добавляем middleware для логирования
	r.Use(middleware.LoggingMiddleware(a.logger))

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/songs", handler.GetSongs).Methods(http.MethodGet)
	api.HandleFunc("/songs/{id}/lyrics", handler.GetLyrics).Methods(http.MethodGet)
	api.HandleFunc("/songs", handler.CreateSong).Methods(http.MethodPost)
	api.HandleFunc("/songs/{id}", handler.UpdateSong).Methods(http.MethodPut)
	api.HandleFunc("/songs/{id}", handler.DeleteSong).Methods(http.MethodDelete)

	// Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Создаем HTTP сервер
	a.httpServer = &http.Server{
		Addr:         ":" + a.config.ServerPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return nil
}

// Run запуск приложения
func (a *App) Run() error {
	a.logger.Info("Starting server", zap.String("port", a.config.ServerPort))
	if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

// Shutdown gracefully останавливает приложение
func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("Shutting down server...")

	if err := a.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	if err := a.db.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}
