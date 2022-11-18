package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aretaja/godevmandb"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
)

// Count ConProviders
func (h *Handler) CountConProviders(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountConProviders(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, map[string]int64{"count": res})
}

// List ConProviders
func (h *Handler) GetConProviders(w http.ResponseWriter, r *http.Request) {
	hlog := httplog.LogEntry(r.Context())
	var p = godevmandb.GetConProvidersParams{
		Limit:  100,
		Offset: 0,
	}

	l, err := strconv.ParseInt(r.FormValue("count"), 10, 32)
	if err != nil {
		hlog.Info().Msg("Invalid count value. Using default")
	} else {
		if l < 100 || l > 0 {
			p.Limit = int32(l)
		}
	}
	o, err := strconv.ParseInt(r.FormValue("start"), 10, 32)
	if err != nil {
		hlog.Info().Msg("Invalid start value. Using default")
	} else {
		if o > 0 {
			p.Offset = int32(o)
		}
	}

	q := godevmandb.New(h.db)
	res, err := q.GetConProviders(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Get ConProvider
func (h *Handler) GetConProvider(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "con_prov_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid provider ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetConProvider(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "Provider not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Create ConProvider
func (h *Handler) CreateConProvider(w http.ResponseWriter, r *http.Request) {
	var p godevmandb.CreateConProviderParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	q := godevmandb.New(h.db)
	res, err := q.CreateConProvider(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusCreated, res)
}

// Update ConProvider
func (h *Handler) UpdateConProvider(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "con_prov_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid provider ID")
		return
	}

	var p godevmandb.UpdateConProviderParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	p.ConProvID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateConProvider(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Delete ConProvider
func (h *Handler) DeleteConProvider(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "con_prov_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid provider ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteConProvider(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
	w.WriteHeader(http.StatusNoContent)
}

// Handlers - Relations
func (h *Handler) GetConProviderConnections(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "con_prov_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid provider ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetConProviderConnections(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}
