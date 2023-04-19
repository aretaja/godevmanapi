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
type xconnect struct {
	UpdatedOn   time.Time `json:"updated_on"`
	CreatedOn   time.Time `json:"created_on"`
	PeerIp      *string   `json:"peer_ip"`
	IfID        *int64    `json:"if_id"`
	PeerIfalias *string   `json:"peer_ifalias"`
	Xname       *string   `json:"xname"`
	Descr       *string   `json:"descr"`
	OpStat      *string   `json:"op_stat"`
	OpStatIn    *string   `json:"op_stat_in"`
	OpStatOut   *string   `json:"op_stat_out"`
	PeerDevID   *int64    `json:"peer_dev_id"`
	VcIdx       int64     `json:"vc_idx"`
	VcID        int64     `json:"vc_id"`
	XcID        int64     `json:"xc_id"`
	DevID       int64     `json:"dev_id"`
}

// Import values from corresponding godevmandb struct
func (r *xconnect) getValues(s godevmandb.Xconnect) {
	r.XcID = s.XcID
	r.DevID = s.DevID
	r.PeerDevID = s.PeerDevID
	r.IfID = s.IfID
	r.VcIdx = s.VcIdx
	r.VcID = s.VcID
	r.PeerIfalias = s.PeerIfalias
	r.Xname = s.Xname
	r.Descr = s.Descr
	r.OpStat = s.OpStat
	r.OpStatIn = s.OpStatIn
	r.OpStatOut = s.OpStatOut
	r.UpdatedOn = s.UpdatedOn
	r.CreatedOn = s.CreatedOn
	r.PeerIp = pgInetToPtr(s.PeerIp)
}

// Return corresponding godevmandb create parameters
func (r *xconnect) createParams() godevmandb.CreateXconnectParams {
	s := godevmandb.CreateXconnectParams{}

	s.DevID = r.DevID
	s.PeerDevID = r.PeerDevID
	s.IfID = r.IfID
	s.VcIdx = r.VcIdx
	s.VcID = r.VcID
	s.PeerIfalias = r.PeerIfalias
	s.Xname = r.Xname
	s.Descr = r.Descr
	s.OpStat = r.OpStat
	s.OpStatIn = r.OpStatIn
	s.OpStatOut = r.OpStatOut
	s.PeerIp = strToPgInet(r.PeerIp)

	return s
}

// Return corresponding godevmandb update parameters
func (r *xconnect) updateParams() godevmandb.UpdateXconnectParams {
	s := godevmandb.UpdateXconnectParams{}

	s.DevID = r.DevID
	s.PeerDevID = r.PeerDevID
	s.IfID = r.IfID
	s.VcIdx = r.VcIdx
	s.VcID = r.VcID
	s.PeerIfalias = r.PeerIfalias
	s.Xname = r.Xname
	s.Descr = r.Descr
	s.OpStat = r.OpStat
	s.OpStatIn = r.OpStatIn
	s.OpStatOut = r.OpStatOut
	s.PeerIp = strToPgInet(r.PeerIp)

	return s
}

