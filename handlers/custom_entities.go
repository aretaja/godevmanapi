package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aretaja/godevmandb"
	"github.com/go-chi/chi/v5"
)

// Count CustomEntities
// @Summary Count custom_entities
// @Description Count number of custom_entities
// @Tags data
// @ID count-custom_entities
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /data/custom_entities/count [GET]
func (h *Handler) CountCustomEntities(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountCustomEntities(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List custom_entities
// @Summary List custom_entities
// @Description List custom_entities info
// @Tags data
// @ID list-custom_entities
// @Param serial_nr_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} godevmandb.CustomEntity
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /data/custom_entities [GET]
func (h *Handler) GetCustomEntities(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetCustomEntitiesParams{
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

	// Serial nr filter
	d := r.FormValue("serial_nr_f")
	if d != "" {
		p.SerialNrF = d
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetCustomEntities(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Get CustomEntity
// @Summary Get customEntity
// @Description Get customEntity info
// @Tags data
// @ID get-customEntity
// @Param cent_id path string true "cent_id"
// @Success 200 {object} godevmandb.CustomEntity
// @Failure 400 {object} StatusResponse "Invalid cent_id"
// @Failure 404 {object} StatusResponse "CustomEntity not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /data/custom_entities/{cent_id} [GET]
func (h *Handler) GetCustomEntity(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "cent_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid customEntity ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetCustomEntity(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "CustomEntity not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Create CustomEntity
// @Summary Create customEntity
// @Description Create customEntity
// @Tags data
// @ID create-customEntity
// @Param Body body godevmandb.CreateCustomEntityParams true "JSON object of godevmandb.CreateCustomEntityParams"
// @Success 201 {object} godevmandb.CustomEntity
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /data/custom_entities [POST]
func (h *Handler) CreateCustomEntity(w http.ResponseWriter, r *http.Request) {
	var p godevmandb.CreateCustomEntityParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	q := godevmandb.New(h.db)
	res, err := q.CreateCustomEntity(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusCreated, res)
}

// Update CustomEntity
// @Summary Update customEntity
// @Description Update customEntity
// @Tags data
// @ID update-customEntity
// @Param cent_id path string true "cent_id"
// @Param Body body godevmandb.UpdateCustomEntityParams true "JSON object of godevmandb.UpdateCustomEntityParams.<br />Ignored fields:<ul><li>cent_id</li></ul>"
// @Success 200 {object} godevmandb.CustomEntity
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /data/custom_entities/{cent_id} [PUT]
func (h *Handler) UpdateCustomEntity(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "cent_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid customEntity ID")
		return
	}

	var p godevmandb.UpdateCustomEntityParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	p.CentID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateCustomEntity(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Delete CustomEntity
// @Summary Delete customEntity
// @Description Delete customEntity
// @Tags data
// @ID delete-customEntity
// @Param cent_id path string true "cent_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid cent_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /data/custom_entities/{cent_id} [DELETE]
func (h *Handler) DeleteCustomEntity(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "cent_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid customEntity ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteCustomEntity(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
