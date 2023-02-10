package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/aretaja/godevmandb"
	"github.com/go-chi/chi/v5"
)

// JSON friendly local type to use in web api. Replaces sql.Null*/pgtype fields
type conType struct {
	UpdatedOn time.Time `json:"updated_on"`
	CreatedOn time.Time `json:"created_on"`
	Notes     *string   `json:"notes"`
	Descr     string    `json:"descr"`
	ConTypeID int64     `json:"con_type_id"`
}

// Import values from corresponding godevmandb struct
func (r *conType) getValues(s godevmandb.ConType) {
	r.ConTypeID = s.ConTypeID
	r.Descr = s.Descr
	r.UpdatedOn = s.UpdatedOn
	r.CreatedOn = s.CreatedOn
	r.Notes = nullStringToPtr(s.Notes)
}

// Return corresponding godevmandb create parameters
func (r *conType) createParams() godevmandb.CreateConTypeParams {
	s := godevmandb.CreateConTypeParams{}

	s.Descr = r.Descr
	s.Notes = strToNullString(r.Notes)

	return s
}

// Return corresponding godevmandb update parameters
func (r *conType) updateParams() godevmandb.UpdateConTypeParams {
	s := godevmandb.UpdateConTypeParams{}

	s.Descr = r.Descr
	s.Notes = strToNullString(r.Notes)

	return s
}

// Count ConTypes
// @Summary Count con_types
// @Description Count number of connection types
// @Tags connections
// @ID count-con_types
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
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
// @Param descr_f query string false "url encoded SQL 'LIKE' operator pattern"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} conType
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
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

	// Descr filter
	d := r.FormValue("descr_f")
	if d != "" {
		p.DescrF = d
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetConTypes(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := []conType{}
	for _, s := range res {
		r := conType{}
		r.getValues(s)
		out = append(out, r)
	}

	RespondJSON(w, r, http.StatusOK, out)
}

// Get ConType
// @Summary Get con_type
// @Description Get connection type info
// @Tags connections
// @ID get-con_type
// @Param con_type_id path string true "con_type_id"
// @Success 200 {object} conType
// @Failure 400 {object} StatusResponse "Invalid con_type_id"
// @Failure 404 {object} StatusResponse "Type not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
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

	out := conType{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Create ConType
// @Summary Create con_type
// @Description Create connection type
// @Tags connections
// @ID create-con_type
// @Param Body body conType true "JSON object of conType.<br />Ignored fields:<ul><li>con_type_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 201 {object} conType
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/types [POST]
func (h *Handler) CreateConType(w http.ResponseWriter, r *http.Request) {
	var pIn conType
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Create parameters for new db record
	p := pIn.createParams()

	q := godevmandb.New(h.db)
	res, err := q.CreateConType(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := conType{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusCreated, out)
}

// Update ConType
// @Summary Update con_type
// @Description Update connection type
// @Tags connections
// @ID update-con_type
// @Param con_type_id path string true "con_type_id"
// @Param Body body conType true "JSON object of conType.<br />Ignored fields:<ul><li>con_type_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 200 {object} conType
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/types/{con_type_id} [PUT]
func (h *Handler) UpdateConType(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "con_type_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid type ID")
		return
	}

	var pIn conType
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Update parameters for new db record
	p := pIn.updateParams()
	p.ConTypeID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateConType(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := conType{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
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
// @Failure 500 {object} StatusResponse "Failde DB transaction"
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
// @Success 200 {array} connection
// @Failure 400 {object} StatusResponse "Invalid con_type_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
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

	out := []connection{}
	for _, s := range res {
		a := connection{}
		a.getValues(s)
		out = append(out, a)
	}

	RespondJSON(w, r, http.StatusOK, out)
}
