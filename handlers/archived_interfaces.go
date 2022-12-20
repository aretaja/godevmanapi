package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/aretaja/godevmandb"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgtype"
	"github.com/jinzhu/copier"
)

// Local type to use with copier. Used for sql Null* field replacement in json response
type ArchivedInterface struct {
	UpdatedOn             time.Time `json:"updated_on"`
	CreatedOn             time.Time `json:"created_on"`
	HostIp6Val            *string   `json:"host_ip6"`
	CiscoOptPowerIndexVal *int64    `json:"cisco_opt_power_index"`
	HostIp4Val            *string   `json:"host_ip4"`
	AliasVal              *string   `json:"alias"`
	TypeEnumVal           *int64    `json:"type_enum"`
	MacVal                *string   `json:"mac"`
	OtnIfIDVal            *int64    `json:"otn_if_id"`
	IfindexVal            *int64    `json:"ifindex"`
	Hostname              string    `json:"hostname"`
	Manufacturer          string    `json:"manufacturer"`
	Model                 string    `json:"model"`
	Descr                 string    `json:"descr"`
	IfaID                 int64     `json:"ifa_id"`
}

func (a *ArchivedInterface) Ifindex(m sql.NullInt64) {
	if m.Valid {
		if v, err := m.Value(); err == nil {
			if res, ok := v.(int64); ok {
				a.IfindexVal = &res
			}
		}
	}
}

func (a *ArchivedInterface) OtnIfID(m sql.NullInt64) {
	if m.Valid {
		if v, err := m.Value(); err == nil {
			if res, ok := v.(int64); ok {
				a.OtnIfIDVal = &res
			}
		}
	}
}

func (a *ArchivedInterface) CiscoOptPowerIndex(m sql.NullInt64) {
	if m.Valid {
		if v, err := m.Value(); err == nil {
			if res, ok := v.(int64); ok {
				a.CiscoOptPowerIndexVal = &res
			}
		}
	}
}

func (a *ArchivedInterface) Alias(m sql.NullString) {
	if m.Valid {
		if v, err := m.Value(); err == nil {
			if res, ok := v.(string); ok {
				a.AliasVal = &res
			}
		}
	}
}

func (a *ArchivedInterface) TypeEnum(m sql.NullInt16) {
	if m.Valid {
		if v, err := m.Value(); err == nil {
			if res, ok := v.(int64); ok {
				a.TypeEnumVal = &res
			}
		}
	}
}

func (a *ArchivedInterface) Mac(m pgtype.Macaddr) {
	if m.Status == 2 {
		if v, err := m.Value(); err == nil {
			res := fmt.Sprintf("%s", v)
			a.MacVal = &res
		}
	}
}

func (a *ArchivedInterface) HostIp4(m pgtype.Inet) {
	if m.Status == 2 {
		if v, err := m.Value(); err == nil {
			res := fmt.Sprintf("%s", v)
			a.HostIp4Val = &res
		}
	}
}

func (a *ArchivedInterface) HostIp6(m pgtype.Inet) {
	if m.Status == 2 {
		if v, err := m.Value(); err == nil {
			res := fmt.Sprintf("%s", v)
			a.HostIp6Val = &res
		}
	}
}

// Count ArchivedInterfaces
// @Summary Count archived_interfaces
// @Description Count number of archived interfaces
// @Tags archived
// @ID count-archived_interfaces
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
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
// @Success 200 {array} ArchivedInterface
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /archived/interfaces [GET]
func (h *Handler) GetArchivedInterfaces(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetArchivedInterfacesParams{
		LimitQ:  1000,
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
	v = r.FormValue("host_ip4_f")
	if v != "" {
		p.HostIp4F = strToPgInet(&v)
	}

	// Host IPv6 filter
	v = r.FormValue("host_ip6_f")
	if v != "" {
		p.HostIp6F = strToPgInet(&v)
	}

	// MAC filter
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
	res, err := q.GetArchivedInterfaces(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := []ArchivedInterface{}
	copier.Copy(&out, &res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Get ArchivedInterface
// @Summary Get archived_interface
// @Description Get archived interface info
// @Tags archived
// @ID get-archived_interface
// @Param ifa_id path string true "ifa_id"
// @Success 200 {object} ArchivedInterface
// @Failure 400 {object} StatusResponse "Invalid ifa_id"
// @Failure 404 {object} StatusResponse "Archived interface not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
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

	out := ArchivedInterface{}
	copier.Copy(&out, &res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Create ArchivedInterface
// @Summary Create archived_interface
// @Description Create archived interface
// @Tags archived
// @ID create-archived_interface
// @Param Body body godevmandb.CreateArchivedInterfaceParams true "JSON object of CreateArchivedInterfaceParams"
// @Success 201 {object} ArchivedInterface
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /archived/interfaces [POST]
func (h *Handler) CreateArchivedInterface(w http.ResponseWriter, r *http.Request) {
	var p godevmandb.CreateArchivedInterfaceParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	q := godevmandb.New(h.db)
	res, err := q.CreateArchivedInterface(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := ArchivedInterface{}
	copier.Copy(&out, &res)

	RespondJSON(w, r, http.StatusCreated, out)
}

// Update ArchivedInterface
// @Summary Update archived_interface
// @Description Update archived interface
// @Tags archived
// @ID update-archived_interface
// @Param ifa_id path string true "ifa_id"
// @Param Body body godevmandb.UpdateArchivedInterfaceParams true "JSON object of UpdateArchivedInterfaceParams"
// @Success 200 {object} ArchivedInterface
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /archived/interfaces/{ifa_id} [PUT]
func (h *Handler) UpdateArchivedInterface(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ifa_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid archived interface ID")
		return
	}

	var p godevmandb.UpdateArchivedInterfaceParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	p.IfaID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateArchivedInterface(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := ArchivedInterface{}
	copier.Copy(&out, &res)

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
// @Failure 500 {object} StatusResponse "Failde DB transaction"
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
