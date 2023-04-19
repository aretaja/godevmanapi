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
type snmpCredential struct {
	UpdatedOn  time.Time                 `json:"updated_on"`
	CreatedOn  time.Time                 `json:"created_on"`
	AuthProto  *godevmandb.SnmpAuthProto `json:"auth_proto"`
	AuthPass   *string                   `json:"auth_pass"`
	SecLevel   *godevmandb.SnmpSecLevel  `json:"sec_level"`
	PrivProto  *godevmandb.SnmpPrivProto `json:"priv_proto"`
	PrivPass   *string                   `json:"priv_pass"`
	Label      string                    `json:"label"`
	AuthName   string                    `json:"auth_name"`
	SnmpCredID int64                     `json:"snmp_snmp_cred_id"`
	Variant    int32                     `json:"variant"`
}

// Import values from corresponding godevmandb struct
func (r *snmpCredential) getValues(s godevmandb.SnmpCredential) error {
	r.SnmpCredID = s.SnmpCredID
	r.Variant = s.Variant
	r.Label = s.Label
	r.AuthName = s.AuthName
	r.AuthProto = nullSnmpAuthProtoToPtr(s.AuthProto)
	r.SecLevel = nullSnmpSecLevelToPtr(s.SecLevel)
	r.PrivProto = nullSnmpPrivProtoToPtr(s.PrivProto)
	r.UpdatedOn = s.UpdatedOn
	r.CreatedOn = s.CreatedOn

	if s.AuthPass != nil {
		enc := *s.AuthPass
		val, err := godevmandb.DecryptStrAes(enc, salt)
		if err != nil {
			return err
		}

		r.AuthPass = &val
	}

	if s.PrivPass != nil {
		enc := *s.PrivPass
		val, err := godevmandb.DecryptStrAes(enc, salt)
		if err != nil {
			return err
		}

		r.PrivPass = &val
	}

	return nil
}

// Return corresponding godevmandb create parameters
func (r *snmpCredential) createParams() (godevmandb.CreateSnmpCredentialParams, error) {
	s := godevmandb.CreateSnmpCredentialParams{}

	s.Variant = r.Variant
	s.Label = r.Label
	s.AuthName = r.AuthName
	s.AuthProto = snmpAuthProtoToNullSnmpAuthProto(r.AuthProto)
	s.SecLevel = snmpSecLevelToNullSnmpSecLevel(r.SecLevel)
	s.PrivProto = snmpPrivProtoToNullSnmpPrivProto(r.PrivProto)

	if r.AuthPass != nil {
		val, err := godevmandb.EncryptStrAes(*r.AuthPass, salt)
		if err != nil {
			return s, err
		}

		*s.AuthPass = val
	}

	if r.PrivPass != nil {
		val, err := godevmandb.EncryptStrAes(*r.PrivPass, salt)
		if err != nil {
			return s, err
		}

		*s.PrivPass = val
	}

	return s, nil
}

// Return corresponding godevmandb update parameters
func (r *snmpCredential) updateParams() (godevmandb.UpdateSnmpCredentialParams, error) {
	s := godevmandb.UpdateSnmpCredentialParams{}

	s.Variant = r.Variant
	s.Label = r.Label
	s.AuthName = r.AuthName
	s.AuthProto = snmpAuthProtoToNullSnmpAuthProto(r.AuthProto)
	s.SecLevel = snmpSecLevelToNullSnmpSecLevel(r.SecLevel)
	s.PrivProto = snmpPrivProtoToNullSnmpPrivProto(r.PrivProto)

	if r.AuthPass != nil {
		val, err := godevmandb.EncryptStrAes(*r.AuthPass, salt)
		if err != nil {
			return s, err
		}

		*s.AuthPass = val
	}

	if r.PrivPass != nil {
		val, err := godevmandb.EncryptStrAes(*r.PrivPass, salt)
		if err != nil {
			return s, err
		}

		*s.PrivPass = val
	}

	return s, nil
}

