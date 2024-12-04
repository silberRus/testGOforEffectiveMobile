package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/testTask/internal/errors"
	"github.com/testTask/internal/models"
	"github.com/testTask/internal/service"
	"go.uber.org/zap"
)

type SongHandler struct {
	service *service.SongService
	logger  *zap.Logger
}

func NewSongHandler(service *service.SongService, logger *zap.Logger) *SongHandler {
	return &SongHandler{
		service: service,
		logger:  logger,
	}
}

// handleError обрабатывает ошибки и возвращает соответствующий HTTP статус
func (h *SongHandler) handleError(w http.ResponseWriter, err error) {
	var status int
	var message string

	if appErr, ok := err.(*errors.Error); ok {
		switch appErr.Type {
		case errors.NotFound:
			status = http.StatusNotFound
			message = appErr.Message
		case errors.BadRequest:
			status = http.StatusBadRequest
			message = appErr.Message
		case errors.Validation:
			status = http.StatusUnprocessableEntity
			message = appErr.Message
		case errors.AlreadyExists:
			status = http.StatusConflict
			message = appErr.Message
		default:
			status = http.StatusInternalServerError
			message = "Internal server error"
		}
	} else {
		status = http.StatusInternalServerError
		message = "Internal server error"
	}

	h.logger.Error("Request error",
		zap.Error(err),
		zap.Int("status", status),
		zap.String("message", message),
	)

	http.Error(w, message, status)
}

// @Summary Get songs with filtering and pagination
// @Description Get list of songs with optional filtering and pagination
// @Tags songs
// @Accept json
// @Produce json
// @Param group_name query string false "Group name"
// @Param song_name query string false "Song name"
// @Param from_date query string false "From date (format: 2006-01-02)"
// @Param to_date query string false "To date (format: 2006-01-02)"
// @Param text query string false "Text content"
// @Param link query string false "Link"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} models.SongsResponse
// @Router /songs [get]
func (h *SongHandler) GetSongs(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling GetSongs request")

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	filter := &models.SongFilter{
		GroupName: r.URL.Query().Get("group_name"),
		SongName:  r.URL.Query().Get("song_name"),
		Text:      r.URL.Query().Get("text"),
		Link:      r.URL.Query().Get("link"),
		Page:      page,
		PageSize:  pageSize,
	}

	// Парсим даты, если они предоставлены
	if fromDateStr := r.URL.Query().Get("from_date"); fromDateStr != "" {
		if fromDate, err := time.Parse("2006-01-02", fromDateStr); err == nil {
			filter.FromDate = &fromDate
		}
	}
	if toDateStr := r.URL.Query().Get("to_date"); toDateStr != "" {
		if toDate, err := time.Parse("2006-01-02", toDateStr); err == nil {
			filter.ToDate = &toDate
		}
	}

	response, err := h.service.GetSongs(r.Context(), filter)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.handleError(w, errors.NewValidation("json encode error", err))
		return
	}
}

// @Summary Get song lyrics
// @Description Get song lyrics with pagination by verses
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} models.LyricsResponse
// @Router /songs/{id}/lyrics [get]
func (h *SongHandler) GetLyrics(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling GetLyrics request")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.handleError(w, errors.NewBadRequest("Invalid song ID", err))
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	response, err := h.service.GetLyrics(r.Context(), id, page, pageSize)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.handleError(w, errors.NewValidation("json encode error", err))
		return
	}
}

// @Summary Create new song
// @Description Create a new song with information from external API
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.SongRequest true "Song information"
// @Success 201 {object} models.Song
// @Router /songs [post]
func (h *SongHandler) CreateSong(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling CreateSong request")

	var req models.SongRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, errors.NewBadRequest("Invalid request body", err))
		return
	}

	song, err := h.service.CreateSong(r.Context(), &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(song)
	if err != nil {
		h.handleError(w, errors.NewValidation("json encode error", err))
		return
	}
}

// @Summary Update song
// @Description Update existing song information
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body models.SongRequest true "Song information"
// @Success 200 {object} models.Song
// @Router /songs/{id} [put]
func (h *SongHandler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling UpdateSong request")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.handleError(w, errors.NewBadRequest("Invalid song ID", err))
		return
	}

	var req models.SongRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, errors.NewBadRequest("Invalid request body", err))
		return
	}

	song, err := h.service.UpdateSong(r.Context(), id, &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(song)
	if err != nil {
		h.handleError(w, errors.NewValidation("json encode error", err))
		return
	}
}

// @Summary Delete song
// @Description Delete a song by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 204 "No Content"
// @Router /songs/{id} [delete]
func (h *SongHandler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling DeleteSong request")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.handleError(w, errors.NewBadRequest("Invalid song ID", err))
		return
	}

	if err := h.service.DeleteSong(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
