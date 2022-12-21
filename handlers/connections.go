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
type connection struct {
	UpdatedOn  time.Time `json:"updated_on"`
	CreatedOn  time.Time `json:"created_on"`
	Hint       *string   `json:"hint"`
	Notes      *string   `json:"notes"`
	ConID      int64     `json:"con_id"`
	SiteID     int64     `json:"site_id"`
	ConProvID  int64     `json:"con_prov_id"`
	ConTypeID  int64     `json:"con_type_id"`
	ConCapID   int64     `json:"con_cap_id"`
	ConClassID int64     `json:"con_class_id"`
	InUse      bool      `json:"in_use"`
}

// Import values from corresponding godevmandb struct
func (r *connection) getValues(s godevmandb.Connection) {
	r.ConID = s.ConID
	r.SiteID = s.SiteID
	r.ConProvID = s.ConProvID
	r.ConTypeID = s.ConTypeID
	r.ConCapID = s.ConCapID
	r.ConClassID = s.ConClassID
	r.InUse = s.InUse
	r.UpdatedOn = s.UpdatedOn
	r.CreatedOn = s.CreatedOn
	r.Hint = nullStringToPtr(s.Hint)
	r.Notes = nullStringToPtr(s.Notes)
}

// Return corresponding godevmandb create parameters
func (r *connection) createParams() godevmandb.CreateConnectionParams {
	s := godevmandb.CreateConnectionParams{}

	s.SiteID = r.SiteID
	s.ConProvID = r.ConProvID
	s.ConTypeID = r.ConTypeID
	s.ConCapID = r.ConCapID
	s.ConClassID = r.ConClassID
	s.InUse = r.InUse
	s.Hint = strToNullString(r.Hint)
	s.Notes = strToNullString(r.Notes)

	return s
}

// Return corresponding godevmandb update parameters
func (r *connection) updateParams() godevmandb.UpdateConnectionParams {
	s := godevmandb.UpdateConnectionParams{}

	s.SiteID = r.SiteID
	s.ConProvID = r.ConProvID
	s.ConTypeID = r.ConTypeID
	s.ConCapID = r.ConCapID
	s.ConClassID = r.ConClassID
	s.InUse = r.InUse
	s.Hint = strToNullString(r.Hint)
	s.Notes = strToNullString(r.Notes)

	return s
}

