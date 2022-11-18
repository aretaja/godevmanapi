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

type Handler struct {
	db  *pgxpool.Pool
	ctx context.Context
}

// Create connection pool
func (h *Handler) Initialize(dbURL string) error {
	h.ctx = context.Background()

	pool, err := pgxpool.Connect(h.ctx, dbURL)
	if err != nil {
		return err
	}
	h.db = pool

	return nil
}

// Root
func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, r, http.StatusOK, map[string]string{"message": "GODEVMANAPI ready!"})
}

// Helpers
func RespondError(w http.ResponseWriter, r *http.Request, code int, message string) {
	hlog := httplog.LogEntry(r.Context())
	hlog.Error().Msg(message)
	res := map[string]string{
		"error":   strconv.Itoa(code),
		"message": message,
	}
	RespondJSON(w, r, code, res)
}

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
