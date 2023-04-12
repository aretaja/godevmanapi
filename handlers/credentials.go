package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aretaja/godevmandb"
	"github.com/go-chi/chi/v5"
)

// Count Credentials
// @Summary Count credentials
// @Description Count number of credentials
// @Tags config
// @ID count-credentials
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /config/credentials/count [GET]
func (h *Handler) CountCredentials(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountCredentials(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List credentials
// @Summary List credentials
// @Description List credentials info
// @Tags config
// @ID list-credentials
// @Param label_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} godevmandb.Credential
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /config/credentials [GET]
func (h *Handler) GetCredentials(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetCredentialsParams{
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
	v := r.FormValue("label_f")
	if v != "" {
		p.LabelF = v
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetCredentials(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// Decrypt secret
	for i, s := range res {
		if s.EncSecret != "" {
			val, err := godevmandb.DecryptStrAes(s.EncSecret, salt)
			if err != nil {
				RespondError(w, r, http.StatusInternalServerError, err.Error())
				return
			}

			res[i].EncSecret = val
		}
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Get Credential
// @Summary Get credential
// @Description Get credential info
// @Tags config
// @ID get-credential
// @Param cred_id path string true "cred_id"
// @Success 200 {object} godevmandb.Credential
// @Failure 400 {object} StatusResponse "Invalid cred_id"
// @Failure 404 {object} StatusResponse "Credential not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /config/credentials/{cred_id} [GET]
func (h *Handler) GetCredential(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "cred_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid credential ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetCredential(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "Credential not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Decrypt secret
	if res.EncSecret != "" {
		val, err := godevmandb.DecryptStrAes(res.EncSecret, salt)
		if err != nil {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		res.EncSecret = val
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Create Credential
// @Summary Create credential
// @Description Create credential
// @Tags config
// @ID create-credential
// @Param Body body godevmandb.CreateCredentialParams true "JSON object of godevmandb.CreateCredentialParams"
// @Success 201 {object} godevmandb.Credential
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /config/credentials [POST]
func (h *Handler) CreateCredential(w http.ResponseWriter, r *http.Request) {
	var p godevmandb.CreateCredentialParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Encrypt secret
	if p.EncSecret != "" {
		val, err := godevmandb.EncryptStrAes(p.EncSecret, salt)
		if err != nil {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		p.EncSecret = val
	}

	q := godevmandb.New(h.db)
	res, err := q.CreateCredential(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// Decrypt secret
	if res.EncSecret != "" {
		val, err := godevmandb.DecryptStrAes(res.EncSecret, salt)
		if err != nil {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		res.EncSecret = val
	}

	RespondJSON(w, r, http.StatusCreated, res)
}

// Update Credential
// @Summary Update credential
// @Description Update credential
// @Tags config
// @ID update-credential
// @Param cred_id path string true "cred_id"
// @Param Body body godevmandb.UpdateCredentialParams true "JSON object of godevmandb.UpdateCredentialParams.<br />Ignored fields:<ul><li>cred_id</li></ul>"
// @Success 200 {object} godevmandb.Credential
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /config/credentials/{cred_id} [PUT]
func (h *Handler) UpdateCredential(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "cred_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid credential ID")
		return
	}

	var p godevmandb.UpdateCredentialParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	p.CredID = id

	// Encrypt secret
	if p.EncSecret != "" {
		val, err := godevmandb.EncryptStrAes(p.EncSecret, salt)
		if err != nil {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		p.EncSecret = val
	}

	q := godevmandb.New(h.db)
	res, err := q.UpdateCredential(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// Decrypt secret
	if res.EncSecret != "" {
		val, err := godevmandb.DecryptStrAes(res.EncSecret, salt)
		if err != nil {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		res.EncSecret = val
	}

	RespondJSON(w, r, http.StatusOK, res)
}

// Delete Credential
// @Summary Delete credential
// @Description Delete credential
// @Tags config
// @ID delete-credential
// @Param cred_id path string true "cred_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid cred_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /config/credentials/{cred_id} [DELETE]
func (h *Handler) DeleteCredential(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "cred_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid credential ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteCredential(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
