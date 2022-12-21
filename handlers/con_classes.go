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
type conClass struct {
	UpdatedOn  time.Time `json:"updated_on"`
	CreatedOn  time.Time `json:"created_on"`
	Notes      *string   `json:"notes"`
	Descr      string    `json:"descr"`
	ConClassID int64     `json:"con_class_id"`
}

// Import values from corresponding godevmandb struct
func (r *conClass) getValues(s godevmandb.ConClass) {
	r.ConClassID = s.ConClassID
	r.Descr = s.Descr
	r.UpdatedOn = s.UpdatedOn
	r.CreatedOn = s.CreatedOn
	r.Notes = nullStringToPtr(s.Notes)
}

// Return corresponding godevmandb create parameters
func (r *conClass) createParams() godevmandb.CreateConClassParams {
	s := godevmandb.CreateConClassParams{}

	s.Descr = r.Descr
	s.Notes = strToNullString(r.Notes)

	return s
}

// Return corresponding godevmandb update parameters
func (r *conClass) updateParams() godevmandb.UpdateConClassParams {
	s := godevmandb.UpdateConClassParams{}

	s.Descr = r.Descr
	s.Notes = strToNullString(r.Notes)

	return s
}

// Count ConClasses
// @Summary Count con_classes
// @Description Count number of connection classes
// @Tags connections
// @ID count-con_classes
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/classes/count [GET]
func (h *Handler) CountConClasses(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountConClasses(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List con_classes
// @Summary List con_classes
// @Description List connection classes info
// @Tags connections
// @ID list-con_classes
// @Param descr_f query string false "url encoded SQL 'LIKE' operator pattern"
// @Param limit query int false "min: 1; max: 1000; default: 1000"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} conClass
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/classes [GET]
func (h *Handler) GetConClasses(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetConClassesParams{
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
	res, err := q.GetConClasses(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := []conClass{}
	for _, s := range res {
		r := conClass{}
		r.getValues(s)
		out = append(out, r)
	}

	RespondJSON(w, r, http.StatusOK, out)
}

// Get ConClass
// @Summary Get con_class
// @Description Get connection class info
// @Tags connections
// @ID get-con_class
// @Param con_class_id path string true "con_class_id"
// @Success 200 {object} conClass
// @Failure 400 {object} StatusResponse "Invalid con_class_id"
// @Failure 404 {object} StatusResponse "Class not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/classes/{con_class_id} [GET]
func (h *Handler) GetConClass(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "con_class_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid class ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetConClass(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "Class not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	out := conClass{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Create ConClass
// @Summary Create con_class
// @Description Create connection class
// @Tags connections
// @ID create-con_class
// @Param Body body conClass true "JSON object of conClass<br />Ignored fields:<ul><li>con_class_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 201 {object} conClass
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/classes [POST]
func (h *Handler) CreateConClass(w http.ResponseWriter, r *http.Request) {
	var pIn conClass
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Create parameters for new db record
	p := pIn.createParams()

	q := godevmandb.New(h.db)
	res, err := q.CreateConClass(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := conClass{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusCreated, out)
}

// Update ConClass
// @Summary Update con_class
// @Description Update connection class
// @Tags connections
// @ID update-con_class
// @Param con_class_id path string true "con_class_id"
// @Param Body body conClass true "JSON object of conClass.<br />Ignored fields:<ul><li>con_class_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 200 {object} conClass
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/classes/{con_class_id} [PUT]
func (h *Handler) UpdateConClass(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "con_class_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid class ID")
		return
	}

	var pIn conClass
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Update parameters for new db record
	p := pIn.updateParams()
	p.ConClassID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateConClass(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := conClass{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Delete ConClass
// @Summary Delete con_class
// @Description Delete connection class
// @Tags connections
// @ID delete-con_class
// @Param con_class_id path string true "con_class_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid con_class_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/classes/{con_class_id} [DELETE]
func (h *Handler) DeleteConClass(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "con_class_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid class ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteConClass(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Relations
// List ConClass Connections
// @Summary List con_class connections
// @Description List connection class connections info
// @Tags connections
// @ID list-con_class-connections
// @Param con_class_id path string true "con_class_id"
// @Success 200 {array} connection
// @Failure 400 {object} StatusResponse "Invalid con_class_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/classes/{con_class_id}/connections [GET]
func (h *Handler) GetConClassConnections(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "con_class_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid class ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetConClassConnections(h.ctx, id)
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
