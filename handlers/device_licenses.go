package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aretaja/godevmandb"
	"github.com/go-chi/chi/v5"
)

// Count DeviceLicenses
// @Summary Count device_licenses
// @Description Count number of device_licenses
// @Tags devices
// @ID count-device_licenses
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/licenses/count [GET]
func (h *Handler) CountDeviceLicenses(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountDeviceLicenses(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List device_licenses
// @Summary List device_licenses
// @Description List device_licenses info
// @Tags devices
// @ID list-device_licenses
// @Param product_f query string false "url encoded SQL 'ILIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param descr_f query string false "url encoded SQL 'ILIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param condition_f query string false "url encoded SQL 'ILIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param installed_le query int false "SQL '<=' operator value"
// @Param installed_ge query int false "SQL '>=' operator value"
// @Param unlocked_le query int false "SQL '<=' operator value"
// @Param unlocked_ge query int false "SQL '>=' operator value"
// @Param tot_inst_le query int false "SQL '<=' operator value"
// @Param tot_inst_ge query int false "SQL '>=' operator value"
// @Param used_le query int false "SQL '<=' operator value"
// @Param used_ge query int false "SQL '>=' operator value"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} godevmandb.DeviceLicense
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/licenses [GET]
func (h *Handler) GetDeviceLicenses(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetDeviceLicensesParams{
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
	if v := r.FormValue("product_f"); v != "" {
		p.ProductF = &v
	}

	if v := r.FormValue("descr_f"); v != "" {
		p.DescrF = &v
	}

	if v := r.FormValue("condition_f"); v != "" {
		p.ConditionF = &v
	}

	if v := r.FormValue("installed_le"); v != "" {
		p.InstalledLe = &v
	}

	if v := r.FormValue("installed_ge"); v != "" {
		p.InstalledGe = &v
	}

	if v := r.FormValue("unlocked_le"); v != "" {
		p.UnlockedLe = &v
	}

	if v := r.FormValue("unlocked_ge"); v != "" {
		p.UnlockedGe = &v
	}

	if v := r.FormValue("tot_inst_le"); v != "" {
		p.TotInstLe = &v
	}

	if v := r.FormValue("tot_inst_ge"); v != "" {
		p.TotInstGe = &v
	}

	if v := r.FormValue("used_le"); v != "" {
		p.UsedLe = &v
	}

	if v := r.FormValue("used_ge"); v != "" {
		p.UsedGe = &v
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetDeviceLicenses(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Get DeviceLicense
// @Summary Get device_license
// @Description Get device_license info
// @Tags devices
// @ID get-device_license
// @Param lic_id path string true "lic_id"
// @Success 200 {object} godevmandb.DeviceLicense
// @Failure 400 {object} StatusResponse "Invalid lic_id"
// @Failure 404 {object} StatusResponse "DeviceLicense not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/licenses/{lic_id} [GET]
func (h *Handler) GetDeviceLicense(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "lic_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device license ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceLicense(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "License not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Create DeviceLicense
// @Summary Create device_license
// @Description Create device_license
// @Tags devices
// @ID create-device_license
// @Param Body body godevmandb.CreateDeviceLicenseParams true "JSON object of godevmandb.DeviceLicenseParams"
// @Success 201 {object} godevmandb.DeviceLicense
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/licenses [POST]
func (h *Handler) CreateDeviceLicense(w http.ResponseWriter, r *http.Request) {
	var p godevmandb.CreateDeviceLicenseParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	q := godevmandb.New(h.db)
	res, err := q.CreateDeviceLicense(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusCreated, res)
}

// Update DeviceLicense
// @Summary Update device_license
// @Description Update device_license
// @Tags devices
// @ID update-device_license
// @Param lic_id path string true "lic_id"
// @Param Body body godevmandb.UpdateDeviceLicenseParams true "JSON object of godevmandb.UpdateDeviceLicenseParams.<br />Ignored fields:<ul><li>lic_id</li></ul>"
// @Success 200 {object} godevmandb.DeviceLicense
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/licenses/{lic_id} [PUT]
func (h *Handler) UpdateDeviceLicense(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "lic_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device license ID")
		return
	}

	var p godevmandb.UpdateDeviceLicenseParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	p.LicID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateDeviceLicense(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Delete DeviceLicense
// @Summary Delete device_license
// @Description Delete device_license
// @Tags devices
// @ID delete-device_license
// @Param lic_id path string true "lic_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid lic_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/licenses/{lic_id} [DELETE]
func (h *Handler) DeleteDeviceLicense(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "lic_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device_license ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteDeviceLicense(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Foreign key
// Get DeviceLicense Device
// @Summary Get device_license device
// @Description Get device_license device info
// @Tags devices
// @ID get-device_license-device
// @Param lic_id path string true "lic_id"
// @Success 200 {object} device
// @Failure 400 {object} StatusResponse "Invalid lic_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/licenses/{lic_id}/device [GET]
func (h *Handler) GetDeviceLicenseDevice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "lic_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid device license ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceLicenseDevice(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := device{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}
