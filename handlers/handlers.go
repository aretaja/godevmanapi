package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/httplog"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Global encryption salt
var salt string

type Handler struct {
	ctx context.Context
	db  *pgxpool.Pool
}

// Create connection pool
func (h *Handler) Initialize(dbURL, s string) error {
	h.ctx = context.Background()
	salt = s

	pool, err := pgxpool.Connect(h.ctx, dbURL)
	if err != nil {
		return err
	}
	h.db = pool

	return nil
}

type StatusResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type CountResponse struct {
	Count int64 `json:"count"`
}

// Welcome
// @Summary Welcome
// @Description Welcome message
// @Tags information
// @ID root
// @Success 200 {object} StatusResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Router / [GET]
func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, r, http.StatusOK, StatusResponse{
		Code:    strconv.Itoa(http.StatusOK),
		Message: "GODEVMANAPI ready!",
	})
}

// Version - dummy function to generate swagger doc
// Endpoint is actually implemented in github.com/aretaja/godevmanapi/app.initializeRoutes()
// @Summary Version
// @Description Return API version info
// @Tags information
// @ID version
// @Success 200 {object} StatusResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Router /version [GET]
func VersionSwagger() {}

// Regular response
func RespondJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	res, err := json.Marshal(payload)
	if err != nil {
		log.Print(err.Error())
		RespondError(w, r, http.StatusInternalServerError, "JSON marshal failed")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}

// Error response
func RespondError(w http.ResponseWriter, r *http.Request, code int, message string) {
	hlog := httplog.LogEntry(r.Context())
	hlog.Error().Msg(message)
	res := StatusResponse{
		Code:    strconv.Itoa(code),
		Message: message,
	}
	RespondJSON(w, r, code, res)
}
