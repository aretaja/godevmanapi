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
type ospfNbr struct {
	UpdatedOn time.Time `json:"updated_on"`
	CreatedOn time.Time `json:"created_on"`
	NbrIp     *string   `json:"nbr_ip"`
	Condition *string   `json:"condition"`
	NbrID     int64     `json:"nbr_id"`
	DevID     int64     `json:"dev_id"`
}

// Import values from corresponding godevmandb struct
func (r *ospfNbr) getValues(s godevmandb.OspfNbr) {
	r.NbrID = s.NbrID
	r.DevID = s.DevID
	r.Condition = s.Condition
	r.UpdatedOn = s.UpdatedOn
	r.CreatedOn = s.CreatedOn
	r.NbrIp = pgInetToPtr(s.NbrIp)
}

// Return corresponding godevmandb create parameters
func (r *ospfNbr) createParams() godevmandb.CreateOspfNbrParams {
	s := godevmandb.CreateOspfNbrParams{}

	s.DevID = r.DevID
	s.Condition = r.Condition
	s.NbrIp = strToPgInet(r.NbrIp)

	return s
}

// Return corresponding godevmandb update parameters
func (r *ospfNbr) updateParams() godevmandb.UpdateOspfNbrParams {
	s := godevmandb.UpdateOspfNbrParams{}

	s.DevID = r.DevID
	s.Condition = r.Condition
	s.NbrIp = strToPgInet(r.NbrIp)

	return s
}

// Count OspfNbrs
// @Summary Count ospf_nbrs
// @Description Count number of ospf_nbrs
// @Tags devices
// @ID count-ospf_nbrs
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/ospf_nbrs/count [GET]
func (h *Handler) CountOspfNbrs(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountOspfNbrs(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List ospf_nbrs
// @Summary List ospf_nbrs
// @Description List ospf_nbrs info
// @Tags devices
// @ID list-ospf_nbrs
// @Param dev_id_f query string false "url encoded SQL 'LIKE' operator pattern"
// @Param condition_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param nbr_ip_f query string false "ip or containing net in CIDR notation"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} ospfNbr
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/ospf_nbrs [GET]
func (h *Handler) GetOspfNbrs(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetOspfNbrsParams{
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

	// SysID filter
	v := r.FormValue("dev_id_f")
	if v != "" {
		p.DevIDF = v
	}

	// Condition filter
	v = r.FormValue("condition_f")
	if v != "" {
		p.ConditionF = &v
	}

	// Host IPv4 filter
	p.NbrIpF = strToPgInet(nil)
	v = r.FormValue("nbr_ip_f")
	if v != "" {
		p.NbrIpF = strToPgInet(&v)
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetOspfNbrs(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := []ospfNbr{}
	for _, s := range res {
		a := ospfNbr{}
		a.getValues(s)
		out = append(out, a)
	}

	RespondJSON(w, r, http.StatusOK, out)
}

// Get OspfNbr
// @Summary Get ospf_nbr
// @Description Get ospf_nbr info
// @Tags devices
// @ID get-ospf_nbr
// @Param nbr_id path string true "nbr_id"
// @Success 200 {object} ospfNbr
// @Failure 400 {object} StatusResponse "Invalid nbr_id"
// @Failure 404 {object} StatusResponse "OspfNbr not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/ospf_nbrs/{nbr_id} [GET]
func (h *Handler) GetOspfNbr(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "nbr_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid ospf_nbr ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetOspfNbr(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "OspfNbr not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	out := ospfNbr{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Create OspfNbr
// @Summary Create ospf_nbr
// @Description Create ospf_nbr
// @Tags devices
// @ID create-ospf_nbr
// @Param Body body ospfNbr true "JSON object of ospfNbr.<br />Ignored fields:<ul><li>nbr_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 201 {object} ospfNbr
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/ospf_nbrs [POST]
func (h *Handler) CreateOspfNbr(w http.ResponseWriter, r *http.Request) {
	var pIn ospfNbr
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Create parameters for new db record
	p := pIn.createParams()

	q := godevmandb.New(h.db)
	res, err := q.CreateOspfNbr(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := ospfNbr{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusCreated, out)
}

// Update OspfNbr
// @Summary Update ospf_nbr
// @Description Update ospf_nbr
// @Tags devices
// @ID update-ospf_nbr
// @Param nbr_id path string true "nbr_id"
// @Param Body body ospfNbr true "JSON object of ospfNbr.<br />Ignored fields:<ul><li>nbr_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 200 {object} ospfNbr
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/ospf_nbrs/{nbr_id} [PUT]
func (h *Handler) UpdateOspfNbr(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "nbr_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid ospf_nbr ID")
		return
	}

	var pIn ospfNbr
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Update parameters for new db record
	p := pIn.updateParams()
	p.DevID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateOspfNbr(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := ospfNbr{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Delete OspfNbr
// @Summary Delete ospf_nbr
// @Description Delete ospf_nbr
// @Tags devices
// @ID delete-ospf_nbr
// @Param nbr_id path string true "nbr_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid nbr_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/ospf_nbrs/{nbr_id} [DELETE]
func (h *Handler) DeleteOspfNbr(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "nbr_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid ospf_nbr ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteOspfNbr(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Foreign key
// Get OspfNbr Device
// @Summary Get ospf_nbr device
// @Description Get ospf_nbr device info
// @Tags devices
// @ID get-ospf_nbr-device
// @Param nbr_id path string true "nbr_id"
// @Success 200 {object} device
// @Failure 400 {object} StatusResponse "Invalid nbr_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/ospf_nbrs/{nbr_id}/device [GET]
func (h *Handler) GetOspfNbrDevice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "nbr_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid ospf_nbr ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetOspfNbrDevice(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := device{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}
