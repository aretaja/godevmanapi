package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aretaja/godevmandb"
	"github.com/go-chi/chi/v5"
)

// Count UserGraphs
// @Summary Count user_graphs
// @Description Count number of user graphs
// @Tags users
// @ID count-user_graphs
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /users/graphs/count [GET]
func (h *Handler) CountUserGraphs(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountUserGraphs(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List user_graphs
// @Summary List user_graphs
// @Description List user graphs info
// @Tags users
// @ID list-user_graphs
// @Param username_f query string false "url encoded SQL '=' operator pattern"
// @Param descr_f query string false "url encoded SQL '=' operator pattern"
// @Param shared_f query bool false "values 'true', 'false'"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} godevmandb.UserGraph
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /users/graphs [GET]
func (h *Handler) GetUserGraphs(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetUserGraphsParams{
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

	if v := r.FormValue("descr_f"); v != "" {
		p.DescrF = v
	}

	if v := r.FormValue("shared_f"); v != "" {
		p.SharedF = v
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetUserGraphs(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Get UserGraph
// @Summary Get user_graph
// @Description Get user graph info
// @Tags users
// @ID get-user_graph
// @Param graph_id path string true "graph_id"
// @Success 200 {object} godevmandb.UserGraph
// @Failure 400 {object} StatusResponse "Invalid graph_id"
// @Failure 404 {object} StatusResponse "Graph not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /users/graphs/{graph_id} [GET]
func (h *Handler) GetUserGraph(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "graph_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid graph ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetUserGraph(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "Graph not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Create UserGraph
// @Summary Create user_graph
// @Description Create user graph
// @Tags users
// @ID create-user_graph
// @Param Body body godevmandb.CreateUserGraphParams true "JSON object of godevmandb.CreateUserGraphParams"
// @Success 201 {object} godevmandb.UserGraph
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /users/graphs [POST]
func (h *Handler) CreateUserGraph(w http.ResponseWriter, r *http.Request) {
	var p godevmandb.CreateUserGraphParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	q := godevmandb.New(h.db)
	res, err := q.CreateUserGraph(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusCreated, res)
}

// Update UserGraph
// @Summary Update user_graph
// @Description Update user graph
// @Tags users
// @ID update-user_graph
// @Param graph_id path string true "graph_id"
// @Param Body body godevmandb.UpdateUserGraphParams true "JSON object of godevmandb.UpdateUserGraphParams.<br />Ignored fields:<ul><li>graph_id</li></ul>"
// @Success 200 {object} godevmandb.UserGraph
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /users/graphs/{graph_id} [PUT]
func (h *Handler) UpdateUserGraph(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "graph_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid graph ID")
		return
	}

	var p godevmandb.UpdateUserGraphParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	p.GraphID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateUserGraph(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Delete UserGraph
// @Summary Delete user_graph
// @Description Delete user graph
// @Tags users
// @ID delete-user_graph
// @Param graph_id path string true "graph_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid graph_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /users/graphs/{graph_id} [DELETE]
func (h *Handler) DeleteUserGraph(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "graph_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid graph ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteUserGraph(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
