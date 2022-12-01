package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

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

// Helpers - Regular response
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

// Helpers - Error response
func RespondError(w http.ResponseWriter, r *http.Request, code int, message string) {
	hlog := httplog.LogEntry(r.Context())
	hlog.Error().Msg(message)
	res := StatusResponse{
		Code:    strconv.Itoa(code),
		Message: message,
	}
	RespondJSON(w, r, code, res)
}

// Helpers - paginateValues
func paginateValues(r *http.Request) []*int32 {
	res := make([]*int32, 2)
	hlog := httplog.LogEntry(r.Context())

	l, err := strconv.ParseInt(r.FormValue("limit"), 10, 32)
	if err != nil {
		hlog.Debug().Msg("Parse limit - " + err.Error())
		hlog.Info().Msg("Invalid limit value. Using default")
	} else {
		if l < 1000 || l > 0 {
			lo := int32(l)
			res[0] = &lo
		} else {
			hlog.Info().Msg("Value of limit value out of range 0 - 1000")
		}
	}

	o, err := strconv.ParseInt(r.FormValue("offset"), 10, 32)
	if err != nil {
		hlog.Debug().Msg("Parse offset - " + err.Error())
		hlog.Info().Msg("Invalid offset value. Using default")
	} else {
		if o > 0 {
			oo := int32(o)
			res[1] = &oo
		} else {
			hlog.Info().Msg("Value of offset value out of range > 0")
		}
	}

	return res
}

// Helpers - Timefilter
func parseTimeFilter(r *http.Request) []time.Time {
	res := make([]time.Time, 4)
	hlog := httplog.LogEntry(r.Context())
	keys := []string{"updated_ge", "updated_le", "created_ge", "created_le"}

	for i := 0; i < 4; i++ {
		v := r.FormValue(keys[i])
		if v == "" {
			continue
		}
		uts, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			hlog.Debug().Msg("Parse " + keys[i] + " - " + err.Error())
		} else {
			res[i] = time.UnixMilli(uts)
		}
	}

	return res
}
