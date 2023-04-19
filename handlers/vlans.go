package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aretaja/godevmandb"
	"github.com/go-chi/chi/v5"
)

// Count Vlans
// @Summary Count vlans
// @Description Count number of vlans
// @Tags devices
// @ID count-vlans
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/vlans/count [GET]
func (h *Handler) CountVlans(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountVlans(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List vlans
// @Summary List vlans
// @Description List vlans info
// @Tags devices
// @ID list-vlans
// @Param vlan_f query string false "url encoded SQL 'LIKE' operator pattern"
// @Param descr_f query string false "url encoded SQL 'ILIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} godevmandb.Vlan
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/vlans [GET]
func (h *Handler) GetVlans(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetVlansParams{
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
	if v := r.FormValue("vlan_f"); v != "" {
		p.VlanF = v
	}

	if v := r.FormValue("descr_f"); v != "" {
		p.DescrF = &v
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetVlans(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Get Vlan
// @Summary Get vlan
// @Description Get vlan info
// @Tags devices
// @ID get-vlan
// @Param v_id path string true "v_id"
// @Success 200 {object} godevmandb.Vlan
// @Failure 400 {object} StatusResponse "Invalid v_id"
// @Failure 404 {object} StatusResponse "Vlan not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/vlans/{v_id} [GET]
func (h *Handler) GetVlan(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "v_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid vlan ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetVlan(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "Neighbor not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Create Vlan
// @Summary Create vlan
// @Description Create vlan
// @Tags devices
// @ID create-vlan
// @Param Body body godevmandb.CreateVlanParams true "JSON object of godevmandb.VlanParams"
// @Success 201 {object} godevmandb.Vlan
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/vlans [POST]
func (h *Handler) CreateVlan(w http.ResponseWriter, r *http.Request) {
	var p godevmandb.CreateVlanParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	q := godevmandb.New(h.db)
	res, err := q.CreateVlan(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusCreated, res)
}

// Update Vlan
// @Summary Update vlan
// @Description Update vlan
// @Tags devices
// @ID update-vlan
// @Param v_id path string true "v_id"
// @Param Body body godevmandb.UpdateVlanParams true "JSON object of godevmandb.UpdateVlanParams.<br />Ignored fields:<ul><li>v_id</li></ul>"
// @Success 200 {object} godevmandb.Vlan
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/vlans/{v_id} [PUT]
func (h *Handler) UpdateVlan(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "v_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid vlan ID")
		return
	}

	var p godevmandb.UpdateVlanParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	p.VID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateVlan(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Delete Vlan
// @Summary Delete vlan
// @Description Delete vlan
// @Tags devices
// @ID delete-vlan
// @Param v_id path string true "v_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid v_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/vlans/{v_id} [DELETE]
func (h *Handler) DeleteVlan(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "v_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid vlan ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteVlan(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Foreign key
// Get Vlan Device
// @Summary Get vlan device
// @Description Get vlan device info
// @Tags devices
// @ID get-vlan-device
// @Param v_id path string true "v_id"
// @Success 200 {object} device
// @Failure 400 {object} StatusResponse "Invalid v_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/vlans/{v_id}/device [GET]
func (h *Handler) GetVlanDevice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "v_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid vlan ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetVlanDevice(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := device{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}
