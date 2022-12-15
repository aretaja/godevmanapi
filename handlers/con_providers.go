package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/aretaja/godevmandb"
	"github.com/go-chi/chi/v5"
	"github.com/jinzhu/copier"
)

// Local type to use with copier. Used for sql Null* field replacement in json response
type ConProvider struct {
	ConProvID int64     `json:"con_prov_id"`
	Descr     string    `json:"descr"`
	NotesVal  *string   `json:"notes"`
	UpdatedOn time.Time `json:"updated_on"`
	CreatedOn time.Time `json:"created_on"`
}

func (a *ConProvider) Notes(m sql.NullString) {
	if m.Valid {
		if v, err := m.Value(); err == nil {
			if res, ok := v.(string); ok {
				a.NotesVal = &res
			}
		}
	}
}

// Count ConProviders
// @Summary Count con_providers
// @Description Count number of connection providers
// @Tags connections
// @ID count-con_providers
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/providers/count [GET]
func (h *Handler) CountConProviders(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountConProviders(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List con_providers
// @Summary List con_providers
// @Description List connection providers info
// @Tags connections
// @ID list-con_providers
// @Param descr_f query string false "url encoded SQL 'LIKE' operator pattern"
// @Param limit query int false "min: 1; max: 1000; default: 1000"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} ConProvider
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/providers [GET]
func (h *Handler) GetConProviders(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetConProvidersParams{
		LimitQ:  100,
		OffsetQ: 0,
	}

	lp := paginateValues(r)
	if lp[0] != nil {
		p.LimitQ = *lp[0]
	}
	if lp[1] != nil {
		p.OffsetQ = *lp[1]
	}

	// Time filter
	tf := parseTimeFilter(r)
	p.UpdatedGe = tf[0]
	p.UpdatedLe = tf[1]
	p.CreatedGe = tf[2]
	p.CreatedLe = tf[3]

	// Descr filter
	d := r.FormValue("descr_f")
	if d != "" {
		p.DescrF = d
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetConProviders(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := []ConProvider{}
	copier.Copy(&out, &res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Get ConProvider
// @Summary Get con_provider
// @Description Get connection provider info
// @Tags connections
// @ID get-con_provider
// @Param con_prov_id path string true "con_prov_id"
// @Success 200 {object} ConProvider
// @Failure 400 {object} StatusResponse "Invalid con_prov_id"
// @Failure 404 {object} StatusResponse "Provider not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/providers/{con_prov_id} [GET]
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

	out := ConProvider{}
	copier.Copy(&out, &res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Create ConProvider
// @Summary Create con_provider
// @Description Create connection provider
// @Tags connections
// @ID create-con_provider
// @Param Body body godevmandb.CreateConProviderParams true "JSON object of CreateConProviderParams"
// @Success 201 {object} ConProvider
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/providers [POST]
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

	out := ConProvider{}
	copier.Copy(&out, &res)

	RespondJSON(w, r, http.StatusCreated, out)
}

// Update ConProvider
// @Summary Update con_provider
// @Description Update connection provider
// @Tags connections
// @ID update-con_provider
// @Param con_prov_id path string true "con_prov_id"
// @Param Body body godevmandb.UpdateConProviderParams true "JSON object of UpdateConProviderParams"
// @Success 200 {object} ConProvider
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/providers/{con_prov_id} [PUT]
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

	out := ConProvider{}
	copier.Copy(&out, &res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Delete ConProvider
// @Summary Delete con_provider
// @Description Delete connection provider
// @Tags connections
// @ID delete-con_provider
// @Param con_prov_id path string true "con_prov_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid con_prov_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/providers/{con_prov_id} [DELETE]
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

	w.WriteHeader(http.StatusNoContent)
}

// Relations
// List ConProvider Connections
// @Summary List con_provider connections
// @Description List connection provider connections info
// @Tags connections
// @ID list-con_provider-connections
// @Param con_prov_id path string true "con_prov_id"
// @Success 200 {array} Connection
// @Failure 400 {object} StatusResponse "Invalid con_prov_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/providers/{con_prov_id}/connections [GET]
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

	out := []Connection{}
	copier.Copy(&out, &res)

	RespondJSON(w, r, http.StatusOK, out)
}
