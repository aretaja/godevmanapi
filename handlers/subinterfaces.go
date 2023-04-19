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
type subinterface struct {
	UpdatedOn time.Time `json:"updated_on"`
	CreatedOn time.Time `json:"created_on"`
	Mac       *string   `json:"mac"`
	Alias     *string   `json:"alias"`
	Oper      *int16    `json:"oper"`
	Adm       *int16    `json:"adm"`
	Speed     *int64    `json:"speed"`
	TypeEnum  *string   `json:"type_enum"`
	Notes     *string   `json:"notes"`
	Ifindex   *int64    `json:"ifindex"`
	IfID      *int64    `json:"if_id"`
	Descr     string    `json:"descr"`
	SifID     int64     `json:"sif_id"`
}

// Import values from corresponding godevmandb struct
func (r *subinterface) getValues(s godevmandb.Subinterface) {
	r.SifID = s.SifID
	r.IfID = s.IfID
	r.Ifindex = s.Ifindex
	r.Descr = s.Descr
	r.Alias = s.Alias
	r.Oper = s.Oper
	r.Adm = s.Adm
	r.Speed = s.Speed
	r.TypeEnum = s.TypeEnum
	r.Notes = s.Notes
	r.UpdatedOn = s.UpdatedOn
	r.CreatedOn = s.CreatedOn
	r.Mac = pgMacaddrToPtr(s.Mac)
}

// Return corresponding godevmandb create parameters
func (r *subinterface) createParams() godevmandb.CreateSubinterfaceParams {
	s := godevmandb.CreateSubinterfaceParams{}

	s.IfID = r.IfID
	s.Ifindex = r.Ifindex
	s.Descr = r.Descr
	s.Alias = r.Alias
	s.Oper = r.Oper
	s.Adm = r.Adm
	s.Speed = r.Speed
	s.TypeEnum = r.TypeEnum
	s.Notes = r.Notes
	s.Mac = strToPgMacaddr(r.Mac)

	return s
}

// Return corresponding godevmandb update parameters
func (r *subinterface) updateParams() godevmandb.UpdateSubinterfaceParams {
	s := godevmandb.UpdateSubinterfaceParams{}

	s.IfID = r.IfID
	s.Ifindex = r.Ifindex
	s.Descr = r.Descr
	s.Alias = r.Alias
	s.Oper = r.Oper
	s.Adm = r.Adm
	s.Speed = r.Speed
	s.TypeEnum = r.TypeEnum
	s.Notes = r.Notes
	s.Mac = strToPgMacaddr(r.Mac)

	return s
}

// Count Subinterfaces
// @Summary Count subinterfaces
// @Description Count number of subinterfaces
// @Tags interfaces
// @ID count-subinterfaces
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/subinterfaces/count [GET]
func (h *Handler) CountSubinterfaces(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountSubinterfaces(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List subinterfaces
// @Summary List subinterfaces
// @Description List subinterfaces info
// @Tags interfaces
// @ID list-subinterfaces
// @Param ifindex_f query string false "url encoded SQL 'LIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param descr_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param alias_f query string false "url encoded SQL 'ILIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param oper_f query string false "url encoded SQL 'LIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param adm_f query string false "url encoded SQL 'LIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param speed_f query string false "url encoded SQL 'LIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param type_enum_f query string false "url encoded SQL 'LIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param notes_f query string false "url encoded SQL 'LIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param mac_f query string false "SQL '=' operator value (MAC address)"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} subinterface
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/subinterfaces [GET]
func (h *Handler) GetSubinterfaces(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetSubinterfacesParams{
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

	if v := r.FormValue("type_enum_f"); v != "" {
		p.TypeEnumF = &v
	}

	p.MacF = strToPgMacaddr(nil)
	if v := r.FormValue("mac_f"); v != "" {
		p.MacF = strToPgMacaddr(&v)
	}

	if v := r.FormValue("notes_f"); v != "" {
		p.NotesF = &v
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetSubinterfaces(h.ctx, p)
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

// Get Subinterface
// @Summary Get interface
// @Description Get subinterface info
// @Tags interfaces
// @ID get-subinterface
// @Param sif_id path string true "sif_id"
// @Success 200 {object} subinterface
// @Failure 400 {object} StatusResponse "Invalid sif_id"
// @Failure 404 {object} StatusResponse "Subinterface not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/subinterfaces/{sif_id} [GET]
func (h *Handler) GetSubinterface(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "sif_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid subinterface ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetSubinterface(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "Subinterface not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	out := subinterface{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Create Subinterface
// @Summary Create subinterface
// @Description Create subinterface
// @Tags interfaces
// @ID create-subinterface
// @Param Body body subinterface true "JSON object of subinterface.<br />Ignored fields:<ul><li>sif_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 201 {object} subinterface
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/subinterfaces [POST]
func (h *Handler) CreateSubinterface(w http.ResponseWriter, r *http.Request) {
	var pIn subinterface
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Create parameters for new db record
	p := pIn.createParams()

	q := godevmandb.New(h.db)
	res, err := q.CreateSubinterface(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := subinterface{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusCreated, out)
}

// Update Subinterface
// @Summary Update subinterface
// @Description Update subinterface
// @Tags interfaces
// @ID update-subinterface
// @Param sif_id path string true "sif_id"
// @Param Body body subinterface true "JSON object of subinterface.<br />Ignored fields:<ul><li>sif_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 200 {object} subinterface
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/subinterfaces/{sif_id} [PUT]
func (h *Handler) UpdateSubinterface(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "sif_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid subinterface ID")
		return
	}

	var pIn subinterface
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Update parameters for new db record
	p := pIn.updateParams()
	p.SifID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateSubinterface(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := subinterface{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Delete Subinterface
// @Summary Delete subinterface
// @Description Delete subinterface
// @Tags interfaces
// @ID delete-subinterface
// @Param sif_id path string true "sif_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid sif_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/subinterfaces/{sif_id} [DELETE]
func (h *Handler) DeleteSubinterface(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "sif_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid subinterface ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteSubinterface(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Foreign key
// Get Subinterface Interface
// @Summary Get subinterface interface
// @Description Get subinterface interface info
// @Tags interfaces
// @ID get-subinterface-interface
// @Param sif_id path string true "sif_id"
// @Success 200 {object} iface
// @Failure 400 {object} StatusResponse "Invalid sif_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /interfaces/subinterfaces/{sif_id}/interface [GET]
func (h *Handler) GetSubinterfaceInterface(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "sif_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid subinterface ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetSubinterfaceInterface(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := iface{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}
