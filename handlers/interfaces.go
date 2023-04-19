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
type iface struct {
	UpdatedOn  time.Time `json:"updated_on"`
	CreatedOn  time.Time `json:"created_on"`
	Adm        *int16    `json:"adm"`
	Mac        *string   `json:"mac"`
	ConID      *int64    `json:"con_id"`
	EntID      *int64    `json:"ent_id"`
	Ifindex    *int64    `json:"ifindex"`
	OtnIfID    *int64    `json:"otn_if_id"`
	Alias      *string   `json:"alias"`
	Oper       *int16    `json:"oper"`
	Parent     *int64    `json:"parent"`
	Speed      *int64    `json:"speed"`
	Minspeed   *int64    `json:"minspeed"`
	TypeEnum   *int16    `json:"type_enum"`
	Descr      string    `json:"descr"`
	IfID       int64     `json:"if_id"`
	DevID      int64     `json:"dev_id"`
	Monstatus  int16     `json:"monstatus"`
	Monerrors  int16     `json:"monerrors"`
	Monload    int16     `json:"monload"`
	Montraffic int16     `json:"montraffic"`
}

// Import values from corresponding godevmandb struct
func (r *iface) getValues(s godevmandb.Interface) {
	r.IfID = s.IfID
	r.ConID = s.ConID
	r.OtnIfID = s.OtnIfID
	r.DevID = s.DevID
	r.EntID = s.EntID
	r.Ifindex = s.Ifindex
	r.Descr = s.Descr
	r.Alias = s.Alias
	r.Oper = s.Oper
	r.Adm = s.Adm
	r.Speed = s.Speed
	r.Minspeed = s.Minspeed
	r.TypeEnum = s.TypeEnum
	r.Monstatus = s.Monstatus
	r.Monerrors = s.Monerrors
	r.Monload = s.Monload
	r.Montraffic = s.Montraffic
	r.UpdatedOn = s.UpdatedOn
	r.CreatedOn = s.CreatedOn
	r.Mac = pgMacaddrToPtr(s.Mac)
}

// Return corresponding godevmandb create parameters
func (r *iface) createParams() godevmandb.CreateInterfaceParams {
	s := godevmandb.CreateInterfaceParams{}

	s.ConID = r.ConID
	s.OtnIfID = r.OtnIfID
	s.DevID = r.DevID
	s.EntID = r.EntID
	s.Ifindex = r.Ifindex
	s.Descr = r.Descr
	s.Alias = r.Alias
	s.Oper = r.Oper
	s.Adm = r.Adm
	s.Speed = r.Speed
	s.Minspeed = r.Minspeed
	s.TypeEnum = r.TypeEnum
	s.Monstatus = r.Monstatus
	s.Monerrors = r.Monerrors
	s.Monload = r.Monload
	s.Montraffic = r.Montraffic
	s.Mac = strToPgMacaddr(r.Mac)

	return s
}

// Return corresponding godevmandb update parameters
func (r *iface) updateParams() godevmandb.UpdateInterfaceParams {
	s := godevmandb.UpdateInterfaceParams{}

	s.ConID = r.ConID
	s.OtnIfID = r.OtnIfID
	s.DevID = r.DevID
	s.EntID = r.EntID
	s.Ifindex = r.Ifindex
	s.Descr = r.Descr
	s.Alias = r.Alias
	s.Oper = r.Oper
	s.Adm = r.Adm
	s.Speed = r.Speed
	s.Minspeed = r.Minspeed
	s.TypeEnum = r.TypeEnum
	s.Monstatus = r.Monstatus
	s.Monerrors = r.Monerrors
	s.Monload = r.Monload
	s.Montraffic = r.Montraffic
	s.Mac = strToPgMacaddr(r.Mac)

	return s
}

