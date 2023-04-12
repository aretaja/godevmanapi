package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aretaja/godevmandb"
	"github.com/go-chi/chi/v5"
)

// Count ConTypes
// @Summary Count con_types
// @Description Count number of connection types
// @Tags connections
// @ID count-con_types
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /connections/types/count [GET]
func (h *Handler) CountConTypes(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountConTypes(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List con_types
// @Summary List con_types
// @Description List connection types info
// @Tags connections
// @ID list-con_types
// @Param descr_f query string false "url encoded SQL 'ILIKE' operator pattern + special value 'isempty'"
// @Param notes_f query string false "url encoded SQL 'ILIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} godevmandb.ConType
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /connections/types [GET]
func (h *Handler) GetConTypes(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetConTypesParams{
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

	if v := r.FormValue("notes_f"); v != "" {
		p.NotesF = &v
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetConTypes(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Get ConType
// @Summary Get con_type
// @Description Get connection type info
// @Tags connections
// @ID get-con_type
// @Param con_type_id path string true "con_type_id"
// @Success 200 {object} godevmandb.ConType
// @Failure 400 {object} StatusResponse "Invalid con_type_id"
// @Failure 404 {object} StatusResponse "Type not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /connections/types/{con_type_id} [GET]
func (h *Handler) GetConType(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "con_type_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid type ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetConType(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "Type not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Create ConType
// @Summary Create con_type
// @Description Create connection type
// @Tags connections
// @ID create-con_type
// @Param Body body godevmandb.CreateConTypeParams true "JSON object of godevmandb.CreateConTypeParams"
// @Success 201 {object} godevmandb.ConType
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /connections/types [POST]
func (h *Handler) CreateConType(w http.ResponseWriter, r *http.Request) {
	var p godevmandb.CreateConTypeParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	q := godevmandb.New(h.db)
	res, err := q.CreateConType(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusCreated, res)
}

// Update ConType
// @Summary Update con_type
// @Description Update connection type
// @Tags connections
// @ID update-con_type
// @Param con_type_id path string true "con_type_id"
// @Param Body body godevmandb.UpdateConTypeParams true "JSON object of godevmandb.UpdateConTypeParams.<br />Ignored fields:<ul><li>con_type_id</li></ul>"
// @Success 200 {object} godevmandb.ConType
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /connections/types/{con_type_id} [PUT]
func (h *Handler) UpdateConType(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "con_type_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid type ID")
		return
	}

	var p godevmandb.UpdateConTypeParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	p.ConTypeID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateConType(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Delete ConType
// @Summary Delete con_type
// @Description Delete connection type
// @Tags connections
// @ID delete-con_type
// @Param con_type_id path string true "con_type_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid con_type_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /connections/types/{con_type_id} [DELETE]
func (h *Handler) DeleteConType(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "con_type_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid type ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteConType(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Relations
// List ConType Connections
// @Summary List con_type connections
// @Description List connection type connections info
// @Tags connections
// @ID list-con_type-connections
// @Param con_type_id path string true "con_type_id"
// @Success 200 {array} godevmandb.Connection
// @Failure 400 {object} StatusResponse "Invalid con_type_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /connections/types/{con_type_id}/connections [GET]
func (h *Handler) GetConTypeConnections(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "con_type_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid type ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetConTypeConnections(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}
