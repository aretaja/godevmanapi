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
type deviceDomain struct {
	UpdatedOn time.Time `json:"updated_on"`
	CreatedOn time.Time `json:"created_on"`
	Descr     string    `json:"descr"`
	DomID     int64     `json:"dom_id"`
}

// Import values from corresponding godevmandb struct
func (r *deviceDomain) getValues(s godevmandb.DeviceDomain) {
	r.DomID = s.DomID
	r.Descr = s.Descr
	r.UpdatedOn = s.UpdatedOn
	r.CreatedOn = s.CreatedOn
}

// Return corresponding godevmandb create parameters
func (r *deviceDomain) createParams() string {
	return r.Descr
}

// Return corresponding godevmandb update parameters
func (r *deviceDomain) updateParams() godevmandb.UpdateDeviceDomainParams {
	s := godevmandb.UpdateDeviceDomainParams{}

	s.Descr = r.Descr

	return s
}

// Count DeviceDomains
// @Summary Count device_domains
// @Description Count number of device domains
// @Tags devices
// @ID count-device_domains
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/domains/count [GET]
func (h *Handler) CountDeviceDomains(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountDeviceDomains(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List device_domains
// @Summary List device_domains
// @Description List device domains info
// @Tags devices
// @ID list-device_domains
// @Param descr_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} deviceDomain
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/domains [GET]
func (h *Handler) GetDeviceDomains(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetDeviceDomainsParams{
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

	// Descr filter
	d := r.FormValue("descr_f")
	if d != "" {
		p.DescrF = d
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetDeviceDomains(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := []deviceDomain{}
	for _, s := range res {
		r := deviceDomain{}
		r.getValues(s)
		out = append(out, r)
	}

	RespondJSON(w, r, http.StatusOK, out)
}

// Get DeviceDomain
// @Summary Get device_domain
// @Description Get device domain info
// @Tags devices
// @ID get-device_domain
// @Param dom_id path string true "dom_id"
// @Success 200 {object} deviceDomain
// @Failure 400 {object} StatusResponse "Invalid dom_id"
// @Failure 404 {object} StatusResponse "Domain not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/domains/{dom_id} [GET]
func (h *Handler) GetDeviceDomain(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dom_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid domain ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceDomain(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "Domain not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	out := deviceDomain{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Create DeviceDomain
// @Summary Create device_domain
// @Description Create device domain
// @Tags devices
// @ID create-device_domain
// @Param Body body deviceDomain true "JSON object of deviceDomain<br />Ignored fields:<ul><li>dom_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 201 {object} deviceDomain
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/domains [POST]
func (h *Handler) CreateDeviceDomain(w http.ResponseWriter, r *http.Request) {
	var pIn deviceDomain
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Create parameters for new db record
	p := pIn.createParams()

	q := godevmandb.New(h.db)
	res, err := q.CreateDeviceDomain(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := deviceDomain{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusCreated, out)
}

// Update DeviceDomain
// @Summary Update device_domain
// @Description Update device domain
// @Tags devices
// @ID update-device_domain
// @Param dom_id path string true "dom_id"
// @Param Body body deviceDomain true "JSON object of deviceDomain.<br />Ignored fields:<ul><li>dom_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 200 {object} deviceDomain
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/domains/{dom_id} [PUT]
func (h *Handler) UpdateDeviceDomain(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dom_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid domain ID")
		return
	}

	var pIn deviceDomain
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Update parameters for new db record
	p := pIn.updateParams()
	p.DomID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateDeviceDomain(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := deviceDomain{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Delete DeviceDomain
// @Summary Delete device_domain
// @Description Delete device domain
// @Tags devices
// @ID delete-device_domain
// @Param dom_id path string true "dom_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid dom_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/domains/{dom_id} [DELETE]
func (h *Handler) DeleteDeviceDomain(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "dom_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid domain ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteDeviceDomain(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Relations
// List DeviceDomain Devices
// @Summary List device_domain devices
// @Description List device domain devices info
// @Tags devices
// @ID list-device_domain-devices
// @Param dom_id path string true "dom_id"
// @Success 200 {array} device
// @Failure 400 {object} StatusResponse "Invalid dom_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/domains/{dom_id}/devices [GET]
// func (h *Handler) GetDeviceDomainDevices(w http.ResponseWriter, r *http.Request) {
// 	id, err := strconv.ParseInt(chi.URLParam(r, "dom_id"), 10, 64)
// 	if err != nil {
// 		RespondError(w, r, http.StatusBadRequest, "Invalid domain ID")
// 		return
// 	}

// 	q := godevmandb.New(h.db)
// 	res, err := q.GetDeviceDomainDevices(h.ctx, id)
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
