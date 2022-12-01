package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aretaja/godevmandb"
	"github.com/go-chi/chi/v5"
)

// Count ConCapacities
// @Summary Count con_capacities
// @Description Count number of connection capacities
// @Tags connections
// @ID count-con_capacities
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/capacities/count [GET]
func (h *Handler) CountConCapacities(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountConCapacities(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List con_capacities
// @Summary List con_capacities
// @Description List connection capacities info
// @Tags connections
// @ID list-con_capacities
// @Param descr_f query string false "url encoded SQL like value"
// @Param limit query int false "min: 1; max: 1000; default: 1000"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} godevmandb.ConCapacity
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/capacities [GET]
func (h *Handler) GetConCapacities(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetConCapacitiesParams{
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
	res, err := q.GetConCapacities(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Get ConCapacity
// @Summary Get capacity
// @Description Get connection capacity info
// @Tags connections
// @ID get-capacity
// @Param con_cap_id path string true "con_cap_id"
// @Success 200 {object} godevmandb.ConCapacity
// @Failure 400 {object} StatusResponse "Invalid con_cap_id"
// @Failure 404 {object} StatusResponse "Provider not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/capacities/{con_cap_id} [GET]
func (h *Handler) GetConCapacity(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "con_cap_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid capacity ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetConCapacity(h.ctx, id)
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

// Create ConCapacity
// @Summary Create capacity
// @Description Create connection capacity
// @Tags connections
// @ID create-capacity
// @Param Body body godevmandb.CreateConCapacityParams true "JSON object of CreateConCapacityParams"
// @Success 201 {object} godevmandb.ConCapacity
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/capacities [POST]
func (h *Handler) CreateConCapacity(w http.ResponseWriter, r *http.Request) {
	var p godevmandb.CreateConCapacityParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	q := godevmandb.New(h.db)
	res, err := q.CreateConCapacity(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusCreated, res)
}

// Update ConCapacity
// @Summary Update capacity
// @Description Update connection capacity
// @Tags connections
// @ID update-capacity
// @Param con_cap_id path string true "con_cap_id"
// @Param Body body godevmandb.UpdateConCapacityParams true "JSON object of UpdateConCapacityParams"
// @Success 200 {object} godevmandb.ConCapacity
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/capacities/{con_cap_id} [PUT]
func (h *Handler) UpdateConCapacity(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "con_cap_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid capacity ID")
		return
	}

	var p godevmandb.UpdateConCapacityParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	p.ConCapID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateConCapacity(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Delete ConCapacity
// @Summary Delete capacity
// @Description Delete connection capacity
// @Tags connections
// @ID delete-capacity
// @Param con_cap_id path string true "con_cap_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid con_cap_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/capacities/{con_cap_id} [DELETE]
func (h *Handler) DeleteConCapacity(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "con_cap_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid capacity ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteConCapacity(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
	w.WriteHeader(http.StatusNoContent)
}

// Relations
// List ConCapacity Connections
// @Summary List capacity connections
// @Description List connection capacity connections info
// @Tags connections
// @ID list-capacity-connections
// @Param con_cap_id path string true "con_cap_id"
// @Success 200 {array} godevmandb.Connection
// @Failure 400 {object} StatusResponse "Invalid con_cap_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/capacities/{con_cap_id}/connections [GET]
func (h *Handler) GetConCapacityConnections(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "con_cap_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid capacity ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetConCapacityConnections(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}
