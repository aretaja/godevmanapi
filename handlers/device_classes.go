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
type deviceClass struct {
	UpdatedOn time.Time `json:"updated_on"`
	CreatedOn time.Time `json:"created_on"`
	Descr     string    `json:"descr"`
	ClassID   int64     `json:"class_id"`
}

// Import values from corresponding godevmandb struct
func (r *deviceClass) getValues(s godevmandb.DeviceClass) {
	r.ClassID = s.ClassID
	r.Descr = s.Descr
	r.UpdatedOn = s.UpdatedOn
	r.CreatedOn = s.CreatedOn
}

// Return corresponding godevmandb create parameters
func (r *deviceClass) createParams() string {
	return r.Descr
}

// Return corresponding godevmandb update parameters
func (r *deviceClass) updateParams() godevmandb.UpdateDeviceClassParams {
	s := godevmandb.UpdateDeviceClassParams{}

	s.Descr = r.Descr

	return s
}

// Count DeviceClasses
// @Summary Count device_classes
// @Description Count number of device classes
// @Tags devices
// @ID count-device_classes
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/classes/count [GET]
func (h *Handler) CountDeviceClasses(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountDeviceClasses(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List device_classes
// @Summary List device_classes
// @Description List device classes info
// @Tags devices
// @ID list-device_classes
// @Param descr_f query string false "url encoded SQL 'LIKE' operator pattern"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} deviceClass
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/classes [GET]
func (h *Handler) GetDeviceClasses(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetDeviceClassesParams{
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
	res, err := q.GetDeviceClasses(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := []deviceClass{}
	for _, s := range res {
		r := deviceClass{}
		r.getValues(s)
		out = append(out, r)
	}

	RespondJSON(w, r, http.StatusOK, out)
}

// Get DeviceClass
// @Summary Get device_class
// @Description Get device class info
// @Tags devices
// @ID get-device_class
// @Param class_id path string true "class_id"
// @Success 200 {object} deviceClass
// @Failure 400 {object} StatusResponse "Invalid class_id"
// @Failure 404 {object} StatusResponse "Class not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/classes/{class_id} [GET]
func (h *Handler) GetDeviceClass(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "class_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid class ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceClass(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "Class not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	out := deviceClass{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Create DeviceClass
// @Summary Create device_class
// @Description Create device class
// @Tags devices
// @ID create-device_class
// @Param Body body deviceClass true "JSON object of deviceClass<br />Ignored fields:<ul><li>class_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 201 {object} deviceClass
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/classes [POST]
func (h *Handler) CreateDeviceClass(w http.ResponseWriter, r *http.Request) {
	var pIn deviceClass
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Create parameters for new db record
	p := pIn.createParams()

	q := godevmandb.New(h.db)
	res, err := q.CreateDeviceClass(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := deviceClass{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusCreated, out)
}

// Update DeviceClass
// @Summary Update device_class
// @Description Update device class
// @Tags devices
// @ID update-device_class
// @Param class_id path string true "class_id"
// @Param Body body deviceClass true "JSON object of deviceClass.<br />Ignored fields:<ul><li>class_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 200 {object} deviceClass
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/classes/{class_id} [PUT]
func (h *Handler) UpdateDeviceClass(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "class_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid class ID")
		return
	}

	var pIn deviceClass
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Update parameters for new db record
	p := pIn.updateParams()
	p.ClassID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateDeviceClass(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := deviceClass{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Delete DeviceClass
// @Summary Delete device_class
// @Description Delete device class
// @Tags devices
// @ID delete-device_class
// @Param class_id path string true "class_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid class_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/classes/{class_id} [DELETE]
func (h *Handler) DeleteDeviceClass(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "class_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid class ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteDeviceClass(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Relations
// List Device Class Device Types
// @Summary List device_class device_types
// @Description List device class device types info
// @Tags devices
// @ID list-device_class-device_types
// @Param class_id path string true "class_id"
// @Success 200 {array} deviceType
// @Failure 400 {object} StatusResponse "Invalid class_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/classes/{class_id}/types [GET]
func (h *Handler) GetDeviceClassTypes(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "class_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid class ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceClassDeviceTypes(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := []deviceType{}
	for _, s := range res {
		a := deviceType{}
		a.getValues(s)
		out = append(out, a)
	}

	RespondJSON(w, r, http.StatusOK, out)
}
