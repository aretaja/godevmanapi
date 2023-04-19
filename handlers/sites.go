package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aretaja/godevmandb"
	"github.com/go-chi/chi/v5"
)

// Count Sites
// @Summary Count sites
// @Description Count number of sites
// @Tags sites
// @ID count-sites
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /sites/count [GET]
func (h *Handler) CountSites(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountSites(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List sites
// @Summary List sites
// @Description List site info
// @Tags sites
// @ID list-sites
// @Param descr_f query string false "url encoded SQL 'ILIKE' operator pattern + special value 'isempty'"
// @Param uident_f query string false "url encoded SQL 'ILIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param area_f query string false "url encoded SQL 'ILIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param addr_f query string false "url encoded SQL 'ILIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param notes_f query string false "url encoded SQL 'ILIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param ext_name_f query string false "url encoded SQL 'ILIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param ext_id_f query string false "url encoded SQL 'ILIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} godevmandb.Site
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /sites [GET]
func (h *Handler) GetSites(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetSitesParams{
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
	if v := r.FormValue("uident_f"); v != "" {
		p.UidentF = &v
	}

	if v := r.FormValue("descr_f"); v != "" {
		p.DescrF = v
	}

	if v := r.FormValue("area_f"); v != "" {
		p.AreaF = &v
	}

	if v := r.FormValue("addr_f"); v != "" {
		p.AddrF = &v
	}

	if v := r.FormValue("notes_f"); v != "" {
		p.NotesF = &v
	}

	if v := r.FormValue("extid_f"); v != "" {
		p.ExtIDF = &v
	}

	if v := r.FormValue("extname_f"); v != "" {
		p.ExtNameF = &v
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetSites(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Get Site
// @Summary Get site
// @Description Get site info
// @Tags sites
// @ID get-site
// @Param site_id path string true "site_id"
// @Success 200 {object} godevmandb.Site
// @Failure 400 {object} StatusResponse "Invalid site_id"
// @Failure 404 {object} StatusResponse "Site not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /sites/{site_id} [GET]
func (h *Handler) GetSite(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "site_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid site ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetSite(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "Site not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Create Site
// @Summary Create site
// @Description Create site
// @Tags sites
// @ID create-site
// @Param Body body godevmandb.CreateSiteParams true "JSON object of godevmandb.CreateSiteParams"
// @Success 201 {object} godevmandb.Site
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /sites [POST]
func (h *Handler) CreateSite(w http.ResponseWriter, r *http.Request) {
	var p godevmandb.CreateSiteParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	q := godevmandb.New(h.db)
	res, err := q.CreateSite(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusCreated, res)
}

// Update Site
// @Summary Update site
// @Description Update site
// @Tags sites
// @ID update-site
// @Param site_id path string true "site_id"
// @Param Body body godevmandb.UpdateSiteParams true "JSON object of godevmandb.UpdateSiteParams.<br />Ignored fields:<ul><li>site_id</li></ul>"
// @Success 200 {object} godevmandb.Site
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /sites/{site_id} [PUT]
func (h *Handler) UpdateSite(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "site_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid site ID")
		return
	}

	var p godevmandb.UpdateSiteParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Update parameters for new db record
	p.SiteID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateSite(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Delete Site
// @Summary Delete site
// @Description Delete site
// @Tags sites
// @ID delete-site
// @Param site_id path string true "site_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid site_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /sites/{site_id} [DELETE]
func (h *Handler) DeleteSite(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "site_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid site ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteSite(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Foreign key
// Get Site Country
// @Summary Get site country
// @Description Get site country info
// @Tags sites
// @ID get-site-country
// @Param site_id path string true "site_id"
// @Success 200 {object} godevmandb.Country
// @Failure 400 {object} StatusResponse "Invalid site_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /sites/{site_id}/country [GET]
func (h *Handler) GetSiteConCountry(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "site_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid site ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetSiteCountry(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Relations
// List Site Devices
// @Summary List site devices
// @Description List site devices info
// @Tags sites
// @ID list-site-devices
// @Param site_id path string true "site_id"
// @Success 200 {array} device
// @Failure 400 {object} StatusResponse "Invalid site_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /sites/{site_id}/devices [GET]
func (h *Handler) GetSiteDevices(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "site_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid site ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetSiteDevices(h.ctx, &id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := []device{}
	for _, s := range res {
		r := device{}
		r.getValues(s)
		out = append(out, r)
	}

	RespondJSON(w, r, http.StatusOK, out)
}

// Relations
// List Site Connections
// @Summary List site connections
// @Description List site connections info
// @Tags sites
// @ID list-site-connections
// @Param site_id path string true "site_id"
// @Success 200 {array} godevmandb.Connection
// @Failure 400 {object} StatusResponse "Invalid site_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /sites/{site_id}/connections [GET]
func (h *Handler) GetSiteConnections(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "site_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid site ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetSiteConnections(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}
