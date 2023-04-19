package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aretaja/godevmandb"
	"github.com/go-chi/chi/v5"
)

// Count Users
// @Summary Count users
// @Description Count number of users
// @Tags users
// @ID count-users
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /users/count [GET]
func (h *Handler) CountUsers(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountUsers(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List users
// @Summary List users
// @Description List users info
// @Tags users
// @ID list-users
// @Param username_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param notes_f query string false "url encoded SQL 'ILIKE' operator pattern + special values 'isnull', 'isempty'"
// @Param userlevel_le query int false "SQL '<=' operator value"
// @Param userlevel_ge query int false "SQL '>=' operator value"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} godevmandb.User
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /users [GET]
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetUsersParams{
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

	if v := r.FormValue("notes_f"); v != "" {
		p.NotesF = &v
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetUsers(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Get User
// @Summary Get user
// @Description Get user info
// @Tags users
// @ID get-user
// @Param username path string true "username"
// @Success 200 {object} godevmandb.User
// @Failure 400 {object} StatusResponse "Invalid username"
// @Failure 404 {object} StatusResponse "User not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /users/{username} [GET]
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "username")

	q := godevmandb.New(h.db)
	res, err := q.GetUser(h.ctx, id)
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

// Create User
// @Summary Create user
// @Description Create user
// @Tags users
// @ID create-user
// @Param Body body godevmandb.CreateUserParams true "JSON object of godevmandb.CreateUserParams"
// @Success 201 {object} godevmandb.User
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /users [POST]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var p godevmandb.CreateUserParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	q := godevmandb.New(h.db)
	res, err := q.CreateUser(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusCreated, res)
}

// Update User
// @Summary Update user
// @Description Update user
// @Tags users
// @ID update-user
// @Param username path string true "username"
// @Param Body body godevmandb.UpdateUserParams true "JSON object of godevmandb.UpdateUserParams.<br />Ignored fields:<ul><li>username</li></ul>"
// @Success 200 {object} godevmandb.User
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /users/{username} [PUT]
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "username")

	var p godevmandb.UpdateUserParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	p.Username = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateUser(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Delete User
// @Summary Delete user
// @Description Delete user
// @Tags users
// @ID delete-user
// @Param username path string true "username"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid username"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /users/{username} [DELETE]
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "username")

	q := godevmandb.New(h.db)
	err := q.DeleteUser(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Relations
// List User Graphs
// @Summary List user graphs
// @Description List user graphs info
// @Tags users
// @ID list-user-graphs
// @Param username path string true "username"
// @Success 200 {array} godevmandb.UserGraph
// @Failure 400 {object} StatusResponse "Invalid username"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /users/{username}/graphs [GET]
func (h *Handler) GetUserUserGraphs(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "username")

	q := godevmandb.New(h.db)
	res, err := q.GetUserUserGraphs(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Relations
// List User Authzs
// @Summary List user authzs
// @Description List user authzs info
// @Tags users
// @ID list-user-authzs
// @Param username path string true "username"
// @Success 200 {array} godevmandb.UserAuthz
// @Failure 400 {object} StatusResponse "Invalid username"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /users/{username}/authzs [GET]
func (h *Handler) GetUserUserAuthzs(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "username")

	q := godevmandb.New(h.db)
	res, err := q.GetUserUserAuthzs(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}