// Count Interfaces
// @Summary Count interfaces
// @Description Count number of interfaces
// @Tags interfaces
// @ID count-interfaces
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/count [GET]
func (h *Handler) CountInterfaces(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountInterfaces(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List interfaces
// @Summary List interfaces
// @Description List interfaces info
// @Tags interfaces
// @ID list-interfaces
// @Param ifindex_f query string false "url encoded SQL 'LIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param descr_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param alias_f query string false "url encoded SQL 'ILIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param oper_f query string false "url encoded SQL 'LIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param adm_f query string false "url encoded SQL 'LIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param speed_f query string false "url encoded SQL 'LIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param minspeed_f query string false "url encoded SQL 'LIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param type_enum_f query string false "url encoded SQL 'LIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param mac_f query string false "SQL '=' operator value (MAC address)"
// @Param monstatus_f query bool false "values 'true', 'false'"
// @Param monerrors_f query bool false "values 'true', 'false'"
// @Param monload_f query bool false "values 'true', 'false'"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} iface
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces [GET]
func (h *Handler) GetInterfaces(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetInterfacesParams{
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

	if v := r.FormValue("alias_f"); v != "" {
		p.AliasF = &v
	}

	if v := r.FormValue("oper_f"); v != "" {
		p.OperF = &v
	}

	if v := r.FormValue("adm_f"); v != "" {
		p.AdmF = &v
	}

	if v := r.FormValue("speed_f"); v != "" {
		p.SpeedF = &v
	}

	if v := r.FormValue("minspeed_f"); v != "" {
		p.MinspeedF = &v
	}

	if v := r.FormValue("type_enum_f"); v != "" {
		p.TypeEnumF = &v
	}

	p.MacF = strToPgMacaddr(nil)
	if v := r.FormValue("mac_f"); v != "" {
		p.MacF = strToPgMacaddr(&v)
	}

	if v := r.FormValue("monstatus_f"); v != "" {
		p.MonstatusF = v
	}

	if v := r.FormValue("monerrors_f"); v != "" {
		p.MonerrorsF = v
	}

	if v := r.FormValue("monload_f"); v != "" {
		p.MonloadF = v
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetInterfaces(h.ctx, p)
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

// Get Interface
// @Summary Get interface
// @Description Get interface info
// @Tags interfaces
// @ID get-interface
// @Param if_id path string true "if_id"
// @Success 200 {object} iface
// @Failure 400 {object} StatusResponse "Invalid if_id"
// @Failure 404 {object} StatusResponse "Interface not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/{if_id} [GET]
func (h *Handler) GetInterface(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "if_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid interface ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetInterface(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "Interface not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	out := iface{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Create Interface
// @Summary Create interface
// @Description Create interface
// @Tags interfaces
// @ID create-interface
// @Param Body body iface true "JSON object of iface.<br />Ignored fields:<ul><li>if_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 201 {object} iface
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces [POST]
func (h *Handler) CreateInterface(w http.ResponseWriter, r *http.Request) {
	var pIn iface
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Create parameters for new db record
	p := pIn.createParams()

	q := godevmandb.New(h.db)
	res, err := q.CreateInterface(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := iface{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusCreated, out)
}

// Update Interface
// @Summary Update interface
// @Description Update interface
// @Tags interfaces
// @ID update-interface
// @Param if_id path string true "if_id"
// @Param Body body iface true "JSON object of iface.<br />Ignored fields:<ul><li>if_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 200 {object} iface
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/{if_id} [PUT]
func (h *Handler) UpdateInterface(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "if_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid interface ID")
		return
	}

	var pIn iface
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Update parameters for new db record
	p := pIn.updateParams()
	p.IfID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateInterface(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := iface{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Delete Interface
// @Summary Delete interface
// @Description Delete interface
// @Tags interfaces
// @ID delete-interface
// @Param if_id path string true "if_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid if_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/{if_id} [DELETE]
func (h *Handler) DeleteInterface(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "if_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid interface ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteInterface(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Foreign key
// Get Interface Connection
// @Summary Get interface connection
// @Description Get interface connection info
// @Tags interfaces
// @ID get-interface-connection
// @Param if_id path string true "if_id"
// @Success 200 {object} godevmandb.Connection
// @Failure 400 {object} StatusResponse "Invalid if_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/{if_id}/connection [GET]
func (h *Handler) GetInterfaceConnection(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "if_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid interface ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetInterfaceConnection(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Foreign key
// Get Interface Parent
// @Summary Get interface parent
// @Description Get interface parent info
// @Tags interfaces
// @ID get-interface-parent
// @Param if_id path string true "if_id"
// @Success 200 {object} iface
// @Failure 400 {object} StatusResponse "Invalid if_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/{if_id}/parent [GET]
func (h *Handler) GetInterfaceParent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "if_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid interface ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetInterfaceParent(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := iface{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Foreign key
// Get Interface Otn Interface
// @Summary Get interface related otn interface
// @Description Get interface otn_if info
// @Tags interfaces
// @ID get-interface-otn_if
// @Param if_id path string true "if_id"
// @Success 200 {object} iface
// @Failure 400 {object} StatusResponse "Invalid if_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/{if_id}/otn_if [GET]
func (h *Handler) GetInterfaceOtnIf(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "if_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid interface ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetInterfaceOtnIf(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := iface{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Foreign key
// Get Interface Device
// @Summary Get interface device
// @Description Get interface device info
// @Tags interfaces
// @ID get-interface-device
// @Param if_id path string true "if_id"
// @Success 200 {object} device
// @Failure 400 {object} StatusResponse "Invalid if_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/{if_id}/device [GET]
func (h *Handler) GetInterfaceDevice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "if_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid interface ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetInterfaceDevice(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := device{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Foreign key
// Get Interface Entity
// @Summary Get interface entity
// @Description Get interface entity info
// @Tags interfaces
// @ID get-interface-entity
// @Param if_id path string true "if_id"
// @Success 200 {object} godevmandb.Entity
// @Failure 400 {object} StatusResponse "Invalid if_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/{if_id}/entity [GET]
func (h *Handler) GetInterfaceEntity(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "if_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid interface ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetInterfaceEntity(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Relations
// List Interface Childs
// @Summary List interface childs
// @Description List interface childs info
// @Tags interfaces
// @ID list-interface-childs
// @Param if_id path string true "if_id"
// @Success 200 {array} iface
// @Failure 400 {object} StatusResponse "Invalid if_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/{if_id}/childs [GET]
func (h *Handler) GetInterfaceChilds(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "if_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid interface ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetInterfaceChilds(h.ctx, id)
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
// List Interface BwStats
// @Summary List interface bw_stats
// @Description List interface int_bw_stats info
// @Tags interfaces
// @ID list-interface-int_bw_stats
// @Param if_id path string true "if_id"
// @Success 200 {array} godevmandb.IntBwStat
// @Failure 400 {object} StatusResponse "Invalid if_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/{if_id}/bw_stats [GET]
func (h *Handler) GetInterfaceIntBwStats(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "if_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid interface ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetInterfaceIntBwStats(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Relations
// List Lower Related Interfaces
// @Summary List lower related interfaces
// @Description List lower related interfaces info
// @Tags interfaces
// @ID list-interface-lower-related
// @Param if_id path string true "if_id"
// @Success 200 {array} iface
// @Failure 400 {object} StatusResponse "Invalid if_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/{if_id}/related_lower [GET]
func (h *Handler) GetInterfaceInterfaceRelationsHigherFor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "if_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid interface ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetInterfaceInterfaceRelationsHigherFor(h.ctx, id)
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
// List Higher Related Interfaces
// @Summary List higher related interfaces
// @Description List higher related interfaces info
// @Tags interfaces
// @ID list-interface-higher-related
// @Param if_id path string true "if_id"
// @Success 200 {array} iface
// @Failure 400 {object} StatusResponse "Invalid if_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/{if_id}/related_higher [GET]
func (h *Handler) GetInterfaceInterfaceRelationsLowerFor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "if_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid interface ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetInterfaceInterfaceRelationsLowerFor(h.ctx, id)
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
// List Interface Vlans
// @Summary List interface vlans
// @Description List interface vlans info
// @Tags interfaces
// @ID list-interface-vlans
// @Param if_id path string true "if_id"
// @Success 200 {array} godevmandb.Vlan
// @Failure 400 {object} StatusResponse "Invalid if_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/{if_id}/vlans [GET]
func (h *Handler) GetInterfaceVlans(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "if_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid interface ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetInterfaceVlans(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Relations
// List Interface Subinterfaces
// @Summary List interface subinterfaces
// @Description List interface subinterfaces info
// @Tags interfaces
// @ID list-interface-subinterfaces
// @Param if_id path string true "if_id"
// @Success 200 {array} subinterface
// @Failure 400 {object} StatusResponse "Invalid if_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/{if_id}/subinterfaces [GET]
func (h *Handler) GetInterfaceSubinterfaces(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "if_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid interface ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetInterfaceSubinterfaces(h.ctx, &id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := []subinterface{}
	for _, s := range res {
		a := subinterface{}
		a.getValues(s)
		out = append(out, a)
	}

	RespondJSON(w, r, http.StatusOK, out)
}

// Relations
// List Interface Xconnects
// @Summary List interface xconnects
// @Description List interface xconnects info
// @Tags interfaces
// @ID list-interface-xconnects
// @Param if_id path string true "if_id"
// @Success 200 {array} xconnect
// @Failure 400 {object} StatusResponse "Invalid if_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/{if_id}/xconnects [GET]
func (h *Handler) GetInterfaceXconnects(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "if_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid interface ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetInterfaceXconnects(h.ctx, &id)
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
