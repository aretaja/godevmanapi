package handlers

import (
	"time"

	"github.com/aretaja/godevmandb"
)

// JSON friendly local type to use in web api. Replaces sql.Null*/pgtype fields
type ospfNbr struct {
	NbrID     int64     `json:"nbr_id"`
	DevID     int64     `json:"dev_id"`
	NbrIp     *string   `json:"nbr_ip"`
	Condition *string   `json:"condition"`
	UpdatedOn time.Time `json:"updated_on"`
	CreatedOn time.Time `json:"created_on"`
}

// Import values from corresponding godevmandb struct
func (r *ospfNbr) getValues(s godevmandb.OspfNbr) {
	r.NbrID = s.NbrID
	r.DevID = s.DevID
	r.Condition = s.Condition
	r.UpdatedOn = s.UpdatedOn
	r.CreatedOn = s.CreatedOn
	r.NbrIp = pgInetToPtr(s.NbrIp)
}

// Return corresponding godevmandb create parameters
func (r *ospfNbr) createParams() godevmandb.CreateOspfNbrParams {
	s := godevmandb.CreateOspfNbrParams{}

	s.DevID = r.DevID
	s.Condition = r.Condition
	s.NbrIp = strToPgInet(r.NbrIp)

	return s
}

// Return corresponding godevmandb update parameters
func (r *ospfNbr) updateParams() godevmandb.UpdateOspfNbrParams {
	s := godevmandb.UpdateOspfNbrParams{}

	s.DevID = r.DevID
	s.Condition = r.Condition
	s.NbrIp = strToPgInet(r.NbrIp)

	return s
}
