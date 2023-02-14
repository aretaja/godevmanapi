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
type device struct {
	UpdatedOn        time.Time `json:"updated_on"`
	CreatedOn        time.Time `json:"created_on"`
	Ip6Addr          *string   `json:"ip6_addr"`
	SysName          *string   `json:"sys_name"`
	Parent           *int64    `json:"parent"`
	Notes            *string   `json:"notes"`
	SnmpRoID         *int64    `json:"snmp_ro_id"`
	ExtModel         *string   `json:"ext_model"`
	SwVersion        *string   `json:"sw_version"`
	SysContact       *string   `json:"sys_contact"`
	SysLocation      *string   `json:"sys_location"`
	SnmpMainID       *int64    `json:"snmp_main_id"`
	Ip4Addr          *string   `json:"ip4_addr"`
	SiteID           *int64    `json:"site_id"`
	SysID            string    `json:"sys_id"`
	HostName         string    `json:"host_name"`
	Source           string    `json:"source"`
	DomID            int64     `json:"dom_id"`
	DevID            int64     `json:"dev_id"`
	Monitor          bool      `json:"monitor"`
	Unresponsive     bool      `json:"unresponsive"`
	ValidationFailed bool      `json:"validation_failed"`
	BackupFailed     bool      `json:"backup_failed"`
	TypeChanged      bool      `json:"type_changed"`
	Backup           bool      `json:"backup"`
	Graph            bool      `json:"graph"`
	Installed        bool      `json:"installed"`
}

// Import values from corresponding godevmandb struct
func (r *device) getValues(s godevmandb.Device) {
	r.DevID = s.DevID
	r.DomID = s.DomID
	r.SysID = s.SysID
	r.HostName = s.HostName
	r.Source = s.Source
	r.Installed = s.Installed
	r.Monitor = s.Monitor
	r.Graph = s.Graph
	r.Backup = s.Backup
	r.TypeChanged = s.TypeChanged
	r.BackupFailed = s.BackupFailed
	r.ValidationFailed = s.ValidationFailed
	r.Unresponsive = s.Unresponsive
	r.UpdatedOn = s.UpdatedOn
	r.CreatedOn = s.CreatedOn
	r.SiteID = nullInt64ToPtr(s.SiteID)
	r.SnmpMainID = nullInt64ToPtr(s.SnmpMainID)
	r.SnmpRoID = nullInt64ToPtr(s.SnmpRoID)
	r.Parent = nullInt64ToPtr(s.Parent)
	r.Ip4Addr = pgInetToPtr(s.Ip4Addr)
	r.Ip6Addr = pgInetToPtr(s.Ip6Addr)
	r.SysName = nullStringToPtr(s.SysName)
	r.SysLocation = nullStringToPtr(s.SysLocation)
	r.SysContact = nullStringToPtr(s.SysContact)
	r.SwVersion = nullStringToPtr(s.SwVersion)
	r.ExtModel = nullStringToPtr(s.ExtModel)
	r.Notes = nullStringToPtr(s.Notes)
}

// Return corresponding godevmandb create parameters
func (r *device) createParams() godevmandb.CreateDeviceParams {
	s := godevmandb.CreateDeviceParams{}

	s.DomID = r.DomID
	s.SysID = r.SysID
	s.HostName = r.HostName
	s.Source = r.Source
	s.Installed = r.Installed
	s.Monitor = r.Monitor
	s.Graph = r.Graph
	s.Backup = r.Backup
	s.TypeChanged = r.TypeChanged
	s.BackupFailed = r.BackupFailed
	s.ValidationFailed = r.ValidationFailed
	s.Unresponsive = r.Unresponsive
	s.SiteID = int64ToNullInt64(r.SiteID)
	s.SnmpMainID = int64ToNullInt64(r.SnmpMainID)
	s.SnmpRoID = int64ToNullInt64(r.SnmpRoID)
	s.Parent = int64ToNullInt64(r.Parent)
	s.Ip4Addr = strToPgInet(r.Ip4Addr)
	s.Ip6Addr = strToPgInet(r.Ip6Addr)
	s.SysName = strToNullString(r.SysName)
	s.SysLocation = strToNullString(r.SysLocation)
	s.SysContact = strToNullString(r.SysContact)
	s.SwVersion = strToNullString(r.SwVersion)
	s.ExtModel = strToNullString(r.ExtModel)
	s.Notes = strToNullString(r.Notes)

	return s
}

