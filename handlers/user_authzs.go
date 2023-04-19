package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aretaja/godevmandb"
	"github.com/go-chi/chi/v5"
)

// Count UserAuthzs
// @Summary Count user_authzs
// @Description Count number of user_authzs
// @Tags users
// @ID count-user_authzs
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /users/authzs/count [GET]
func (h *Handler) CountUserAuthzs(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountUserAuthzs(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List user_authzs
// @Summary List user_authzs
// @Description List user_authzs info
// @Tags users
// @ID list-user_authzs
// @Param username_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param userlevel_le query int false "SQL '<=' operator value"
// @Param userlevel_ge query int false "SQL '>=' operator value"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} godevmandb.UserAuthz
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /users/authzs [GET]
func (h *Handler) GetUserAuthzs(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetUserAuthzsParams{
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
	if v := r.FormValue("username_f"); v != "" {
		p.UsernameF = v
	}

	if v := r.FormValue("userlevel_ge"); v != "" {
		p.UserlevelGe = v
	}

	if v := r.FormValue("userlevel_le"); v != "" {
		p.UserlevelLe = v
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetUserAuthzs(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Get UserAuthz
// @Summary Get user_authz
// @Description Get user_authz info
// @Tags users
// @ID get-user_authz
// @Param username path string true "username"
// @Param dom_id path string true "dom_id"
// @Success 200 {object} godevmandb.UserAuthz
// @Failure 400 {object} StatusResponse "Invalid username/dom_id"
// @Failure 404 {object} StatusResponse "UserAuthz not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /users/authzs/{username}/{dom_id} [GET]
func (h *Handler) GetUserAuthz(w http.ResponseWriter, r *http.Request) {
	var p godevmandb.GetUserAuthzParams
	id, err := strconv.ParseInt(chi.URLParam(r, "dom_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid domain ID")
		return
	}

	p.Username = chi.URLParam(r, "username")
	p.DomID = id

	q := godevmandb.New(h.db)
	res, err := q.GetUserAuthz(h.ctx, p)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "User not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Create UserAuthz
// @Summary Create user_authz
// @Description Create user_authz
// @Tags users
// @ID create-user_authz
// @Param Body body godevmandb.CreateUserAuthzParams true "JSON object of godevmandb.CreateUserAuthzParams"
// @Success 201 {object} godevmandb.UserAuthz
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /users/authzs [POST]
func (h *Handler) CreateUserAuthz(w http.ResponseWriter, r *http.Request) {
	var p godevmandb.CreateUserAuthzParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	q := godevmandb.New(h.db)
	res, err := q.CreateUserAuthz(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusCreated, res)
}

// Update UserAuthz
// @Summary Update user_authz
// @Description Update user_authz
// @Tags users
// @ID update-user_authz
// @Param username path string true "username"
// @Param dom_id path string true "dom_id"
// @Param Body body godevmandb.UpdateUserAuthzParams true "JSON object of godevmandb.UpdateUserAuthzParams.<br />Ignored fields:<ul><li>username</li><li>dom_id</li></ul>"
// @Success 200 {object} godevmandb.UserAuthz
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /users/authzs/{username}/{dom_id} [PUT]
func (h *Handler) UpdateUserAuthz(w http.ResponseWriter, r *http.Request) {
	var p godevmandb.UpdateUserAuthzParams
	id, err := strconv.ParseInt(chi.URLParam(r, "dom_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid domain ID")
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	p.Username = chi.URLParam(r, "username")
	p.DomID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateUserAuthz(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Delete UserAuthz
// @Summary Delete user_authz
// @Description Delete user_authz
// @Tags users
// @ID delete-user_authz
// @Param username path string true "username"
// @Param dom_id path string true "dom_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid username/dom_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /users/authzs/{username}/{dom_id} [DELETE]
func (h *Handler) DeleteUserAuthz(w http.ResponseWriter, r *http.Request) {
	var p godevmandb.DeleteUserAuthzParams
	id, err := strconv.ParseInt(chi.URLParam(r, "dom_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid domain ID")
		return
	}

	p.Username = chi.URLParam(r, "username")
	p.DomID = id

	q := godevmandb.New(h.db)
	err = q.DeleteUserAuthz(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Foreign key
// Get UserAuthz DeviceDomain
// @Summary Get user_authz device_domain
// @Description Get user_authz device_domain info
// @Tags users
// @ID get-user_authz-device-domain
// @Param username path string true "username"
// @Param dom_id path string true "dom_id"
// @Success 200 {object} godevmandb.DeviceDomain
// @Failure 400 {object} StatusResponse "Invalid username/dom_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /users/authzs/{username}/{dom_id}/device_domain [GET]
func (h *Handler) GetUserAuthzDeviceDomain(w http.ResponseWriter, r *http.Request) {
	var p godevmandb.GetUserAuthzDeviceDomainParams
	id, err := strconv.ParseInt(chi.URLParam(r, "dom_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid username and/or domain ID")
		return
	}

	p.Username = chi.URLParam(r, "username")
	p.DomID = id

	q := godevmandb.New(h.db)
	res, err := q.GetUserAuthzDeviceDomain(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}
