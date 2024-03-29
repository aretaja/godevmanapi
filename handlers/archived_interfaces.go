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
type archivedInterface struct {
	UpdatedOn          time.Time `json:"updated_on"`
	CreatedOn          time.Time `json:"created_on"`
	HostIp6            *string   `json:"host_ip6"`
	CiscoOptPowerIndex *string   `json:"cisco_opt_power_index"`
	HostIp4            *string   `json:"host_ip4"`
	Alias              *string   `json:"alias"`
	TypeEnum           *int16    `json:"type_enum"`
	Mac                *string   `json:"mac"`
	OtnIfID            *int64    `json:"otn_if_id"`
	Ifindex            *int64    `json:"ifindex"`
	Hostname           string    `json:"hostname"`
	Manufacturer       string    `json:"manufacturer"`
	Model              string    `json:"model"`
	Descr              string    `json:"descr"`
	IfaID              int64     `json:"ifa_id"`
}

// Import values from corresponding godevmandb struct
func (r *archivedInterface) getValues(s godevmandb.ArchivedInterface) {
	r.IfaID = s.IfaID
	r.Hostname = s.Hostname
	r.Manufacturer = s.Manufacturer
	r.Model = s.Model
	r.Descr = s.Descr
	r.UpdatedOn = s.UpdatedOn
	r.CreatedOn = s.CreatedOn
	r.Ifindex = s.Ifindex
	r.OtnIfID = s.OtnIfID
	r.CiscoOptPowerIndex = s.CiscoOptPowerIndex
	r.HostIp4 = pgInetToPtr(s.HostIp4)
	r.HostIp6 = pgInetToPtr(s.HostIp6)
	r.Alias = s.Alias
	r.TypeEnum = s.TypeEnum
	r.Mac = pgMacaddrToPtr(s.Mac)
}

// Return corresponding godevmandb create parameters
func (r *archivedInterface) createParams() godevmandb.CreateArchivedInterfaceParams {
	s := godevmandb.CreateArchivedInterfaceParams{}

	s.Hostname = r.Hostname
	s.Manufacturer = r.Manufacturer
	s.Model = r.Model
	s.Descr = r.Descr
	s.Ifindex = r.Ifindex
	s.OtnIfID = r.OtnIfID
	s.CiscoOptPowerIndex = r.CiscoOptPowerIndex
	s.HostIp4 = strToPgInet(r.HostIp4)
	s.HostIp6 = strToPgInet(r.HostIp6)
	s.Alias = r.Alias
	s.TypeEnum = r.TypeEnum
	s.Mac = strToPgMacaddr(r.Mac)

	return s
}

// Return corresponding godevmandb update parameters
func (r *archivedInterface) updateParams() godevmandb.UpdateArchivedInterfaceParams {
	s := godevmandb.UpdateArchivedInterfaceParams{}

	s.Hostname = r.Hostname
	s.Manufacturer = r.Manufacturer
	s.Model = r.Model
	s.Descr = r.Descr
	s.Ifindex = r.Ifindex
	s.OtnIfID = r.OtnIfID
	s.CiscoOptPowerIndex = r.CiscoOptPowerIndex
	s.HostIp4 = strToPgInet(r.HostIp4)
	s.HostIp6 = strToPgInet(r.HostIp6)
	s.Alias = r.Alias
	s.TypeEnum = r.TypeEnum
	s.Mac = strToPgMacaddr(r.Mac)

	return s
}