// Return corresponding godevmandb update parameters
func (r *device) updateParams() godevmandb.UpdateDeviceParams {
	s := godevmandb.UpdateDeviceParams{}

	s.DomID = r.DomID
	s.SysID = r.SysID
	s.HostName = r.HostName
	s.Source = r.Source
	s.Installed = r.Installed
	s.Monitor = r.Monitor
	s.Graph = r.Graph
	s.Backup = r.Backup
	s.TypeChanged = r.TypeChanged
	s.BackupFailed = r.BackupFailed
	s.ValidationFailed = r.ValidationFailed
	s.Unresponsive = r.Unresponsive
	s.SiteID = int64ToNullInt64(r.SiteID)
	s.SnmpMainID = int64ToNullInt64(r.SnmpMainID)
	s.SnmpRoID = int64ToNullInt64(r.SnmpRoID)
	s.Parent = int64ToNullInt64(r.Parent)
	s.Ip4Addr = strToPgInet(r.Ip4Addr)
	s.Ip6Addr = strToPgInet(r.Ip6Addr)
	s.SysName = strToNullString(r.SysName)
	s.SysLocation = strToNullString(r.SysLocation)
	s.SysContact = strToNullString(r.SysContact)
	s.SwVersion = strToNullString(r.SwVersion)
	s.ExtModel = strToNullString(r.ExtModel)
	s.Notes = strToNullString(r.Notes)

	return s
}

// Count Devices
// @Summary Count devices
// @Description Count number of devices
// @Tags devices
// @ID count-devices
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/count [GET]
func (h *Handler) CountDevices(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountDevices(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List devices
// @Summary List devices
// @Description List devices info
// @Tags devices
// @ID list-devices
// @Param sys_id_f query string false "url encoded SQL 'LIKE' operator pattern"
// @Param host_name_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param name_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param sw_version_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param notes_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param ip4_addr_f query string false "ip or containing net in CIDR notation"
// @Param ip6_addr_f query string false "ip or containing net in CIDR notation"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} device
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices [GET]
func (h *Handler) GetDevices(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetDevicesParams{
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
	v := r.FormValue("sys_id_f")
	if v != "" {
		p.SysIDF = v
	}

	// Hostname filter
	v = r.FormValue("host_name_f")
	if v != "" {
		p.HostNameF = v
	}

	// Software filter
	v = r.FormValue("sw_version_f")
	if v != "" {
		p.SwVersionF = strToNullString(&v)
	}

	// Notes filter
	v = r.FormValue("notes_f")
	if v != "" {
		p.NotesF = strToNullString(&v)
	}

	// Name filter
	v = r.FormValue("name_f")
	if v != "" {
		p.NameF = strToNullString(&v)
	}

	// Host IPv4 filter
	p.Ip4AddrF = strToPgInet(nil)
	v = r.FormValue("ip4_addr_f")
	if v != "" {
		p.Ip4AddrF = strToPgInet(&v)
	}

	// Host IPv6 filter
	p.Ip6AddrF = strToPgInet(nil)
	v = r.FormValue("ip6_addr_f")
	if v != "" {
		p.Ip6AddrF = strToPgInet(&v)
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetDevices(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := []device{}
	for _, s := range res {
		a := device{}
		a.getValues(s)
		out = append(out, a)
	}

	RespondJSON(w, r, http.StatusOK, out)
}

// Get Device
// @Summary Get device
// @Description Get device info
// @Tags devices
// @ID get-device
// @Param dev_id path string true "dev_id"
// @Success 200 {object} device
// @Failure 400 {object} StatusResponse "Invalid dev_id"
// @Failure 404 {object} StatusResponse "Device not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/{dev_id} [GET]
func (h *Handler) GetDevice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dev_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDevice(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "Device not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	out := device{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Create Device
// @Summary Create device
// @Description Create device
// @Tags devices
// @ID create-device
// @Param Body body device true "JSON object of device.<br />Ignored fields:<ul><li>dev_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 201 {object} device
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices [POST]
func (h *Handler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var pIn device
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Create parameters for new db record
	p := pIn.createParams()

	q := godevmandb.New(h.db)
	res, err := q.CreateDevice(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := device{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusCreated, out)
}

// Update Device
// @Summary Update device
// @Description Update device
// @Tags devices
// @ID update-device
// @Param dev_id path string true "dev_id"
// @Param Body body device true "JSON object of device.<br />Ignored fields:<ul><li>dev_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 200 {object} device
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/{dev_id} [PUT]
func (h *Handler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dev_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device ID")
		return
	}

	var pIn device
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
	res, err := q.UpdateDevice(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := device{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Delete Device
// @Summary Delete device
// @Description Delete device
// @Tags devices
// @ID delete-device
// @Param dev_id path string true "dev_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid dev_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/{dev_id} [DELETE]
func (h *Handler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dev_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteDevice(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
