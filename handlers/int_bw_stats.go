package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aretaja/godevmandb"
	"github.com/go-chi/chi/v5"
)

// Count Int Bw Stats
// @Summary Count int bw stats
// @Description Count number of int bw stats
// @Tags interfaces
// @ID count-bw_stats
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/bw_stats/count [GET]
func (h *Handler) CountIntBwStats(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountIntBwStats(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List int bw stats
// @Summary List int bw stats
// @Description List int bw stats info
// @Tags interfaces
// @ID list-bw_stats
// @Param if_group_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param to50in_le query int false "SQL '<=' operator value"
// @Param to50in_ge query int false "SQL '>=' operator value"
// @Param to75in_le query int false "SQL '<=' operator value"
// @Param to75in_ge query int false "SQL '>=' operator value"
// @Param to90in_le query int false "SQL '<=' operator value"
// @Param to90in_ge query int false "SQL '>=' operator value"
// @Param to100in_le query int false "SQL '<=' operator value"
// @Param to100in_ge query int false "SQL '>=' operator value"
// @Param to50out_le query int false "SQL '<=' operator value"
// @Param to50out_ge query int false "SQL '>=' operator value"
// @Param to75out_le query int false "SQL '<=' operator value"
// @Param to75out_ge query int false "SQL '>=' operator value"
// @Param to90out_le query int false "SQL '<=' operator value"
// @Param to90out_ge query int false "SQL '>=' operator value"
// @Param to100out_le query int false "SQL '<=' operator value"
// @Param to100out_ge query int false "SQL '>=' operator value"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} godevmandb.IntBwStat
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/bw_stats [GET]
func (h *Handler) GetIntBwStats(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetIntBwStatsParams{
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
	if v := r.FormValue("if_group_f"); v != "" {
		p.IfGroupF = v
	}

	if v := r.FormValue("to50in_le"); v != "" {
		p.To50inLe = &v
	}

	if v := r.FormValue("to50in_ge"); v != "" {
		p.To50inGe = &v
	}

	if v := r.FormValue("to75in_le"); v != "" {
		p.To75inLe = &v
	}

	if v := r.FormValue("to75in_ge"); v != "" {
		p.To75inGe = &v
	}

	if v := r.FormValue("to90in_le"); v != "" {
		p.To90inLe = &v
	}

	if v := r.FormValue("to90in_ge"); v != "" {
		p.To90inGe = &v
	}

	if v := r.FormValue("to100in_le"); v != "" {
		p.To100inLe = &v
	}

	if v := r.FormValue("to100in_ge"); v != "" {
		p.To100inGe = &v
	}

	if v := r.FormValue("to50out_le"); v != "" {
		p.To50outLe = &v
	}

	if v := r.FormValue("to50out_ge"); v != "" {
		p.To50outGe = &v
	}

	if v := r.FormValue("to75out_le"); v != "" {
		p.To75outLe = &v
	}

	if v := r.FormValue("to75out_ge"); v != "" {
		p.To75outGe = &v
	}

	if v := r.FormValue("to90out_le"); v != "" {
		p.To90outLe = &v
	}

	if v := r.FormValue("to90out_ge"); v != "" {
		p.To90outGe = &v
	}

	if v := r.FormValue("to100out_le"); v != "" {
		p.To100inLe = &v
	}

	if v := r.FormValue("to100out_ge"); v != "" {
		p.To100inGe = &v
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetIntBwStats(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Get IntBwStat
// @Summary Get IntBwStat
// @Description Get IntBwStat info
// @Tags interfaces
// @ID get-int_bw_stat
// @Param bw_id path string true "bw_id"
// @Success 200 {object} godevmandb.IntBwStat
// @Failure 400 {object} StatusResponse "Invalid bw_id"
// @Failure 404 {object} StatusResponse "Stat not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/bw_stats/{bw_id} [GET]
func (h *Handler) GetIntBwStat(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "bw_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid Stat ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetIntBwStat(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "Stat not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Create IntBwStat
// @Summary Create IntBwStat
// @Description Create IntBwStat
// @Tags interfaces
// @ID create-int_bw_stat
// @Param Body body godevmandb.CreateIntBwStatParams true "JSON object of godevmandb.CreateIntBwStatParams"
// @Success 201 {object} godevmandb.IntBwStat
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/bw_stats [POST]
func (h *Handler) CreateIntBwStat(w http.ResponseWriter, r *http.Request) {
	var p godevmandb.CreateIntBwStatParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	q := godevmandb.New(h.db)
	res, err := q.CreateIntBwStat(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusCreated, res)
}

// Update IntBwStat
// @Summary Update IntBwStat
// @Description Update IntBwStat
// @Tags interfaces
// @ID update-int_bw_stat
// @Param bw_id path string true "bw_id"
// @Param Body body godevmandb.UpdateIntBwStatParams true "JSON object of godevmandb.UpdateIntBwStatParams.<br />Ignored fields:<ul><li>bw_id</li></ul>"
// @Success 200 {object} godevmandb.IntBwStat
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/bw_stats/{bw_id} [PUT]
func (h *Handler) UpdateIntBwStat(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "bw_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid Stat ID")
		return
	}

	var p godevmandb.UpdateIntBwStatParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	p.BwID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateIntBwStat(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Delete IntBwStat
// @Summary Delete IntBwStat
// @Description Delete IntBwStat
// @Tags interfaces
// @ID delete-int_bw_stat
// @Param bw_id path string true "bw_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid bw_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/bw_stats/{bw_id} [DELETE]
func (h *Handler) DeleteIntBwStat(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "bw_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid Stat ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteIntBwStat(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Foreign key
// Get Int Bw Stat Interface
// @Summary Get bw_stats interface
// @Description Get bw_stats interface info
// @Tags interfaces
// @ID get-bw_stats-interface
// @Param bw_id path string true "bw_id"
// @Success 200 {object} iface
// @Failure 400 {object} StatusResponse "Invalid bw_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/bw_stats/{bw_id}/interface [GET]
func (h *Handler) GetIntBwStatInterface(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "bw_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid stat ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetIntBwStatInterface(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := iface{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}
