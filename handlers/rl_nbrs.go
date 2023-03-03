package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aretaja/godevmandb"
	"github.com/go-chi/chi/v5"
)

// Count RlNbrs
// @Summary Count rl_nbrs
// @Description Count number of radio link neighbors
// @Tags devices
// @ID count-rl_nbrs
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/rl_nbrs/count [GET]
func (h *Handler) CountRlNbrs(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountRlNbrs(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List rl_nbrs
// @Summary List rl_nbrs
// @Description List radio link neighbors info
// @Tags devices
// @ID list-rl_nbrs
// @Param dev_id_f query string false "url encoded SQL 'LIKE' operator pattern"
// @Param nbr_ent_id_f query string false "url encoded SQL 'LIKE' operator pattern"
// @Param nbr_sysname_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} godevmandb.RlNbr
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/rl_nbrs [GET]
func (h *Handler) GetRlNbrs(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetRlNbrsParams{
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
	if v := r.FormValue("dev_id_f"); v != "" {
		p.DevIDF = v
	}

	if v := r.FormValue("nbr_ent_id_f"); v != "" {
		p.NbrEntIDF = &v
	}
	if v := r.FormValue("nbr_sysname_f"); v != "" {
		p.NbrSysnameF = v
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetRlNbrs(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Get RlNbr
// @Summary Get rl_nbr
// @Description Get radio link neighbor info
// @Tags devices
// @ID get-rl_nbr
// @Param nbr_id path string true "nbr_id"
// @Success 200 {object} godevmandb.RlNbr
// @Failure 400 {object} StatusResponse "Invalid nbr_id"
// @Failure 404 {object} StatusResponse "Domain not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/rl_nbrs/{nbr_id} [GET]
func (h *Handler) GetRlNbr(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "nbr_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid neighbor ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetRlNbr(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "Neighbor not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Create RlNbr
// @Summary Create rl_nbr
// @Description Create radio link neighbor
// @Tags devices
// @ID create-rl_nbr
// @Param Body body godevmandb.CreateRlNbrParams true "JSON object of godevmandb.RlNbrParams"
// @Success 201 {object} godevmandb.RlNbr
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/rl_nbrs [POST]
func (h *Handler) CreateRlNbr(w http.ResponseWriter, r *http.Request) {
	var p godevmandb.CreateRlNbrParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	q := godevmandb.New(h.db)
	res, err := q.CreateRlNbr(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusCreated, res)
}

// Update RlNbr
// @Summary Update rl_nbr
// @Description Update radio link neighbor
// @Tags devices
// @ID update-rl_nbr
// @Param nbr_id path string true "nbr_id"
// @Param Body body godevmandb.UpdateRlNbrParams true "JSON object of godevmandb.UpdateRlNbrParams.<br />Ignored fields:<ul><li>nbr_id</li></ul>"
// @Success 200 {object} godevmandb.RlNbr
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/rl_nbrs/{nbr_id} [PUT]
func (h *Handler) UpdateRlNbr(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "nbr_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid neighbor ID")
		return
	}

	var p godevmandb.UpdateRlNbrParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	p.NbrID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateRlNbr(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Delete RlNbr
// @Summary Delete rl_nbr
// @Description Delete radio link neighbor
// @Tags devices
// @ID delete-rl_nbr
// @Param nbr_id path string true "nbr_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid nbr_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/rl_nbrs/{nbr_id} [DELETE]
func (h *Handler) DeleteRlNbr(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "nbr_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid neighbor ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteRlNbr(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Foreign key
// Get RlNbr Device
// @Summary Get rl_nbr device
// @Description Get rl_nbr device info
// @Tags devices
// @ID get-rl_nbr-device
// @Param nbr_id path string true "nbr_id"
// @Success 200 {object} device
// @Failure 400 {object} StatusResponse "Invalid nbr_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/rl_nbrs/{nbr_id}/device [GET]
func (h *Handler) GetRlNbrDevice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "nbr_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid rl_nbr ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetRlNbrDevice(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := device{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Foreign key
// Get RlNbr Entity
// @Summary Get rl_nbr entity
// @Description Get rl_nbr entity info
// @Tags devices
// @ID get-rl_nbr-entity
// @Param nbr_id path string true "nbr_id"
// @Success 200 {object} godevmandb.Entity
// @Failure 400 {object} StatusResponse "Invalid nbr_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/rl_nbrs/{nbr_id}/entity [GET]
func (h *Handler) GetRlNbrEntity(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "nbr_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid rl_nbr ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetRlNbrEntity(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}
