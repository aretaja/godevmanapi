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
type archivedSubinterface struct {
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
func (r *archivedSubinterface) getValues(s godevmandb.ArchivedSubinterface) {
	r.SifaID = s.SifaID
	r.Hostname = s.Hostname
	r.Descr = s.Descr
	r.UpdatedOn = s.UpdatedOn
	r.CreatedOn = s.CreatedOn
	r.Ifindex = s.Ifindex
	r.HostIp4 = pgInetToPtr(s.HostIp4)
	r.HostIp6 = pgInetToPtr(s.HostIp6)
	r.Alias = s.Alias
	r.Notes = s.Notes
	r.Type = s.Type
	r.Mac = pgMacaddrToPtr(s.Mac)
}

// Return corresponding godevmandb create parameters
func (r *archivedSubinterface) createParams() godevmandb.CreateArchivedSubinterfaceParams {
	s := godevmandb.CreateArchivedSubinterfaceParams{}

	s.Hostname = r.Hostname
	s.Descr = r.Descr
	s.Ifindex = r.Ifindex
	s.HostIp4 = strToPgInet(r.HostIp4)
	s.HostIp6 = strToPgInet(r.HostIp6)
	s.Alias = r.Alias
	s.Notes = r.Notes
	s.Type = r.Type
	s.Mac = strToPgMacaddr(r.Mac)

	return s
}

// Return corresponding godevmandb update parameters
func (r *archivedSubinterface) updateParams() godevmandb.UpdateArchivedSubinterfaceParams {
	s := godevmandb.UpdateArchivedSubinterfaceParams{}

	s.Hostname = r.Hostname
	s.Descr = r.Descr
	s.Ifindex = r.Ifindex
	s.HostIp4 = strToPgInet(r.HostIp4)
	s.HostIp6 = strToPgInet(r.HostIp6)
	s.Alias = r.Alias
	s.Notes = r.Notes
	s.Type = r.Type
	s.Mac = strToPgMacaddr(r.Mac)

	return s
}

// Count ArchivedSubinterfaces
// @Summary Count archived_subinterfaces
// @Description Count number of archived subinterfaces
// @Tags archived
// @ID count-archived_subinterfaces
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /archived/subinterfaces/count [GET]
func (h *Handler) CountArchivedSubinterfaces(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountArchivedSubinterfaces(h.ctx)
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
// @Param ifindex_f query string false "url encoded SQL 'LIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param descr_f query string false "url encoded SQL 'ILIKE' operator pattern + special value 'isempty'"
// @Param parent_descr_f query string false "url encoded SQL 'ILIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param alias_f query string false "url encoded SQL 'ILIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param type_f query string false "url encoded SQL 'ILIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param mac_f query string false "SQL '=' operator value (MAC address)"
// @Param hostname_f query string false "url encoded SQL 'ILIKE' operator pattern + special value 'isempty'"
// @Param host_ip4_f query string false "ip or containing net in CIDR notation"
// @Param host_ip6_f query string false "ip or containing net in CIDR notation"
// @Param notes_f query string false "url encoded SQL 'ILIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} archivedSubinterface
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /archived/subinterfaces [GET]
func (h *Handler) GetArchivedSubinterfaces(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetArchivedSubinterfacesParams{
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

	if v := r.FormValue("descr_f"); v != "" {
		p.DescrF = v
	}

	if v := r.FormValue("parent_descr_f"); v != "" {
		p.ParentDescrF = &v
	}

	if v := r.FormValue("alias_f"); v != "" {
		p.AliasF = &v
	}

	if v := r.FormValue("type_f"); v != "" {
		p.TypeF = &v
	}

	p.MacF = strToPgMacaddr(nil)
	if v := r.FormValue("mac_f"); v != "" {
		p.MacF = strToPgMacaddr(&v)
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

	if v := r.FormValue("notes_f"); v != "" {
		p.NotesF = &v
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetArchivedSubinterfaces(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := []archivedSubinterface{}
	for _, s := range res {
		a := archivedSubinterface{}
		a.getValues(s)
		out = append(out, a)
	}

	RespondJSON(w, r, http.StatusOK, out)
}

// Get ArchivedSubinterface
// @Summary Get archived_subinterface
// @Description Get archived subinterface info
// @Tags archived
// @ID get-archived_subinterface
// @Param sifa_id path string true "sifa_id"
// @Success 200 {object} archivedSubinterface
// @Failure 400 {object} StatusResponse "Invalid sifa_id"
// @Failure 404 {object} StatusResponse "Archived subinterface not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /archived/subinterfaces/{sifa_id} [GET]
func (h *Handler) GetArchivedSubinterface(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "sifa_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid archived subinterface ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetArchivedSubinterface(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "Archived subinterface not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	out := archivedSubinterface{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Create ArchivedSubinterface
// @Summary Create archived_subinterface
// @Description Create archived subinterface
// @Tags archived
// @ID create-archived_subinterface
// @Param Body body archivedSubinterface true "JSON object of archivedSubinterface.<br />Ignored fields:<ul><li>sifa_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 201 {object} archivedSubinterface
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /archived/subinterfaces [POST]
func (h *Handler) CreateArchivedSubinterface(w http.ResponseWriter, r *http.Request) {
	var pIn archivedSubinterface
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Create parameters for new db record
	p := pIn.createParams()

	q := godevmandb.New(h.db)
	res, err := q.CreateArchivedSubinterface(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := archivedSubinterface{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusCreated, out)
}

// Update ArchivedSubinterface
// @Summary Update archived_subinterface
// @Description Update archived subinterface
// @Tags archived
// @ID update-archived_subinterface
// @Param sifa_id path string true "sifa_id"
// @Param Body body archivedSubinterface true "JSON object of archivedSubinterface.<br />Ignored fields:<ul><li>sifa_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 200 {object} archivedSubinterface
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /archived/subinterfaces/{sifa_id} [PUT]
func (h *Handler) UpdateArchivedSubinterface(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "sifa_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid archived subinterface ID")
		return
	}

	var pIn archivedSubinterface
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
	res, err := q.UpdateArchivedSubinterface(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := archivedSubinterface{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Delete ArchivedSubinterface
// @Summary Delete archived_subinterface
// @Description Delete archived subinterface
// @Tags archived
// @ID delete-archived_subinterface
// @Param sifa_id path string true "sifa_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid sifa_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /archived/subinterfaces/{sifa_id} [DELETE]
func (h *Handler) DeleteArchivedSubinterface(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "sifa_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid archived subinterface ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteArchivedSubinterface(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