// Count Xconnects
// @Summary Count xconnects
// @Description Count number of xconnects
// @Tags devices
// @ID count-xconnects
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/xconnects/count [GET]
func (h *Handler) CountXconnects(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountXconnects(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List xconnects
// @Summary List xconnects
// @Description List xconnects info
// @Tags devices
// @ID list-xconnects
// @Param vc_idx_f query string false "url encoded SQL '=' operator pattern"
// @Param vc_id_f query string false "url encoded SQL '=' operator pattern"
// @Param peer_ifalias_f query string false "url encoded SQL 'ILIKE' operator pattern + special value 'isnull', 'isempty'"
// @Param xname_f query string false "url encoded SQL 'ILIKE' operator pattern + special value 'isnull', 'isempty'"
// @Param descr_f query string false "url encoded SQL 'ILIKE' operator pattern + special value 'isnull', 'isempty'"
// @Param op_stat_f query string false "url encoded SQL 'ILIKE' operator pattern + special value 'isnull', 'isempty'"
// @Param op_stat_in_f query string false "url encoded SQL 'ILIKE' operator pattern + special value 'isnull', 'isempty'"
// @Param op_stat_out_f query string false "url encoded SQL 'ILIKE' operator pattern + special value 'isnull', 'isempty'"
// @Param peer_ip_f query string false "ip or containing net in CIDR notation"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} xconnect
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/xconnects [GET]
func (h *Handler) GetXconnects(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetXconnectsParams{
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
	if v := r.FormValue("vc_idx_f"); v != "" {
		p.VcIdxF = v
	}

	if v := r.FormValue("vc_id_f"); v != "" {
		p.VcIDF = v
	}

	if v := r.FormValue("peer_ifalias_f"); v != "" {
		p.PeerIfaliasF = &v
	}

	if v := r.FormValue("xname_f"); v != "" {
		p.XnameF = &v
	}

	if v := r.FormValue("descr_f"); v != "" {
		p.DescrF = &v
	}

	if v := r.FormValue("op_stat_f"); v != "" {
		p.OpStatF = &v
	}

	if v := r.FormValue("op_stat_in_f"); v != "" {
		p.OpStatInF = &v
	}

	if v := r.FormValue("op_stat_out_f"); v != "" {
		p.OpStatOutF = &v
	}

	p.PeerIpF = strToPgInet(nil)
	if v := r.FormValue("peer_ip_f"); v != "" {
		p.PeerIpF = strToPgInet(&v)
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetXconnects(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := []xconnect{}
	for _, s := range res {
		a := xconnect{}
		a.getValues(s)
		out = append(out, a)
	}

	RespondJSON(w, r, http.StatusOK, out)
}

// Get Xconnect
// @Summary Get xconnect
// @Description Get xconnect info
// @Tags devices
// @ID get-xconnect
// @Param xc_id path string true "xc_id"
// @Success 200 {object} xconnect
// @Failure 400 {object} StatusResponse "Invalid xc_id"
// @Failure 404 {object} StatusResponse "Xconnect not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/xconnects/{xc_id} [GET]
func (h *Handler) GetXconnect(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "xc_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid xconnect ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetXconnect(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "Xconnect not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	out := xconnect{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Create Xconnect
// @Summary Create xconnect
// @Description Create xconnect
// @Tags devices
// @ID create-xconnect
// @Param Body body xconnect true "JSON object of xconnect.<br />Ignored fields:<ul><li>xc_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 201 {object} xconnect
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/xconnects [POST]
func (h *Handler) CreateXconnect(w http.ResponseWriter, r *http.Request) {
	var pIn xconnect
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Create parameters for new db record
	p := pIn.createParams()

	q := godevmandb.New(h.db)
	res, err := q.CreateXconnect(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := xconnect{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusCreated, out)
}

// Update Xconnect
// @Summary Update xconnect
// @Description Update xconnect
// @Tags devices
// @ID update-xconnect
// @Param xc_id path string true "xc_id"
// @Param Body body xconnect true "JSON object of xconnect.<br />Ignored fields:<ul><li>xc_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 200 {object} xconnect
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/xconnects/{xc_id} [PUT]
func (h *Handler) UpdateXconnect(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "xc_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid xconnect ID")
		return
	}

	var pIn xconnect
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
	res, err := q.UpdateXconnect(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := xconnect{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Delete Xconnect
// @Summary Delete xconnect
// @Description Delete xconnect
// @Tags devices
// @ID delete-xconnect
// @Param xc_id path string true "xc_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid xc_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/xconnects/{xc_id} [DELETE]
func (h *Handler) DeleteXconnect(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "xc_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid xconnect ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteXconnect(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Foreign key
// Get Xconnect Device
// @Summary Get xconnect device
// @Description Get xconnect device info
// @Tags devices
// @ID get-xconnect-device
// @Param xc_id path string true "xc_id"
// @Success 200 {object} device
// @Failure 400 {object} StatusResponse "Invalid xc_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/xconnects/{xc_id}/device [GET]
func (h *Handler) GetXconnectDevice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "xc_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid xconnect ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetXconnectDevice(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := device{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Foreign key
// Get Xconnect peer Device
// @Summary Get xconnect peer device
// @Description Get xconnect peer device info
// @Tags devices
// @ID get-xconnect-peer-device
// @Param xc_id path string true "xc_id"
// @Success 200 {object} device
// @Failure 400 {object} StatusResponse "Invalid xc_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/xconnects/{xc_id}/peer_device [GET]
func (h *Handler) GetXconnectPeerDevice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "xc_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid xconnect ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetXconnectPeerDevice(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := device{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Foreign key
// Get Xconnect Interface
// @Summary Get xconnect interface
// @Description Get xconnect interface info
// @Tags devices
// @ID get-xconnect-interface
// @Param xc_id path string true "xc_id"
// @Success 200 {object} device
// @Failure 400 {object} StatusResponse "Invalid xc_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/xconnects/{xc_id}/interface [GET]
func (h *Handler) GetXconnectInterface(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "xc_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid xconnect ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetXconnectInterface(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := iface{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}
