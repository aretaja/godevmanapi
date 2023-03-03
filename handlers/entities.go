package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aretaja/godevmandb"
	"github.com/go-chi/chi/v5"
)

// Count Entities
// @Summary Count entities
// @Description Count number of entities
// @Tags entities
// @ID count-entities
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /entities/count [GET]
func (h *Handler) CountEntities(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountEntities(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List entities
// @Summary List entities
// @Description List entities info
// @Tags entities
// @ID list-entities
// @Param slot_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param descr_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param model_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param w_product_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param hw_revision_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param serial_nr_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param sw_product_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param sw_revision_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param manufacturer_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} godevmandb.Entity
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /entities [GET]
func (h *Handler) GetEntities(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetEntitiesParams{
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

	// Slot filter
	d := r.FormValue("slot_f")
	if d != "" {
		p.SlotF = d
	}

	// Descr filter
	d = r.FormValue("descr_f")
	if d != "" {
		p.DescrF = d
	}

	// Model filter
	d = r.FormValue("model_f")
	if d != "" {
		p.ModelF = d
	}

	// HwProduct filter
	d = r.FormValue("hw_product_f")
	if d != "" {
		p.HwProductF = d
	}

	// HwRevision filter
	d = r.FormValue("hw_revision_f")
	if d != "" {
		p.HwRevisionF = d
	}

	// Serial nr filter
	d = r.FormValue("serial_nr_f")
	if d != "" {
		p.SerialNrF = d
	}

	// SwProduct filter
	d = r.FormValue("sw_product_f")
	if d != "" {
		p.SwProductF = d
	}

	// SwRevision filter
	d = r.FormValue("sw_revision_f")
	if d != "" {
		p.SwRevisionF = d
	}

	// Manufacturer filter
	d = r.FormValue("manufacturer_f")
	if d != "" {
		p.ManufacturerF = d
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetEntities(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Get Entity
// @Summary Get Entity
// @Description Get Entity info
// @Tags entities
// @ID get-Entity
// @Param ent_id path string true "ent_id"
// @Success 200 {object} godevmandb.Entity
// @Failure 400 {object} StatusResponse "Invalid ent_id"
// @Failure 404 {object} StatusResponse "Entity not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /entities/{ent_id} [GET]
func (h *Handler) GetEntity(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ent_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid Entity ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetEntity(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "Entity not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Create Entity
// @Summary Create Entity
// @Description Create Entity
// @Tags entities
// @ID create-Entity
// @Param Body body godevmandb.CreateEntityParams true "JSON object of godevmandb.CreateEntityParams"
// @Success 201 {object} godevmandb.Entity
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /entities [POST]
func (h *Handler) CreateEntity(w http.ResponseWriter, r *http.Request) {
	var p godevmandb.CreateEntityParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	q := godevmandb.New(h.db)
	res, err := q.CreateEntity(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusCreated, res)
}

// Update Entity
// @Summary Update Entity
// @Description Update Entity
// @Tags entities
// @ID update-Entity
// @Param ent_id path string true "ent_id"
// @Param Body body godevmandb.UpdateEntityParams true "JSON object of godevmandb.UpdateEntityParams.<br />Ignored fields:<ul><li>ent_id</li></ul>"
// @Success 200 {object} godevmandb.Entity
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /entities/{ent_id} [PUT]
func (h *Handler) UpdateEntity(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ent_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid Entity ID")
		return
	}

	var p godevmandb.UpdateEntityParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	p.EntID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateEntity(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Delete Entity
// @Summary Delete Entity
// @Description Delete Entity
// @Tags entities
// @ID delete-Entity
// @Param ent_id path string true "ent_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid ent_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /entities/{ent_id} [DELETE]
func (h *Handler) DeleteEntity(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ent_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid Entity ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteEntity(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Foreign key
// Get Entity Device
// @Summary Get entity device
// @Description Get entity device info
// @Tags entities
// @ID get-entity-device
// @Param ent_id path string true "ent_id"
// @Success 200 {object} device
// @Failure 400 {object} StatusResponse "Invalid ent_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /entities/{ent_id}/device [GET]
func (h *Handler) GetEntityDevice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ent_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid entity ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetEntityDevice(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := device{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Foreign key
// Get Entity Parent
// @Summary Get entity parent
// @Description Get entity parent info
// @Tags entities
// @ID get-entity-parent
// @Param ent_id path string true "ent_id"
// @Success 200 {object} godevmandb.Entity
// @Failure 400 {object} StatusResponse "Invalid ent_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /entities/{ent_id}/parent [GET]
func (h *Handler) GetEntityParent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ent_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid entity ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetEntityParent(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Relations
// List Entity Childs
// @Summary List entity childs
// @Description List connection entity childs info
// @Tags entities
// @ID list-entity-childs
// @Param ent_id path string true "ent_id"
// @Success 200 {array} godevmandb.Entity
// @Failure 400 {object} StatusResponse "Invalid ent_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /entities/{ent_id}/childs [GET]
func (h *Handler) GetEntityChilds(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ent_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid entity ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetEntityChilds(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Relations
// List Entity PhyIndexes
// @Summary List entity phy indexes
// @Description List connection entity phy indexes info
// @Tags entities
// @ID list-entity-phy-indexes
// @Param ent_id path string true "ent_id"
// @Success 200 {array} godevmandb.EntityPhyIndex
// @Failure 400 {object} StatusResponse "Invalid ent_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /entities/{ent_id}/entity_phy_indexes [GET]
func (h *Handler) GetEntityEntityPhyIndexes(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ent_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid entity ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetEntityEntityPhyIndexes(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Relations
// List Entity Interfaces
// @Summary List entity interfaces
// @Description List connection entity interfaces info
// @Tags entities
// @ID list-entity-interfaces
// @Param ent_id path string true "ent_id"
// @Success 200 {array} iface
// @Failure 400 {object} StatusResponse "Invalid ent_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /entities/{ent_id}/interfaces [GET]
func (h *Handler) GetEntityInterfaces(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ent_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid entity ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetEntityInterfaces(h.ctx, &id)
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
// List Entity RlfNbrs
// @Summary List entity rl_nbrs
// @Description List connection entity rl_nbrs info
// @Tags entities
// @ID list-entity-rl_nbrs
// @Param ent_id path string true "ent_id"
// @Success 200 {array} godevmandb.RlNbr
// @Failure 400 {object} StatusResponse "Invalid ent_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /entities/{ent_id}/rl_nbrs [GET]
func (h *Handler) GetEntityRlfNbrs(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ent_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid entity ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetEntityRlfNbrs(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}