// Count ArchivedInterfaces
// @Summary Count archived_interfaces
// @Description Count number of archived interfaces
// @Tags archived
// @ID count-archived_interfaces
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /archived/interfaces/count [GET]
func (h *Handler) CountArchivedInterfaces(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountArchivedInterfaces(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List archived_interfaces
// @Summary List archived_interfaces
// @Description List archived interfaces info
// @Tags archived
// @ID list-archived_interfaces
// @Param ifindex_f query string false "url encoded SQL 'LIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param otn_if_id_f query string false "url encoded SQL 'LIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param cisco_opt_power_index_f query string false "url encoded SQL 'LIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param hostname_f query string false "url encoded SQL 'ILIKE' operator pattern + special value 'isempty'"
// @Param host_ip4_f query string false "ip or containing net in CIDR notation"
// @Param host_ip6_f query string false "ip or containing net in CIDR notation"
// @Param manufacturer_f query string false "url encoded SQL 'ILIKE' operator pattern + special value 'isempty'"
// @Param model_f query string false "url encoded SQL 'ILIKE' operator pattern + special value 'isempty'"
// @Param descr_f query string false "url encoded SQL 'ILIKE' operator pattern + special value 'isempty'"
// @Param alias_f query string false "url encoded SQL 'ILIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param type_enum_f query string false "url encoded SQL 'LIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param mac_f query string false "SQL '=' operator value (MAC address)"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} archivedInterface
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /archived/interfaces [GET]
func (h *Handler) GetArchivedInterfaces(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetArchivedInterfacesParams{
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
	if v := r.FormValue("ifindex_f"); v != "" {
		p.IfindexF = &v
	}

	if v := r.FormValue("otn_if_id_f"); v != "" {
		p.OtnIfIDF = &v
	}

	if v := r.FormValue("cisco_opt_power_index_f"); v != "" {
		p.CiscoOptPowerIndexF = &v
	}

	if v := r.FormValue("hostname_f"); v != "" {
		p.HostnameF = v
	}

	p.HostIp4F = strToPgInet(nil)
	if v := r.FormValue("host_ip4_f"); v != "" {
		p.HostIp4F = strToPgInet(&v)
	}

	p.HostIp6F = strToPgInet(nil)
	if v := r.FormValue("host_ip6_f"); v != "" {
		p.HostIp6F = strToPgInet(&v)
	}

	if v := r.FormValue("manufacturer_f"); v != "" {
		p.ManufacturerF = v
	}

	if v := r.FormValue("model_f"); v != "" {
		p.ModelF = v
	}

	if v := r.FormValue("descr_f"); v != "" {
		p.DescrF = v
	}

	if v := r.FormValue("alias_f"); v != "" {
		p.AliasF = &v
	}

	if v := r.FormValue("type_enum_f"); v != "" {
		p.TypeEnumF = &v
	}

	p.MacF = strToPgMacaddr(nil)
	if v := r.FormValue("mac_f"); v != "" {
		p.MacF = strToPgMacaddr(&v)
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetArchivedInterfaces(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := []archivedInterface{}
	for _, s := range res {
		a := archivedInterface{}
		a.getValues(s)
		out = append(out, a)
	}

	RespondJSON(w, r, http.StatusOK, out)
}

// Get ArchivedInterface
// @Summary Get archived_interface
// @Description Get archived interface info
// @Tags archived
// @ID get-archived_interface
// @Param ifa_id path string true "ifa_id"
// @Success 200 {object} archivedInterface
// @Failure 400 {object} StatusResponse "Invalid ifa_id"
// @Failure 404 {object} StatusResponse "Archived interface not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /archived/interfaces/{ifa_id} [GET]
func (h *Handler) GetArchivedInterface(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ifa_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid archived interface ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetArchivedInterface(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "Archived interface not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	out := archivedInterface{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Create ArchivedInterface
// @Summary Create archived_interface
// @Description Create archived interface
// @Tags archived
// @ID create-archived_interface
// @Param Body body archivedInterface true "JSON object of archivedInterface.<br />Ignored fields:<ul><li>ifa_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 201 {object} archivedInterface
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /archived/interfaces [POST]
func (h *Handler) CreateArchivedInterface(w http.ResponseWriter, r *http.Request) {
	var pIn archivedInterface
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Create parameters for new db record
	p := pIn.createParams()

	q := godevmandb.New(h.db)
	res, err := q.CreateArchivedInterface(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := archivedInterface{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusCreated, out)
}

// Update ArchivedInterface
// @Summary Update archived_interface
// @Description Update archived interface
// @Tags archived
// @ID update-archived_interface
// @Param ifa_id path string true "ifa_id"
// @Param Body body archivedInterface true "JSON object of archivedInterface.<br />Ignored fields:<ul><li>ifa_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 200 {object} archivedInterface
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /archived/interfaces/{ifa_id} [PUT]
func (h *Handler) UpdateArchivedInterface(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ifa_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid archived interface ID")
		return
	}

	var pIn archivedInterface
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Update parameters for new db record
	p := pIn.updateParams()
	p.IfaID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateArchivedInterface(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := archivedInterface{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Delete ArchivedInterface
// @Summary Delete archived_interface
// @Description Delete archived interface
// @Tags archived
// @ID delete-archived_interface
// @Param ifa_id path string true "ifa_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid ifa_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /archived/interfaces/{ifa_id} [DELETE]
func (h *Handler) DeleteArchivedInterface(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ifa_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid archived interface ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteArchivedInterface(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
