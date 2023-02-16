package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aretaja/godevmandb"
	"github.com/go-chi/chi/v5"
)

// Count DeviceCredentials
// @Summary Count credentials
// @Description Count number of credentials
// @Tags devices
// @ID count-device_credentials
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/credentials/count [GET]
func (h *Handler) CountDeviceCredentials(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountDeviceCredentials(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List DeviceCredentials
// @Summary List device_credentials
// @Description List device credentials info
// @Tags devices
// @ID list-device_credentials
// @Param username_f query string false "url encoded SQL 'LIKE' operator pattern"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} godevmandb.DeviceCredential
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/credentials [GET]
func (h *Handler) GetDeviceCredentials(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetDeviceCredentialsParams{
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
	d := r.FormValue("username_f")
	if d != "" {
		p.UsernameF = d
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetDeviceCredentials(h.ctx, p)
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

// Get DeviceCredential
// @Summary Get device_credential
// @Description Get device credential info
// @Tags devices
// @ID get-device_credential
// @Param cred_id path string true "cred_id"
// @Success 200 {object} godevmandb.DeviceCredential
// @Failure 400 {object} StatusResponse "Invalid cred_id"
// @Failure 404 {object} StatusResponse "Credential not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/credentials/{cred_id} [GET]
func (h *Handler) GetDeviceCredential(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "cred_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid credential ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetDeviceCredential(h.ctx, id)
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

// Create DeviceCredential
// @Summary Create device_credential
// @Description Create device credential
// @Tags devices
// @ID create-device_credential
// @Param Body body godevmandb.CreateDeviceCredentialParams true "JSON object of godevmandb.CreateDeviceCredentialParams"
// @Success 201 {object} godevmandb.DeviceCredential
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/credentials [POST]
func (h *Handler) CreateDeviceCredential(w http.ResponseWriter, r *http.Request) {
	var p godevmandb.CreateDeviceCredentialParams
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
	res, err := q.CreateDeviceCredential(h.ctx, p)
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

// Update DeviceCredential
// @Summary Update device_credential
// @Description Update device credential
// @Tags devices
// @ID update-device_credential
// @Param cred_id path string true "cred_id"
// @Param Body body godevmandb.UpdateDeviceCredentialParams true "JSON object of godevmandb.UpdateDeviceCredentialParams.<br />Ignored fields:<ul><li>cred_id</li></ul>"
// @Success 200 {object} godevmandb.DeviceCredential
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/credentials/{cred_id} [PUT]
func (h *Handler) UpdateDeviceCredential(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "cred_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid credential ID")
		return
	}

	var p godevmandb.UpdateDeviceCredentialParams
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
	res, err := q.UpdateDeviceCredential(h.ctx, p)

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

// Delete DeviceCredential
// @Summary Delete device_credential
// @Description Delete device credential
// @Tags devices
// @ID delete-device_credential
// @Param cred_id path string true "cred_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid cred_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /devices/credentials/{cred_id} [DELETE]
func (h *Handler) DeleteDeviceCredential(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "cred_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid credential ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteDeviceCredential(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
