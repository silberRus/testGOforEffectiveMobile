package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/testTask/internal/errors"
	"github.com/testTask/internal/models"
)

type SongRepository interface {
	GetSongs(ctx context.Context, filter *models.SongFilter) (*models.SongsResponse, error)
	GetSongByID(ctx context.Context, id int) (*models.Song, error)
	CreateSong(ctx context.Context, song *models.Song) (*models.Song, error)
	UpdateSong(ctx context.Context, song *models.Song) (*models.Song, error)
	DeleteSong(ctx context.Context, id int) error
}

type PostgresSongRepository struct {
	db *sql.DB
}

func NewPostgresSongRepository(db *sql.DB) SongRepository {
	return &PostgresSongRepository{db: db}
}

// GetSongs получает список песен
func (r *PostgresSongRepository) GetSongs(ctx context.Context, filter *models.SongFilter) (*models.SongsResponse, error) {
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}

	// Получаем общее количество записей
	var totalItems int
	err := r.db.QueryRowContext(ctx, countSongsQuery,
		filter.GroupName,
		filter.SongName,
		filter.FromDate,
		filter.ToDate,
		filter.Text,
		filter.Link,
	).Scan(&totalItems)
	if err != nil {
		return nil, errors.NewInternal("failed to count songs", err)
	}

	// Вычисляем общее количество страниц
	totalPages := (totalItems + filter.PageSize - 1) / filter.PageSize

	// Если запрошенная страница больше общего количества страниц, возвращаем ошибку
	if filter.Page > totalPages {
		return nil, errors.NewNotFound(fmt.Sprintf("page %d does not exist, total pages: %d", filter.Page, totalPages), nil)
	}

	offset := (filter.Page - 1) * filter.PageSize

	// Получаем записи для текущей страницы
	rows, err := r.db.QueryContext(ctx, getSongsQuery,
		filter.GroupName,
		filter.SongName,
		filter.FromDate,
		filter.ToDate,
		filter.Text,
		filter.Link,
		filter.PageSize,
		offset,
	)
	if err != nil {
		return nil, errors.NewInternal("failed to query songs", err)
	}
	defer rows.Close()

	// Собираем список песен
	songs := make([]models.Song, 0)
	for rows.Next() {
		var song models.Song
		err := rows.Scan(
			&song.ID,
			&song.GroupName,
			&song.SongName,
			&song.ReleaseDate,
			&song.Text,
			&song.Link,
			&song.CreatedAt,
			&song.UpdatedAt,
		)
		if err != nil {
			return nil, errors.NewInternal("failed to scan song", err)
		}
		songs = append(songs, song)
	}

	// Возвращаем список песен
	return &models.SongsResponse{
		Songs:       songs,
		CurrentPage: filter.Page,
		TotalPages:  totalPages,
		TotalItems:  totalItems,
		PageSize:    filter.PageSize,
	}, nil
}

// GetSongByID получает информацию о песне по ее ID
func (r *PostgresSongRepository) GetSongByID(ctx context.Context, id int) (*models.Song, error) {
	var song models.Song
	err := r.db.QueryRowContext(ctx, getSongByIDQuery, id).Scan(
		&song.ID,
		&song.GroupName,
		&song.SongName,
		&song.ReleaseDate,
		&song.Text,
		&song.Link,
		&song.CreatedAt,
		&song.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.NewNotFound("song not found", err)
	}
	if err != nil {
		return nil, errors.NewInternal("failed to get song", err)
	}

	return &song, nil
}

// CreateSong создает новую песню
func (r *PostgresSongRepository) CreateSong(ctx context.Context, song *models.Song) (*models.Song, error) {
	// Проверяем, существует ли уже песня с такими данными
	var exists bool
	err := r.db.QueryRowContext(ctx, checkSongExistsForCreateQuery,
		song.GroupName,
		song.SongName,
	).Scan(&exists)
	if err != nil {
		return nil, errors.NewInternal("failed to check song existence", err)
	}
	if exists {
		return nil, errors.NewAlreadyExists("song with this group name and song name already exists", nil)
	}

	err = r.db.QueryRowContext(ctx, createSongQuery,
		song.GroupName,
		song.SongName,
		song.ReleaseDate,
		song.Text,
		song.Link,
	).Scan(
		&song.ID,
		&song.GroupName,
		&song.SongName,
		&song.ReleaseDate,
		&song.Text,
		&song.Link,
		&song.CreatedAt,
		&song.UpdatedAt,
	)
	if err != nil {
		return nil, errors.NewInternal("failed to create song", err)
	}

	return song, nil
}

// UpdateSong обновляет информацию о песне
func (r *PostgresSongRepository) UpdateSong(ctx context.Context, song *models.Song) (*models.Song, error) {
	// Проверяем, существует ли уже песня с такими данными
	var exists bool
	err := r.db.QueryRowContext(ctx, checkSongExistsQuery,
		song.GroupName,
		song.SongName,
		song.ID,
	).Scan(&exists)
	if err != nil {
		return nil, errors.NewInternal("failed to check song existence", err)
	}
	if exists {
		return nil, errors.NewAlreadyExists("song with this group name and song name already exists", nil)
	}

	err = r.db.QueryRowContext(ctx, updateSongQuery,
		song.GroupName,
		song.SongName,
		song.ReleaseDate,
		song.Text,
		song.Link,
		song.ID,
	).Scan(
		&song.ID,
		&song.GroupName,
		&song.SongName,
		&song.ReleaseDate,
		&song.Text,
		&song.Link,
		&song.CreatedAt,
		&song.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.NewNotFound("song not found", err)
	}
	if err != nil {
		return nil, errors.NewInternal("failed to update song", err)
	}

	return song, nil
}

// DeleteSong удаляет песню
func (r *PostgresSongRepository) DeleteSong(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, deleteSongQuery, id)
	if err != nil {
		return errors.NewInternal("failed to delete song", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewInternal("failed to get rows affected", err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFound("song not found", nil)
	}

	return nil
}
