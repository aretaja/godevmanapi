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
type archivedSubInterface struct {
	UpdatedOn time.Time `json:"updated_on"`
	CreatedOn time.Time `json:"created_on"`
	HostIp6   *string   `json:"host_ip6"`
	HostIp4   *string   `json:"host_ip4"`
	Alias     *string   `json:"alias"`
	Notes     *string   `json:"notes"`
	Type      *string   `json:"type"`
	Mac       *string   `json:"mac"`
	Ifindex   *int64    `json:"ifindex"`
	Hostname  string    `json:"hostname"`
	Descr     string    `json:"descr"`
	SifaID    int64     `json:"sifa_id"`
}

// Import values from corresponding godevmandb struct
func (r *archivedSubInterface) getValues(s godevmandb.ArchivedSubInterface) {
	r.SifaID = s.SifaID
	r.Hostname = s.Hostname
	r.Descr = s.Descr
	r.UpdatedOn = s.UpdatedOn
	r.CreatedOn = s.CreatedOn
	r.Ifindex = nullInt64ToPtr(s.Ifindex)
	r.HostIp4 = pgInetToPtr(s.HostIp4)
	r.HostIp6 = pgInetToPtr(s.HostIp6)
	r.Alias = nullStringToPtr(s.Alias)
	r.Notes = nullStringToPtr(s.Notes)
	r.Type = nullStringToPtr(s.Type)
	r.Mac = pgMacaddrToPtr(s.Mac)
}

// Return corresponding godevmandb create parameters
func (r *archivedSubInterface) createParams() godevmandb.CreateArchivedSubInterfaceParams {
	s := godevmandb.CreateArchivedSubInterfaceParams{}

	s.Hostname = r.Hostname
	s.Descr = r.Descr
	s.Ifindex = int64ToNullInt64(r.Ifindex)
	s.HostIp4 = strToPgInet(r.HostIp4)
	s.HostIp6 = strToPgInet(r.HostIp6)
	s.Alias = strToNullString(r.Alias)
	s.Notes = strToNullString(r.Notes)
	s.Type = nullStringToPtr(r.Type)
	s.Mac = strToPgMacaddr(r.Mac)

	return s
}

// Return corresponding godevmandb update parameters
func (r *archivedSubInterface) updateParams() godevmandb.UpdateArchivedSubInterfaceParams {
	s := godevmandb.UpdateArchivedSubInterfaceParams{}

	s.Hostname = r.Hostname
	s.Descr = r.Descr
	s.Ifindex = int64ToNullInt64(r.Ifindex)
	s.HostIp4 = strToPgInet(r.HostIp4)
	s.HostIp6 = strToPgInet(r.HostIp6)
	s.Alias = strToNullString(r.Alias)
	s.Notes = strToNullString(r.Notes)
	s.Type = nullStringToPtr(r.Type)
	s.Mac = strToPgMacaddr(r.Mac)

	return s
}

