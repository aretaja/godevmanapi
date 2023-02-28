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
type ipInterface struct {
	UpdatedOn time.Time `json:"updated_on"`
	CreatedOn time.Time `json:"created_on"`
	Ifindex   *int64    `json:"ifindex"`
	IpAddr    *string   `json:"ip_addr"`
	Descr     *string   `json:"descr"`
	Alias     *string   `json:"alias"`
	IpID      int64     `json:"ip_id"`
	DevID     int64     `json:"dev_id"`
}

// Import values from corresponding godevmandb struct
func (r *ipInterface) getValues(s godevmandb.IpInterface) {
	r.IpID = s.IpID
	r.DevID = s.DevID
	r.Ifindex = s.Ifindex
	r.Descr = s.Descr
	r.Alias = s.Alias
	r.UpdatedOn = s.UpdatedOn
	r.CreatedOn = s.CreatedOn
	r.IpAddr = pgInetToPtr(s.IpAddr)
}

// Return corresponding godevmandb create parameters
func (r *ipInterface) createParams() godevmandb.CreateIpInterfaceParams {
	s := godevmandb.CreateIpInterfaceParams{}

	s.DevID = r.DevID
	s.Ifindex = r.Ifindex
	s.Descr = r.Descr
	s.Alias = r.Alias
	s.IpAddr = strToPgInet(r.IpAddr)

	return s
}

// Return corresponding godevmandb update parameters
func (r *ipInterface) updateParams() godevmandb.UpdateIpInterfaceParams {
	s := godevmandb.UpdateIpInterfaceParams{}

	s.DevID = r.DevID
	s.Ifindex = r.Ifindex
	s.Descr = r.Descr
	s.Alias = r.Alias
	s.IpAddr = strToPgInet(r.IpAddr)

	return s
}

// Count IpInterfaces
// @Summary Count ip_interfaces
// @Description Count number of ip_interfaces
// @Tags ip_interfaces
// @ID count-ip_interfaces
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /ip_interfaces/count [GET]
func (h *Handler) CountIpInterfaces(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountIpInterfaces(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List ip_interfaces
// @Summary List ip_interfaces
// @Description List ip_interfaces info
// @Tags ip_interfaces
// @ID list-ip_interfaces
// @Param dev_id_f query string false "url encoded SQL 'LIKE' operator pattern"
// @Param ifindex_f query string false "url encoded SQL 'LIKE' operator pattern"
// @Param descr_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param alias_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param ip_addr_f query string false "ip or containing net in CIDR notation"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} ipInterface
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /ip_interfaces [GET]
func (h *Handler) GetIpInterfaces(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetIpInterfacesParams{
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

	// Software filter
	v = r.FormValue("ifindex_f")
	if v != "" {
		p.IfindexF = &v
	}

	// Notes filter
	v = r.FormValue("descr_f")
	if v != "" {
		p.DescrF = &v
	}

	// Name filter
	v = r.FormValue("alias_f")
	if v != "" {
		p.AliasF = &v
	}

	// Host IPv4 filter
	p.IpAddrF = strToPgInet(nil)
	v = r.FormValue("ip_addr_f")
	if v != "" {
		p.IpAddrF = strToPgInet(&v)
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetIpInterfaces(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := []ipInterface{}
	for _, s := range res {
		a := ipInterface{}
		a.getValues(s)
		out = append(out, a)
	}

	RespondJSON(w, r, http.StatusOK, out)
}

// Get IpInterface
// @Summary Get ip_interface
// @Description Get ip_interface info
// @Tags ip_interfaces
// @ID get-ip_interface
// @Param ip_id path string true "ip_id"
// @Success 200 {object} ipInterface
// @Failure 400 {object} StatusResponse "Invalid ip_id"
// @Failure 404 {object} StatusResponse "IpInterface not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /ip_interfaces/{ip_id} [GET]
func (h *Handler) GetIpInterface(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ip_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid ip_interface ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetIpInterface(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "IpInterface not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	out := ipInterface{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Create IpInterface
// @Summary Create ip_interface
// @Description Create ip_interface
// @Tags ip_interfaces
// @ID create-ip_interface
// @Param Body body ipInterface true "JSON object of ipInterface.<br />Ignored fields:<ul><li>ip_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 201 {object} ipInterface
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /ip_interfaces [POST]
func (h *Handler) CreateIpInterface(w http.ResponseWriter, r *http.Request) {
	var pIn ipInterface
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Create parameters for new db record
	p := pIn.createParams()

	q := godevmandb.New(h.db)
	res, err := q.CreateIpInterface(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := ipInterface{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusCreated, out)
}

// Update IpInterface
// @Summary Update ip_interface
// @Description Update ip_interface
// @Tags ip_interfaces
// @ID update-ip_interface
// @Param ip_id path string true "ip_id"
// @Param Body body ipInterface true "JSON object of ipInterface.<br />Ignored fields:<ul><li>ip_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 200 {object} ipInterface
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /ip_interfaces/{ip_id} [PUT]
func (h *Handler) UpdateIpInterface(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ip_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid ip_interface ID")
		return
	}

	var pIn ipInterface
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
	res, err := q.UpdateIpInterface(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := ipInterface{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Delete IpInterface
// @Summary Delete ip_interface
// @Description Delete ip_interface
// @Tags ip_interfaces
// @ID delete-ip_interface
// @Param ip_id path string true "ip_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid ip_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /ip_interfaces/{ip_id} [DELETE]
func (h *Handler) DeleteIpInterface(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ip_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid ip_interface ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteIpInterface(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Foreign key
// Get IpInterface Device
// @Summary Get ip_interface device
// @Description Get ip_interface device info
// @Tags ip_interfaces
// @ID get-ip_interface-device
// @Param ip_id path string true "ip_id"
// @Success 200 {object} device
// @Failure 400 {object} StatusResponse "Invalid ip_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /ip_interfaces/{ip_id}/device [GET]
func (h *Handler) GetIpInterfaceDevice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "if_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid ip_interface ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetIpInterfaceDevice(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := device{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}
