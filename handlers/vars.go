package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aretaja/godevmandb"
	"github.com/go-chi/chi/v5"
)

// Count Vars
// @Summary Count vars
// @Description Count number of vars
// @Tags config
// @ID count-vars
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /config/vars/count [GET]
func (h *Handler) CountVars(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountVars(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List vars
// @Summary List vars
// @Description List vars info
// @Tags config
// @ID list-vars
// @Param descr_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param content_f query string false "url encoded SQL 'ILIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param notes_f query string false "url encoded SQL 'ILIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} godevmandb.Var
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /config/vars [GET]
func (h *Handler) GetVars(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetVarsParams{
		LimitQ:  100,
		OffsetQ: 0,
	}

	lp := paginateValues(r)
	if lp[0] != nil {
		if *lp[0] < 100 {
			p.LimitQ = *lp[0]
		}
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

	// Filters
	if v := r.FormValue("descr_f"); v != "" {
		p.DescrF = v
	}

	if v := r.FormValue("content_f"); v != "" {
		p.ContentF = &v
	}

	if v := r.FormValue("notes_f"); v != "" {
		p.NotesF = &v
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetVars(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Get Var
// @Summary Get var
// @Description Get var info
// @Tags config
// @ID get-var
// @Param descr path string true "descr"
// @Success 200 {object} godevmandb.Var
// @Failure 400 {object} StatusResponse "Invalid descr"
// @Failure 404 {object} StatusResponse "Var not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /config/vars/{descr} [GET]
func (h *Handler) GetVar(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "descr")

	q := godevmandb.New(h.db)
	res, err := q.GetVar(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "Var not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Create Var
// @Summary Create var
// @Description Create var
// @Tags config
// @ID create-var
// @Param Body body godevmandb.CreateVarParams true "JSON object of godevmandb.CreateVarParams"
// @Success 201 {object} godevmandb.Var
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /config/vars [POST]
func (h *Handler) CreateVar(w http.ResponseWriter, r *http.Request) {
	var p godevmandb.CreateVarParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	q := godevmandb.New(h.db)
	res, err := q.CreateVar(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusCreated, res)
}

// Update Var
// @Summary Update var
// @Description Update var
// @Tags config
// @ID update-var
// @Param descr path string true "descr"
// @Param Body body godevmandb.UpdateVarParams true "JSON object of godevmandb.UpdateVarParams.<br />Ignored fields:<ul><li>descr</li></ul>"
// @Success 200 {object} godevmandb.Var
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /config/vars/{descr} [PUT]
func (h *Handler) UpdateVar(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "descr")

	var p godevmandb.UpdateVarParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	p.Descr = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateVar(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Delete Var
// @Summary Delete var
// @Description Delete var
// @Tags config
// @ID delete-var
// @Param descr path string true "descr"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid descr"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /config/vars/{descr} [DELETE]
func (h *Handler) DeleteVar(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "descr")

	q := godevmandb.New(h.db)
	err := q.DeleteVar(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