// Count ArchivedSubInterfaces
// @Summary Count archived_subinterfaces
// @Description Count number of archived subinterfaces
// @Tags archived
// @ID count-archived_subinterfaces
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /archived/subinterfaces/count [GET]
func (h *Handler) CountArchivedSubInterfaces(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountArchivedSubInterfaces(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List archived_subinterfaces
// @Summary List archived_subinterfaces
// @Description List archived subinterfaces info
// @Tags archived
// @ID list-archived_subinterfaces
// @Param ifindex_f query string false "url encoded SQL 'LIKE' operator pattern"
// @Param hostname_f query string false "url encoded SQL 'LIKE' operator pattern"
// @Param descr_f query string false "url encoded SQL 'LIKE' operator pattern"
// @Param alias_f query string false "url encoded SQL 'LIKE' operator pattern"
// @Param host_ip4_f query string false "ip or containing net in CIDR notation"
// @Param host_ip6_f query string false "ip or containing net in CIDR notation"
// @Param mac_f query string false "SQL '=' operator value (MAC address)"
// @Param limit query int false "min: 1; max: 1000; default: 1000"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} archivedSubInterface
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /archived/subinterfaces [GET]
func (h *Handler) GetArchivedSubInterfaces(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetArchivedSubInterfacesParams{
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

	// Ifindex filter
	v := r.FormValue("ifindex_f")
	if v != "" {
		p.IfindexF = strToNullString(&v)
	}

	// Alias filter
	v = r.FormValue("alias_f")
	if v != "" {
		p.AliasF = strToNullString(&v)
	}

	// Host IPv4 filter
	p.HostIp4F = strToPgInet(nil)
	v = r.FormValue("host_ip4_f")
	if v != "" {
		p.HostIp4F = strToPgInet(&v)
	}

	// Host IPv6 filter
	p.HostIp6F = strToPgInet(nil)
	v = r.FormValue("host_ip6_f")
	if v != "" {
		p.HostIp6F = strToPgInet(&v)
	}

	// MAC filter
	p.MacF = strToPgMacaddr(nil)
	v = r.FormValue("mac_f")
	if v != "" {
		p.MacF = strToPgMacaddr(&v)
	}

	// Hostname filter
	v = r.FormValue("hostname_f")
	if v != "" {
		p.HostnameF = v
	}

	// Descr filter
	v = r.FormValue("descr_f")
	if v != "" {
		p.DescrF = v
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetArchivedSubInterfaces(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := []archivedSubInterface{}
	for _, s := range res {
		a := archivedSubInterface{}
		a.getValues(s)
		out = append(out, a)
	}

	RespondJSON(w, r, http.StatusOK, out)
}

// Get ArchivedSubInterface
// @Summary Get archived_subinterface
// @Description Get archived subinterface info
// @Tags archived
// @ID get-archived_subinterface
// @Param sifa_id path string true "sifa_id"
// @Success 200 {object} archivedSubInterface
// @Failure 400 {object} StatusResponse "Invalid sifa_id"
// @Failure 404 {object} StatusResponse "Archived subinterface not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /archived/subinterfaces/{sifa_id} [GET]
func (h *Handler) GetArchivedSubInterface(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "sifa_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid archived subinterface ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetArchivedSubInterface(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "Archived subinterface not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	out := archivedSubInterface{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Create ArchivedSubInterface
// @Summary Create archived_subinterface
// @Description Create archived subinterface
// @Tags archived
// @ID create-archived_subinterface
// @Param Body body archivedSubInterface true "JSON object of archivedSubInterface.<br />Ignored fields:<ul><li>sifa_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 201 {object} archivedSubInterface
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /archived/subinterfaces [POST]
func (h *Handler) CreateArchivedSubInterface(w http.ResponseWriter, r *http.Request) {
	var pIn archivedSubInterface
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Create parameters for new db record
	p := pIn.createParams()

	q := godevmandb.New(h.db)
	res, err := q.CreateArchivedSubInterface(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := archivedSubInterface{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusCreated, out)
}

// Update ArchivedSubInterface
// @Summary Update archived_subinterface
// @Description Update archived subinterface
// @Tags archived
// @ID update-archived_subinterface
// @Param sifa_id path string true "sifa_id"
// @Param Body body archivedSubInterface true "JSON object of archivedSubInterface.<br />Ignored fields:<ul><li>sifa_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 200 {object} archivedSubInterface
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /archived/subinterfaces/{sifa_id} [PUT]
func (h *Handler) UpdateArchivedSubInterface(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "sifa_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid archived subinterface ID")
		return
	}

	var pIn archivedSubInterface
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Update parameters for new db record
	p := pIn.updateParams()
	p.SifaID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateArchivedSubInterface(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := archivedSubInterface{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Delete ArchivedSubInterface
// @Summary Delete archived_subinterface
// @Description Delete archived subinterface
// @Tags archived
// @ID delete-archived_subinterface
// @Param sifa_id path string true "sifa_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid sifa_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /archived/subinterfaces/{sifa_id} [DELETE]
func (h *Handler) DeleteArchivedSubInterface(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "sifa_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid archived subinterface ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteArchivedSubInterface(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
