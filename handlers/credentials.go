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
type credential struct {
	UpdatedOn time.Time `json:"updated_on"`
	CreatedOn time.Time `json:"created_on"`
	Username  *string   `json:"username"`
	EncSecret string    `json:"enc_secret"`
	Label     string    `json:"label"`
	CredID    int64     `json:"cred_id"`
}

// Import values from corresponding godevmandb struct
func (r *credential) getValues(s godevmandb.Credential) error {
	r.CredID = s.CredID
	r.Label = s.Label
	r.UpdatedOn = s.UpdatedOn
	r.CreatedOn = s.CreatedOn
	r.Username = nullStringToPtr(s.Username)

	val, err := godevmandb.DecryptStrAes(s.EncSecret, salt)
	if err != nil {
		return err
	}

	r.EncSecret = val

	return nil
}

// Return corresponding godevmandb create parameters
func (r *credential) createParams() (godevmandb.CreateCredentialParams, error) {
	s := godevmandb.CreateCredentialParams{}

	s.Label = r.Label
	s.Username = strToNullString(r.Username)

	val, err := godevmandb.EncryptStrAes(r.EncSecret, salt)
	if err != nil {
		return s, err
	}

	s.EncSecret = val

	return s, nil
}

// Return corresponding godevmandb update parameters
func (r *credential) updateParams() (godevmandb.UpdateCredentialParams, error) {
	s := godevmandb.UpdateCredentialParams{}

	s.Label = r.Label
	s.Username = strToNullString(r.Username)

	val, err := godevmandb.EncryptStrAes(r.EncSecret, salt)
	if err != nil {
		return s, err
	}

	s.EncSecret = val

	return s, nil
}

// Count Credentials
// @Summary Count credentials
// @Description Count number of credentials
// @Tags data
// @ID count-credentials
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /data/credentials/count [GET]
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
// @Tags data
// @ID list-credentials
// @Param label_f query string false "url encoded SQL 'LIKE' operator pattern"
// @Param limit query int false "min: 1; max: 100; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} credential
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /data/credentials [GET]
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

	// Descr filter
	d := r.FormValue("label_f")
	if d != "" {
		p.LabelF = d
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetCredentials(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := []credential{}
	for _, s := range res {
		r := credential{}
		r.getValues(s)
		out = append(out, r)
	}

	RespondJSON(w, r, http.StatusOK, out)
}

// Get Credential
// @Summary Get credential
// @Description Get credential info
// @Tags data
// @ID get-credential
// @Param cred_id path string true "cred_id"
// @Success 200 {object} credential
// @Failure 400 {object} StatusResponse "Invalid cred_id"
// @Failure 404 {object} StatusResponse "Credential not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /data/credentials/{cred_id} [GET]
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

	out := credential{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Create Credential
// @Summary Create credential
// @Description Create credential
// @Tags data
// @ID create-credential
// @Param Body body credential true "JSON object of credential.<br />Ignored fields:<ul><li>cred_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 201 {object} credential
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /data/credentials [POST]
func (h *Handler) CreateCredential(w http.ResponseWriter, r *http.Request) {
	var pIn credential
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Create parameters for new db record
	p, err := pIn.createParams()
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.CreateCredential(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := credential{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusCreated, out)
}

// Update Credential
// @Summary Update credential
// @Description Update credential
// @Tags data
// @ID update-credential
// @Param cred_id path string true "cred_id"
// @Param Body body credential true "JSON object of credential.<br />Ignored fields:<ul><li>cred_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 200 {object} credential
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /data/credentials/{cred_id} [PUT]
func (h *Handler) UpdateCredential(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "cred_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid credential ID")
		return
	}

	var pIn credential
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Update parameters for new db record
	p, err := pIn.updateParams()
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	p.CredID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateCredential(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := credential{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Delete Credential
// @Summary Delete credential
// @Description Delete credential
// @Tags data
// @ID delete-credential
// @Param cred_id path string true "cred_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid cred_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /data/credentials/{cred_id} [DELETE]
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
