package app

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// initDatabase инициализирует подключение к базе данных
func (a *App) initDatabase() error {
	connStr := a.config.GetDBConnString()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	a.db = db
	return nil
}

// runMigrations запускает миграции базы данных
func (a *App) runMigrations() error {
	driver, err := postgres.WithInstance(a.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		a.logger.Error("Failed to create migration instance", zap.Error(err))
		// Продолжим работу, даже если есть проблемы с миграциями
		return nil
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		a.logger.Error("Failed to apply migrations", zap.Error(err))
		// Проверим, есть ли таблица songs
		var exists bool
		err = a.db.QueryRow(checkTableExistsQuery).Scan(&exists)

		if err != nil {
			a.logger.Error("Failed to check if songs table exists", zap.Error(err))
		} else if !exists {
			// Если таблицы нет, создаем её вручную
			_, err = a.db.Exec(createSongsTableQuery)
			if err != nil {
				a.logger.Error("Failed to create songs table", zap.Error(err))
				return fmt.Errorf("failed to create songs table: %w", err)
			}
			a.logger.Info("Created songs table manually")
		}
	}
	return nil
}

// createDatabaseIfNotExists проверяет существование базы данных и создает её, если она не существует
func (a *App) createDatabaseIfNotExists() error {
	// Подключаемся к базе postgres для создания нашей базы данных
	postgresConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		a.config.DBHost,
		a.config.DBPort,
		a.config.DBUser,
		a.config.DBPassword,
	)

	postgresDB, err := sql.Open("postgres", postgresConnStr)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres database: %w", err)
	}
	defer postgresDB.Close()

	// Проверяем существование базы данных
	var exists bool
	err = postgresDB.QueryRow(checkDatabaseExistsQuery, a.config.DBName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	// Если база данных не существует, создаем её
	if !exists {
		_, err = postgresDB.Exec(fmt.Sprintf("CREATE DATABASE %s", a.config.DBName))
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		a.logger.Info(fmt.Sprintf("Database %s created successfully", a.config.DBName))
	}

	return nil
}
