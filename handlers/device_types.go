package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"github.com/aretaja/godevmandb"
	"github.com/go-chi/chi/v5"
)

// JSON friendly local type to use in web api. Replaces sql.Null*/pgtype fields
type deviceType struct {
	UpdatedOn    time.Time `json:"updated_on"`
	CreatedOn    time.Time `json:"created_on"`
	SysID        string    `json:"sys_id"`
	Manufacturer string    `json:"manufacturer"`
	Model        string    `json:"model"`
	ClassID      int64     `json:"class_id"`
	SnmpVer      int16     `json:"snmp_ver"`
	Hc           bool      `json:"hc"`
}

// Import values from corresponding godevmandb struct
func (r *deviceType) getValues(s godevmandb.DeviceType) {
	r.SysID = s.SysID
	r.ClassID = s.ClassID
	r.Manufacturer = s.Manufacturer
	r.Model = s.Model
	r.Hc = s.Hc
	r.SnmpVer = s.SnmpVer
	r.UpdatedOn = s.UpdatedOn
	r.CreatedOn = s.CreatedOn
}

// Return corresponding godevmandb create parameters
func (r *deviceType) createParams() godevmandb.CreateDeviceTypeParams {
	s := godevmandb.CreateDeviceTypeParams{}

	s.SysID = r.SysID
	s.ClassID = r.ClassID
	s.Manufacturer = r.Manufacturer
	s.Model = r.Model
	s.Hc = r.Hc
	s.SnmpVer = r.SnmpVer

	return s
}

// Return corresponding godevmandb update parameters
func (r *deviceType) updateParams() godevmandb.UpdateDeviceTypeParams {
	s := godevmandb.UpdateDeviceTypeParams{}

	s.ClassID = r.ClassID
	s.Manufacturer = r.Manufacturer
	s.Model = r.Model
	s.Hc = r.Hc
	s.SnmpVer = r.SnmpVer

	return s
}

// Count DeviceTypes
// @Summary Count device_types
// @Description Count number of device types
// @Tags devices
// @ID count-device_types
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/types/count [GET]
func (h *Handler) CountDeviceTypes(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountDeviceTypes(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List device_types
// @Summary List device_types
// @Description List device types info
// @Tags devices
// @ID list-device_types
// @Param sys_id_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param manufacturer_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param model_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param limit query int false "min: 1; max: 100; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} deviceType
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/types [GET]
func (h *Handler) GetDeviceTypes(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetDeviceTypesParams{
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
	str := r.FormValue("sys_id_f")
	if str != "" {
		p.SysIDF = str
	}

	// Manufacturer filter
	str = r.FormValue("manufacturer_f")
	if str != "" {
		p.ManufacturerF = str
	}

	// Model filter
	str = r.FormValue("model_f")
	if str != "" {
		p.ModelF = str
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetDeviceTypes(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := []deviceType{}
	for _, s := range res {
		r := deviceType{}
		r.getValues(s)
		out = append(out, r)
	}

	RespondJSON(w, r, http.StatusOK, out)
}

// Get DeviceType
// @Summary Get device_type
// @Description Get device type info
// @Tags devices
// @ID get-device_type
// @Param sys_id path string true "sys_id"
// @Success 200 {object} deviceType
// @Failure 400 {object} StatusResponse "Invalid sys_id"
// @Failure 404 {object} StatusResponse "Domain not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/types/{sys_id} [GET]
func (h *Handler) GetDeviceType(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "sys_id")

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceType(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "Domain not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	out := deviceType{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Create DeviceType
// @Summary Create device_type
// @Description Create device type
// @Tags devices
// @ID create-device_type
// @Param Body body deviceType true "JSON object of deviceType<br />Ignored fields:<ul><li>updated_on</li><li>created_on</li></ul>sys_id must match ^[\w-\.]+$ regex"
// @Success 201 {object} deviceType
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/types [POST]
func (h *Handler) CreateDeviceType(w http.ResponseWriter, r *http.Request) {
	var pIn deviceType
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Validate sys_id
	pattern := regexp.MustCompile(`^[\w-\.]+$`)
	if !pattern.MatchString(pIn.SysID) {
		RespondError(w, r, http.StatusBadRequest, "Invalid sys_id")
		return
	}

	// Create parameters for new db record
	p := pIn.createParams()

	q := godevmandb.New(h.db)
	res, err := q.CreateDeviceType(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := deviceType{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusCreated, out)
}

// Update DeviceType
// @Summary Update device_type
// @Description Update device type
// @Tags devices
// @ID update-device_type
// @Param sys_id path string true "sys_id"
// @Param Body body deviceType true "JSON object of deviceType.<br />Ignored fields:<ul><li>sys_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 200 {object} deviceType
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/types/{sys_id} [PUT]
func (h *Handler) UpdateDeviceType(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "sys_id")

	var pIn deviceType
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Update parameters for new db record
	p := pIn.updateParams()
	p.SysID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateDeviceType(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := deviceType{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Delete DeviceType
// @Summary Delete device_type
// @Description Delete device type
// @Tags devices
// @ID delete-device_type
// @Param sys_id path string true "sys_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid sys_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/types/{sys_id} [DELETE]
func (h *Handler) DeleteDeviceType(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "sys_id")

	q := godevmandb.New(h.db)
	err := q.DeleteDeviceType(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Foreign key
// Get Device Type Class
// @Summary Get device type class
// @Description Get device type class info
// @Tags devices
// @ID get-device-type-class
// @Param sys_id path string true "sys_id"
// @Success 200 {object} deviceClass
// @Failure 400 {object} StatusResponse "Invalid sys_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/types/{sys_id}/class [GET]
func (h *Handler) GetDeviceTypeClass(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "sys_id")

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceTypeDeviceClass(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := deviceClass{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Relations
// List DeviceType Devices
// @Summary List device_type devices
// @Description List device type devices info
// @Tags devices
// @ID list-device_type-devices
// @Param sys_id path string true "sys_id"
// @Success 200 {array} device
// @Failure 400 {object} StatusResponse "Invalid sys_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/types/{sys_id}/devices [GET]
// func (h *Handler) GetDeviceTypeDevices(w http.ResponseWriter, r *http.Request) {
// 	id, err := strconv.ParseInt(chi.URLParam(r, "sys_id"), 10, 64)
// 	if err != nil {
// 		RespondError(w, r, http.StatusBadRequest, "Invalid type ID")
// 		return
// 	}

// 	q := godevmandb.New(h.db)
// 	res, err := q.GetDeviceTypeDevices(h.ctx, id)
// 	if err != nil {
// 		RespondError(w, r, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	out := []device{}
// 	for _, s := range res {
// 		a := device{}
// 		a.getValues(s)
// 		out = append(out, a)
// 	}

// 	RespondJSON(w, r, http.StatusOK, out)
// }
