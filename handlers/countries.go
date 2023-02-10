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
type country struct {
	UpdatedOn time.Time `json:"updated_on"`
	CreatedOn time.Time `json:"created_on"`
	Code      string    `json:"code"`
	Descr     string    `json:"descr"`
	CountryID int64     `json:"country_id"`
}

// Import values from corresponding godevmandb struct
func (r *country) getValues(s godevmandb.Country) {
	r.CountryID = s.CountryID
	r.Descr = s.Descr
	r.Code = s.Code
	r.UpdatedOn = s.UpdatedOn
	r.CreatedOn = s.CreatedOn
}

// Return corresponding godevmandb create parameters
func (r *country) createParams() godevmandb.CreateCountryParams {
	s := godevmandb.CreateCountryParams{}

	s.Descr = r.Descr
	s.Code = r.Code

	return s
}

// Return corresponding godevmandb update parameters
func (r *country) updateParams() godevmandb.UpdateCountryParams {
	s := godevmandb.UpdateCountryParams{}

	s.Descr = r.Descr
	s.Code = r.Code

	return s
}

// Count Countries
// @Summary Count countries
// @Description Count number of countries
// @Tags sites
// @ID count-countries
// @Success 200 {object} CountResponse
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /sites/countries/count [GET]
func (h *Handler) CountCountries(w http.ResponseWriter, r *http.Request) {
	q := godevmandb.New(h.db)
	res, err := q.CountCountries(h.ctx)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, r, http.StatusOK, CountResponse{Count: res})
}

// List countries
// @Summary List countries
// @Description List countries info
// @Tags sites
// @ID list-countries
// @Param descr_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param code_f query string false "url encoded SQL 'ILIKE' operator pattern"
// @Param limit query int false "min: 1; max: 1000; default: 100"
// @Param offset query int false "default: 0"
// @Param updated_ge query int false "record update time >= (unix timestamp in milliseconds)"
// @Param updated_le query int false "record update time <= (unix timestamp in milliseconds)"
// @Param created_ge query int false "record creation time >= (unix timestamp in milliseconds)"
// @Param created_le query int false "record creation time <= (unix timestamp in milliseconds)"
// @Success 200 {array} country
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /sites/countries [GET]
func (h *Handler) GetCountries(w http.ResponseWriter, r *http.Request) {
	// Pagination
	var p = godevmandb.GetCountriesParams{
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
	d := r.FormValue("descr_f")
	if d != "" {
		p.DescrF = d
	}

	// Code filter
	c := r.FormValue("code_f")
	if c != "" {
		p.CodeF = c
	}

	// Query DB
	q := godevmandb.New(h.db)
	res, err := q.GetCountries(h.ctx, p)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := []country{}
	for _, s := range res {
		r := country{}
		r.getValues(s)
		out = append(out, r)
	}

	RespondJSON(w, r, http.StatusOK, out)
}

// Get Country
// @Summary Get country
// @Description Get country info
// @Tags sites
// @ID get-country
// @Param country_id path string true "country_id"
// @Success 200 {object} country
// @Failure 400 {object} StatusResponse "Invalid country_id"
// @Failure 404 {object} StatusResponse "Country not found"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /sites/countries/{country_id} [GET]
func (h *Handler) GetCountry(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "country_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid country ID")
		return
	}

	q := godevmandb.New(h.db)
	res, err := q.GetCountry(h.ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			RespondError(w, r, http.StatusNotFound, "Country not found")
		} else {
			RespondError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	out := country{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Create Country
// @Summary Create country
// @Description Create country
// @Tags sites
// @ID create-country
// @Param Body body country true "JSON object of country.<br />Ignored fields:<ul><li>country_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 201 {object} country
// @Failure 400 {object} StatusResponse "Invalid request payload"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /sites/countries [POST]
func (h *Handler) CreateCountry(w http.ResponseWriter, r *http.Request) {
	var pIn country
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Create parameters for new db record
	p := pIn.createParams()

	q := godevmandb.New(h.db)
	res, err := q.CreateCountry(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := country{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusCreated, out)
}

// Update Country
// @Summary Update country
// @Description Update country
// @Tags sites
// @ID update-country
// @Param country_id path string true "country_id"
// @Param Body body country true "JSON object of country.<br />Ignored fields:<ul><li>country_id</li><li>updated_on</li><li>created_on</li></ul>"
// @Success 200 {object} country
// @Failure 400 {object} StatusResponse "Invalid request"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /sites/countries/{country_id} [PUT]
func (h *Handler) UpdateCountry(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "country_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid country ID")
		return
	}

	var pIn country
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pIn); err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Update parameters for new db record
	p := pIn.updateParams()
	p.CountryID = id

	q := godevmandb.New(h.db)
	res, err := q.UpdateCountry(h.ctx, p)

	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	out := country{}
	out.getValues(res)

	RespondJSON(w, r, http.StatusOK, out)
}

// Delete Country
// @Summary Delete country
// @Description Delete country
// @Tags sites
// @ID delete-country
// @Param country_id path string true "country_id"
// @Success 204
// @Failure 400 {object} StatusResponse "Invalid country_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /sites/countries/{country_id} [DELETE]
func (h *Handler) DeleteCountry(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "country_id"), 10, 64)
	if err != nil {
		RespondError(w, r, http.StatusBadRequest, "Invalid country ID")
		return
	}

	q := godevmandb.New(h.db)
	err = q.DeleteCountry(h.ctx, id)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Relations
// List Country Sites
// @Summary List country sites
// @Description List country sites info
// @Tags sites
// @ID list-country-sites
// @Param country_id path string true "country_id"
// @Success 200 {array} connection
// @Failure 400 {object} StatusResponse "Invalid country_id"
// @Failure 404 {object} StatusResponse "Invalid route error"
// @Failure 405 {object} StatusResponse "Invalid method error"
// @Failure 500 {object} StatusResponse "Failde DB transaction"
// @Router /sites/countries/{country_id}/sites [GET]
// func (h *Handler) GetCountrySites(w http.ResponseWriter, r *http.Request) {
// 	id, err := strconv.ParseInt(chi.URLParam(r, "country_id"), 10, 64)
// 	if err != nil {
// 		RespondError(w, r, http.StatusBadRequest, "Invalid country ID")
// 		return
// 	}

// 	q := godevmandb.New(h.db)
// 	res, err := q.GetCountrySites(h.ctx, id)
// 	if err != nil {
// 		RespondError(w, r, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	out := []site{}
// 	for _, s := range res {
// 		a := site{}
// 		a.getValues(s)
// 		out = append(out, a)
// 	}

// 	RespondJSON(w, r, http.StatusOK, out)
// }
