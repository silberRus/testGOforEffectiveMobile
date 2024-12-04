package service

import (
	"context"
	"github.com/testTask/internal/errors"
	"strings"
	"time"

	"github.com/testTask/internal/models"
	"github.com/testTask/internal/repository"
	"go.uber.org/zap"
)

type SongService struct {
	repo   repository.SongRepository
	logger *zap.Logger
}

func NewSongService(repo repository.SongRepository, logger *zap.Logger) *SongService {
	return &SongService{
		repo:   repo,
		logger: logger,
	}
}

// GetSongs получает список песен с опциональным фильтром
func (s *SongService) GetSongs(ctx context.Context, filter *models.SongFilter) (*models.SongsResponse, error) {
	s.logger.Info("Getting songs with filter",
		zap.String("group", filter.GroupName),
		zap.String("song", filter.SongName),
		zap.Any("fromDate", filter.FromDate),
		zap.Any("toDate", filter.ToDate),
		zap.String("text", filter.Text),
		zap.String("link", filter.Link),
		zap.Int("page", filter.Page),
		zap.Int("pageSize", filter.PageSize))

	return s.repo.GetSongs(ctx, filter)
}

// GetLyrics получает текст песни с опциональным фильтром
func (s *SongService) GetLyrics(ctx context.Context, id, page, pageSize int) (*models.LyricsResponse, error) {
	s.logger.Info("Getting lyrics",
		zap.Int("songId", id),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	song, err := s.repo.GetSongByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if song.Text == "" {
		return nil, errors.NewLyricsNotFound("lyrics not found", nil)
	}

	verses := strings.Split(song.Text, "\n\n")
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	totalPages := (len(verses) + pageSize - 1) / pageSize
	if page > totalPages {
		return nil, errors.NewNotFound("page out of range", nil)
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	if end > len(verses) {
		end = len(verses)
	}

	return &models.LyricsResponse{
		Text:        strings.Join(verses[start:end], "\n\n"),
		CurrentPage: page,
		TotalPages:  totalPages,
		PageSize:    pageSize,
	}, nil
}

// CreateSong создает новую песню
func (s *SongService) CreateSong(ctx context.Context, req *models.SongRequest) (*models.Song, error) {
	s.logger.Info("Creating new song",
		zap.String("group", req.GroupName),
		zap.String("song", req.SongName))

	// Создаем песню с предоставленными данными
	song := &models.Song{
		GroupName:   req.GroupName,
		SongName:    req.SongName,
		ReleaseDate: time.Now(),
		Text:        req.Text,
		Link:        req.Link,
	}

	return s.repo.CreateSong(ctx, song)
}

// UpdateSong обновляет существующую песню
func (s *SongService) UpdateSong(ctx context.Context, id int, req *models.SongRequest) (*models.Song, error) {
	s.logger.Info("Updating song",
		zap.Int("id", id),
		zap.String("group", req.GroupName),
		zap.String("song", req.SongName))

	song, err := s.repo.GetSongByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Обновляем только предоставленные поля
	if req.GroupName != "" {
		song.GroupName = req.GroupName
	}
	if req.SongName != "" {
		song.SongName = req.SongName
	}
	if req.Text != "" {
		song.Text = req.Text
	}
	if req.Link != "" {
		song.Link = req.Link
	}
	song.UpdatedAt = time.Now()

	return s.repo.UpdateSong(ctx, song)
}

// DeleteSong удаляет существующую песню
func (s *SongService) DeleteSong(ctx context.Context, id int) error {
	s.logger.Info("Deleting song", zap.Int("id", id))
	return s.repo.DeleteSong(ctx, id)
}
