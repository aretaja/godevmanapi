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
	r.SiteID = s.SiteID
	r.SnmpMainID = s.SnmpMainID
	r.SnmpRoID = s.SnmpRoID
	r.Parent = s.Parent
	r.Ip4Addr = pgInetToPtr(s.Ip4Addr)
	r.Ip6Addr = pgInetToPtr(s.Ip6Addr)
	r.SysName = s.SysName
	r.SysLocation = s.SysLocation
	r.SysContact = s.SysContact
	r.SwVersion = s.SwVersion
	r.ExtModel = s.ExtModel
	r.Notes = s.Notes
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
	s.SiteID = r.SiteID
	s.SnmpMainID = r.SnmpMainID
	s.SnmpRoID = r.SnmpRoID
	s.Parent = r.Parent
	s.Ip4Addr = strToPgInet(r.Ip4Addr)
	s.Ip6Addr = strToPgInet(r.Ip6Addr)
	s.SysName = r.SysName
	s.SysLocation = r.SysLocation
	s.SysContact = r.SysContact
	s.SwVersion = r.SwVersion
	s.ExtModel = r.ExtModel
	s.Notes = r.Notes

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
	s.SiteID = r.SiteID
	s.SnmpMainID = r.SnmpMainID
	s.SnmpRoID = r.SnmpRoID
	s.Parent = r.Parent
	s.Ip4Addr = strToPgInet(r.Ip4Addr)
	s.Ip6Addr = strToPgInet(r.Ip6Addr)
	s.SysName = r.SysName
	s.SysLocation = r.SysLocation
	s.SysContact = r.SysContact
	s.SwVersion = r.SwVersion
	s.ExtModel = r.ExtModel
	s.Notes = r.Notes

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
// @Failure 500 {object} StatusResponse "Failed DB transaction"
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
// @Param source_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param sys_name_f query string false "url encoded SQL 'ILIKE' operator pattern + special value 'isnull', 'isempty'"
// @Param sw_version_f query string false "url encoded SQL 'ILIKE' operator pattern + special value 'isnull', 'isempty'"
// @Param notes_f query string false "url encoded SQL 'ILIKE' operator pattern + special value 'isnull', 'isempty'"
// @Param ext_model_f query string false "url encoded SQL 'ILIKE' operator pattern + special value 'isnull', 'isempty'"
// @Param ip4_addr_f query string false "ip or containing net in CIDR notation"
// @Param ip6_addr_f query string false "ip or containing net in CIDR notation"
// @Param installed_f query bool false "values 'true', 'false'"
// @Param monitor_f query bool false "values 'true', 'false'"
// @Param graph_f query bool false "values 'true', 'false'"
// @Param backup_f query bool false "values 'true', 'false'"
// @Param type_changed_f query bool false "values 'true', 'false'"
// @Param backup_failed_f query bool false "values 'true', 'false'"
// @Param validation_failed_f query bool false "values 'true', 'false'"
// @Param unresponsive_f query bool false "values 'true', 'false'"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} device
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
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

	// Filters
	if v := r.FormValue("sys_id_f"); v != "" {
		p.SysIDF = v
	}

	if v := r.FormValue("host_name_f"); v != "" {
		p.HostNameF = v
	}

	if v := r.FormValue("source_f"); v != "" {
		p.SourceF = v
	}

	if v := r.FormValue("sw_version_f"); v != "" {
		p.SwVersionF = &v
	}

	if v := r.FormValue("notes_f"); v != "" {
		p.NotesF = &v
	}

	if v := r.FormValue("sys_name_f"); v != "" {
		p.SysNameF = &v
	}

	if v := r.FormValue("ext_model_f"); v != "" {
		p.ExtModelF = &v
	}

	p.Ip4AddrF = strToPgInet(nil)
	if v := r.FormValue("ip4_addr_f"); v != "" {
		p.Ip4AddrF = strToPgInet(&v)
	}

	p.Ip6AddrF = strToPgInet(nil)
	if v := r.FormValue("ip4_aip6_addr_fddr_f"); v != "" {
		p.Ip6AddrF = strToPgInet(&v)
	}

	if v := r.FormValue("installed_f"); v != "" {
		p.InstalledF = v
	}

	if v := r.FormValue("monitor_f"); v != "" {
		p.MonitorF = v
	}

	if v := r.FormValue("graph_f"); v != "" {
		p.GraphF = v
	}

	if v := r.FormValue("backup_f"); v != "" {
		p.BackupF = v
	}

	if v := r.FormValue("type_changed_f"); v != "" {
		p.TypeChangedF = v
	}

	if v := r.FormValue("backup_failed_f"); v != "" {
		p.BackupFailedF = v
	}

	if v := r.FormValue("validation_failed_f"); v != "" {
		p.ValidationFailedF = v
	}

	if v := r.FormValue("unresponsive_f"); v != "" {
		p.UnresponsiveF = v
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
// @Failure 500 {object} StatusResponse "Failed DB transaction"
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
// @Failure 500 {object} StatusResponse "Failed DB transaction"
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
// @Failure 500 {object} StatusResponse "Failed DB transaction"
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
// @Failure 500 {object} StatusResponse "Failed DB transaction"
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

// Foreign key
// Get Device DeviceDomain
// @Summary Get device device_domain
// @Description Get device device_domain info
// @Tags devices
// @ID get-device-device-domain
// @Param dev_id path string true "dev_id"
// @Success 200 {object} godevmandb.DeviceDomain
// @Failure 400 {object} StatusResponse "Invalid dev_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/{dev_id}/domain [GET]
func (h *Handler) GetDeviceDeviceDomain(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dev_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceDeviceDomain(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Foreign key
// Get Device DeviceType
// @Summary Get device device_type
// @Description Get device device_type info
// @Tags devices
// @ID get-device-device-type
// @Param dev_id path string true "dev_id"
// @Success 200 {object} godevmandb.DeviceType
// @Failure 400 {object} StatusResponse "Invalid dev_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/{dev_id}/type [GET]
func (h *Handler) GetDeviceDeviceType(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dev_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceDeviceType(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Foreign key
// Get Device Parent
// @Summary Get device parent
// @Description Get device parent info
// @Tags devices
// @ID get-device-parent
// @Param dev_id path string true "dev_id"
// @Success 200 {object} device
// @Failure 400 {object} StatusResponse "Invalid dev_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/{dev_id}/parent [GET]
func (h *Handler) GetDeviceParent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dev_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceParent(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := device{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Foreign key
// Get Device Site
// @Summary Get device site
// @Description Get device site info
// @Tags devices
// @ID get-device-site
// @Param dev_id path string true "dev_id"
// @Success 200 {object} godevmandb.Site
// @Failure 400 {object} StatusResponse "Invalid dev_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/{dev_id}/site [GET]
func (h *Handler) GetDeviceSite(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dev_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceSite(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Foreign key
// Get Device SnmpCredentialsMain
// @Summary Get device snmp_credentials_main
// @Description Get device snmp_credentials_main info
// @Tags devices
// @ID get-device-snmp-credentials-main
// @Param dev_id path string true "dev_id"
// @Success 200 {object} snmpCredential
// @Failure 400 {object} StatusResponse "Invalid dev_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/{dev_id}/snmp_credentials_main [GET]
func (h *Handler) GetDeviceSnmpCredentialsMain(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dev_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceSnmpCredentialsMain(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := snmpCredential{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Foreign key
// Get Device SnmpCredentialsRo
// @Summary Get device snmp_credentials_ro
// @Description Get device snmp_credentials_ro info
// @Tags devices
// @ID get-device-snmp-credentials-ro
// @Param dev_id path string true "dev_id"
// @Success 200 {object} snmpCredential
// @Failure 400 {object} StatusResponse "Invalid dev_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/{dev_id}/snmp_credentials_ro [GET]
func (h *Handler) GetDeviceSnmpCredentialsRo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dev_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceSnmpCredentialsRo(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := snmpCredential{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Relations
// List Device Childs
// @Summary List device childs
// @Description List device childs info
// @Tags devices
// @ID list-device-childs
// @Param dev_id path string true "dev_id"
// @Success 200 {array} device
// @Failure 400 {object} StatusResponse "Invalid dev_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/{dev_id}/childs [GET]
func (h *Handler) GetDeviceChilds(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dev_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceChilds(h.ctx, id)
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

// Relations
// List Device Credentials
// @Summary List device credentials
// @Description List device credentials info
// @Tags devices
// @ID list-device-credentials
// @Param dev_id path string true "dev_id"
// @Success 200 {array} godevmandb.DeviceCredential
// @Failure 400 {object} StatusResponse "Invalid dev_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/{dev_id}/credentials [GET]
func (h *Handler) GetDeviceDeviceCredentials(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dev_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceDeviceCredentials(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// Decrypt secret
	for i, s := range res {
		if s.EncSecret != "" {
			val, err := godevmandb.DecryptStrAes(s.EncSecret, salt)
			if err != nil {
				RespondError(w, r, http.StatusInternalServerError, err.Error())
				return
			}

			res[i].EncSecret = val
		}
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Relations
// List Device Extensions
// @Summary List device extensions
// @Description List device extensions info
// @Tags devices
// @ID list-device-extensions
// @Param dev_id path string true "dev_id"
// @Success 200 {array} godevmandb.DeviceExtension
// @Failure 400 {object} StatusResponse "Invalid dev_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/{dev_id}/extensions [GET]
func (h *Handler) GetDeviceDeviceExtensions(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dev_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceDeviceExtensions(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Relations
// List Device Licenses
// @Summary List device licenses
// @Description List device licenses info
// @Tags devices
// @ID list-device-licenses
// @Param dev_id path string true "dev_id"
// @Success 200 {array} godevmandb.DeviceLicense
// @Failure 400 {object} StatusResponse "Invalid dev_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/{dev_id}/licenses [GET]
func (h *Handler) GetDeviceDeviceLicenses(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dev_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceDeviceLicenses(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Relations
// List Device State
// @Summary List device state
// @Description List device state info
// @Tags devices
// @ID list-device-state
// @Param dev_id path string true "dev_id"
// @Success 200 {array} godevmandb.DeviceState
// @Failure 400 {object} StatusResponse "Invalid dev_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/{dev_id}/state [GET]
func (h *Handler) GetDeviceDeviceState(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dev_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceDeviceState(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Relations
// List Device Entities
// @Summary List device entities
// @Description List device entities info
// @Tags devices
// @ID list-device-entities
// @Param dev_id path string true "dev_id"
// @Success 200 {array} godevmandb.Entity
// @Failure 400 {object} StatusResponse "Invalid dev_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/{dev_id}/entities [GET]
func (h *Handler) GetDeviceEntities(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dev_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceEntities(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Relations
// List Device Interfaces
// @Summary List device interfaces
// @Description List device interfaces info
// @Tags devices
// @ID list-device-interfaces
// @Param dev_id path string true "dev_id"
// @Success 200 {array} iface
// @Failure 400 {object} StatusResponse "Invalid dev_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/{dev_id}/interfaces [GET]
func (h *Handler) GetDeviceInterfaces(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dev_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceInterfaces(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := []iface{}
	for _, s := range res {
		a := iface{}
		a.getValues(s)
		out = append(out, a)
	}

	RespondJSON(w, r, http.StatusOK, out)
}

// Relations
// List Device Ip Interfaces
// @Summary List device ip interfaces
// @Description List device ip interfaces info
// @Tags devices
// @ID list-device-ip-interfaces
// @Param dev_id path string true "dev_id"
// @Success 200 {array} ipInterface
// @Failure 400 {object} StatusResponse "Invalid dev_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/{dev_id}/ip_interfaces [GET]
func (h *Handler) GetDeviceIpInterfaces(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dev_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceIpInterfaces(h.ctx, id)
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

// Relations
// List Device Ospf Neighbors
// @Summary List device ospf nbrs
// @Description List device ospf nbrs info
// @Tags devices
// @ID list-device-ospf-nbrs
// @Param dev_id path string true "dev_id"
// @Success 200 {array} ospfNbr
// @Failure 400 {object} StatusResponse "Invalid dev_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/{dev_id}/ospf_nbrs [GET]
func (h *Handler) GetDeviceOspfNbrs(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dev_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceOspfNbrs(h.ctx, id)
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

// Relations
// List Device Peer Xconnects
// @Summary List device peer xconnects
// @Description List device peer xconnects info
// @Tags devices
// @ID list-device-peer-xconnects
// @Param dev_id path string true "dev_id"
// @Success 200 {array} xconnect
// @Failure 400 {object} StatusResponse "Invalid dev_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/{dev_id}/peer_xconnects [GET]
func (h *Handler) GetDevicePeerXconnects(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dev_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDevicePeerXconnects(h.ctx, &id)
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

// Relations
// List Device Radio Link Neighbors
// @Summary List device rl nbrs
// @Description List device rl nbrs info
// @Tags devices
// @ID list-device-rl-nbrs
// @Param dev_id path string true "dev_id"
// @Success 200 {array} godevmandb.RlNbr
// @Failure 400 {object} StatusResponse "Invalid dev_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/{dev_id}/rl_nbrs [GET]
func (h *Handler) GetDeviceRlNbrs(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dev_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceRlNbrs(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Relations
// List Device Vlans
// @Summary List device vlans
// @Description List device vlans info
// @Tags devices
// @ID list-device-vlans
// @Param dev_id path string true "dev_id"
// @Success 200 {array} godevmandb.Vlan
// @Failure 400 {object} StatusResponse "Invalid dev_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/{dev_id}/vlans [GET]
func (h *Handler) GetDeviceVlans(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dev_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceVlans(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Relations
// List Device Xconnects
// @Summary List device xconnects
// @Description List device xconnects info
// @Tags devices
// @ID list-device-xconnects
// @Param dev_id path string true "dev_id"
// @Success 200 {array} xconnect
// @Failure 400 {object} StatusResponse "Invalid dev_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /devices/{dev_id}/xconnects [GET]
func (h *Handler) GetDeviceXconnects(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dev_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceXconnects(h.ctx, id)
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