// Count SnmpCredentials
// @Summary Count snmp_credentials
// @Description Count number of snmp credentials
// @Tags config
// @ID count-snmp_credentials
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /config/snmp_credentials/count [GET]
func (h *Handler) CountSnmpCredentials(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountSnmpCredentials(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List SnmpCredentials
// @Summary List snmp_credentials
// @Description List snmp credentials info
// @Tags config
// @ID list-snmp_credentials
// @Param label_f query string false "url encoded SQL 'LIKE' operator pattern"
// @Param variant_f query string false "SQL '=' operator value"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} snmpCredential
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /config/snmp_credentials [GET]
func (h *Handler) GetSnmpCredentials(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetSnmpCredentialsParams{
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
	if v := r.FormValue("label_f"); v != "" {
		p.LabelF = v
	}

	if v := r.FormValue("variant_f"); v != "" {
		p.VariantF = v
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetSnmpCredentials(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := []snmpCredential{}
	for _, s := range res {
		r := snmpCredential{}
		r.getValues(s)
		out = append(out, r)
	}

	RespondJSON(w, r, http.StatusOK, out)
}

// Get SnmpCredential
// @Summary Get snmp_credential
// @Description Get snmp credential info
// @Tags config
// @ID get-snmp_credential
// @Param snmp_cred_id path string true "snmp_cred_id"
// @Success 200 {object} snmpCredential
// @Failure 400 {object} StatusResponse "Invalid snmp_cred_id"
// @Failure 404 {object} StatusResponse "Credential not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /config/snmp_credentials/{snmp_cred_id} [GET]
func (h *Handler) GetSnmpCredential(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "snmp_cred_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid credential ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetSnmpCredential(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "Credential not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	out := snmpCredential{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Create SnmpCredential
// @Summary Create snmp_credential
// @Description Create snmp credential
// @Tags config
// @ID create-snmp_credential
// @Param Body body snmpCredential true "JSON object of credential.<br />Ignored fields:<ul><li>snmp_cred_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 201 {object} snmpCredential
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /config/snmp_credentials [POST]
func (h *Handler) CreateSnmpCredential(w http.ResponseWriter, r *http.Request) {
	var pIn snmpCredential
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
	res, err := q.CreateSnmpCredential(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := snmpCredential{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusCreated, out)
}

// Update SnmpCredential
// @Summary Update snmp_credential
// @Description Update snmp credential
// @Tags config
// @ID update-snmp_credential
// @Param snmp_cred_id path string true "snmp_cred_id"
// @Param Body body snmpCredential true "JSON object of credential.<br />Ignored fields:<ul><li>snmp_cred_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 200 {object} snmpCredential
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /config/snmp_credentials/{snmp_cred_id} [PUT]
func (h *Handler) UpdateSnmpCredential(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "snmp_cred_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid credential ID")
		return
	}

	var pIn snmpCredential
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

	p.SnmpCredID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateSnmpCredential(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := snmpCredential{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Delete SnmpCredential
// @Summary Delete snmp_credential
// @Description Delete snmp credential
// @Tags config
// @ID delete-snmp_credential
// @Param snmp_cred_id path string true "snmp_cred_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid snmp_cred_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /config/snmp_credentials/{snmp_cred_id} [DELETE]
func (h *Handler) DeleteSnmpCredential(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "snmp_cred_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid credential ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteSnmpCredential(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Relations
// List SnmpCredential Main Devices
// @Summary List snmp_credential main devices
// @Description List snmp credential main devices info
// @Tags config
// @ID list-snmp-credential-devices
// @Param snmp_cred_id path string true "snmp_cred_id"
// @Success 200 {array} device
// @Failure 400 {object} StatusResponse "Invalid snmp_cred_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /config/snmp_credentials/{snmp_cred_id}/main_devices [GET]
func (h *Handler) GetSnmpCredentialsMainDevices(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "snmp_cred_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid credential ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetSnmpCredentialsMainDevices(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := []device{}
	for _, s := range res {
		a := device{}
		a.getValues(s)
		out = append(out, a)
	}

	RespondJSON(w, r, http.StatusOK, out)
}

// Relations
// List SnmpCredential Read Only Devices
// @Summary List snmp_credential ro devices
// @Description List snmp credential read only devices info
// @Tags config
// @ID list-snmp-credential-ro-devices
// @Param snmp_cred_id path string true "snmp_cred_id"
// @Success 200 {array} device
// @Failure 400 {object} StatusResponse "Invalid snmp_cred_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failed DB transaction"
// @Router /config/snmp_credentials/{snmp_cred_id}/ro_devices [GET]
func (h *Handler) GetSnmpCredentialsRoDevices(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "snmp_cred_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid credential ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetSnmpCredentialsRoDevices(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := []device{}
	for _, s := range res {
		a := device{}
		a.getValues(s)
		out = append(out, a)
	}

	RespondJSON(w, r, http.StatusOK, out)
}
