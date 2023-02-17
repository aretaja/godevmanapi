package handlers

import (
	"time"

	"github.com/aretaja/godevmandb"
)

// JSON friendly local type to use in web api. Replaces sql.Null*/pgtype fields
type ipInterface struct {
	UpdatedOn time.Time `json:"updated_on"`
	CreatedOn time.Time `json:"created_on"`
	Ifindex   *int64    `json:"ifindex"`
	IpAddr    *string   `json:"ip_addr"`
	Descr     *string   `json:"descr"`
	Alias     *string   `json:"alias"`
	IpID      int64     `json:"ip_id"`
	DevID     int64     `json:"dev_id"`
}

// Import values from corresponding godevmandb struct
func (r *ipInterface) getValues(s godevmandb.IpInterface) {
	r.IpID = s.IpID
	r.DevID = s.DevID
	r.Ifindex = s.Ifindex
	r.Descr = s.Descr
	r.Alias = s.Alias
	r.UpdatedOn = s.UpdatedOn
	r.CreatedOn = s.CreatedOn
	r.IpAddr = pgInetToPtr(s.IpAddr)
}

// Return corresponding godevmandb create parameters
func (r *ipInterface) createParams() godevmandb.CreateIpInterfaceParams {
	s := godevmandb.CreateIpInterfaceParams{}

	s.DevID = r.DevID
	s.Ifindex = r.Ifindex
	s.Descr = r.Descr
	s.Alias = r.Alias
	s.IpAddr = strToPgInet(r.IpAddr)

	return s
}

// Return corresponding godevmandb update parameters
func (r *ipInterface) updateParams() godevmandb.UpdateIpInterfaceParams {
	s := godevmandb.UpdateIpInterfaceParams{}

	s.DevID = r.DevID
	s.Ifindex = r.Ifindex
	s.Descr = r.Descr
	s.Alias = r.Alias
	s.IpAddr = strToPgInet(r.IpAddr)

	return s
}