// Count Connections
// @Summary Count connections
// @Description Count number of connections
// @Tags connections
// @ID count-connections
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/count [GET]
func (h *Handler) CountConnections(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountConnections(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List connections
// @Summary List connections
// @Description List connection info
// @Tags connections
// @ID list-connections
// @Param hint_f query string false "url encoded SQL 'LIKE' operator pattern"
// @Param limit query int false "min: 1; max: 1000; default: 1000"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} connection
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections [GET]
func (h *Handler) GetConnections(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetConnectionsParams{
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
	d := r.FormValue("hint_f")
	if d != "" {
		p.HintF = d
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetConnections(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := []connection{}
	for _, s := range res {
		r := connection{}
		r.getValues(s)
		out = append(out, r)
	}

	RespondJSON(w, r, http.StatusOK, out)
}

// Get Connection
// @Summary Get connection
// @Description Get connection info
// @Tags connections
// @ID get-connection
// @Param con_id path string true "con_id"
// @Success 200 {object} connection
// @Failure 400 {object} StatusResponse "Invalid con_id"
// @Failure 404 {object} StatusResponse "Connection not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/{con_id} [GET]
func (h *Handler) GetConnection(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "con_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid connection ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetConnection(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "Connection not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	out := connection{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Create Connection
// @Summary Create connection
// @Description Create connection
// @Tags connections
// @ID create-connection
// @Param Body body connection true "JSON object of connection.<br />Ignored fields:<ul><li>con_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 201 {object} connection
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections [POST]
func (h *Handler) CreateConnection(w http.ResponseWriter, r *http.Request) {
	var pIn connection
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Create parameters for new db record
	p := pIn.createParams()

	q := godevmandb.New(h.db)
	res, err := q.CreateConnection(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := connection{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusCreated, out)
}

// Update Connection
// @Summary Update connection
// @Description Update connection
// @Tags connections
// @ID update-connection
// @Param con_id path string true "con_id"
// @Param Body body connection true "JSON object of connection.<br />Ignored fields:<ul><li>con_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 200 {object} connection
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/{con_id} [PUT]
func (h *Handler) UpdateConnection(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "con_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid connection ID")
		return
	}

	var pIn connection
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Update parameters for new db record
	p := pIn.updateParams()
	p.ConID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateConnection(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := connection{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Delete Connection
// @Summary Delete connection
// @Description Delete connection
// @Tags connections
// @ID delete-connection
// @Param con_id path string true "con_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid con_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/{con_id} [DELETE]
func (h *Handler) DeleteConnection(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "con_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid connection ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteConnection(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Foreign key
// Get Connection Capacitiy
// @Summary Get connection capacity
// @Description Get connection capacity info
// @Tags connections
// @ID get-connection-capacity
// @Param con_id path string true "con_id"
// @Success 200 {object} conCapacity
// @Failure 400 {object} StatusResponse "Invalid con_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/{con_id}/capacity [GET]
func (h *Handler) GetConnectionConCapacitiy(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "con_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid connection ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetConnectionConCapacitiy(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := conCapacity{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Foreign key
// Get Connection Class
// @Summary Get connection class
// @Description Get connection class info
// @Tags connections
// @ID get-connection-class
// @Param con_id path string true "con_id"
// @Success 200 {object} conClass
// @Failure 400 {object} StatusResponse "Invalid con_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/{con_id}/class [GET]
func (h *Handler) GetConnectionConClass(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "con_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid connection ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetConnectionConClass(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := conClass{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Foreign key
// Get Connection Provider
// @Summary Get connection provider
// @Description Get connection provider info
// @Tags connections
// @ID get-connection-provider
// @Param con_id path string true "con_id"
// @Success 200 {object} conProvider
// @Failure 400 {object} StatusResponse "Invalid con_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/{con_id}/provider [GET]
func (h *Handler) GetConnectionConProvider(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "con_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid connection ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetConnectionConProvider(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := conProvider{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Foreign key
// Get Connection Site
// @Summary Get connection site
// @Description Get connection site info
// @Tags connections
// @ID get-connection-site
// @Param con_id path string true "con_id"
// @Success 200 {object} site
// @Failure 400 {object} StatusResponse "Invalid con_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/{con_id}/site [GET]
// func (h *Handler) GetConnectionSite(w http.ResponseWriter, r *http.Request) {
// 	id, err := strconv.ParseInt(chi.URLParam(r, "con_id"), 10, 64)
// 	if err != nil {
// 		RespondError(w, r, http.StatusBadRequest, "Invalid connection ID")
// 		return
// 	}

// 	q := godevmandb.New(h.db)
// 	res, err := q.GetConnectionSite(h.ctx, id)
// 	if err != nil {
// 		RespondError(w, r, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	out := site{}
//  out.getValues(res)

// 	RespondJSON(w, r, http.StatusOK, out)
// }

// Foreign key
// Get Connection Type
// @Summary Get connection type
// @Description Get connection type info
// @Tags connections
// @ID get-connection-type
// @Param con_id path string true "con_id"
// @Success 200 {object} conType
// @Failure 400 {object} StatusResponse "Invalid con_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/{con_id}/type [GET]
func (h *Handler) GetConnectionConType(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "con_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid connection ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetConnectionConType(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := conType{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Relations
// List Connection Interfaces
// @Summary List connection interfaces
// @Description List connection interfaces info
// @Tags connections
// @ID list-connection-interfaces
// @Param con_id path string true "con_id"
// @Success 200 {array} interface
// @Failure 400 {object} StatusResponse "Invalid con_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /connections/{con_id}/interfaces [GET]
// func (h *Handler) GetConnectionInterfaces(w http.ResponseWriter, r *http.Request) {
// 	id, err := strconv.ParseInt(chi.URLParam(r, "con_id"), 10, 64)
// 	if err != nil {
// 		RespondError(w, r, http.StatusBadRequest, "Invalid connection ID")
// 		return
// 	}

// 	q := godevmandb.New(h.db)
// 	res, err := q.GetConnectionInterfaces(h.ctx, sql.NullInt64{Int64: id, Valid: true})
// 	if err != nil {
// 		RespondError(w, r, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	out := []interface{}
//  for _, s := range res {
//  	r := interface{}
//  	r.getValues(s)
//  	out = append(out, r)
//  }

// 	RespondJSON(w, r, http.StatusOK, out)
// }
