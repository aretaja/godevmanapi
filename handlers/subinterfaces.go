package handlers

import (
	"time"

	"github.com/aretaja/godevmandb"
)

// JSON friendly local type to use in web api. Replaces sql.Null*/pgtype fields
type subinterface struct {
	UpdatedOn time.Time `json:"updated_on"`
	CreatedOn time.Time `json:"created_on"`
	Mac       *string   `json:"mac"`
	Alias     *string   `json:"alias"`
	Oper      *int16    `json:"oper"`
	Adm       *int16    `json:"adm"`
	Speed     *int64    `json:"speed"`
	TypeEnum  *string   `json:"type_enum"`
	Notes     *string   `json:"notes"`
	Ifindex   *int64    `json:"ifindex"`
	IfID      *int64    `json:"if_id"`
	Descr     string    `json:"descr"`
	SifID     int64     `json:"sif_id"`
}

// Import values from corresponding godevmandb struct
func (r *subinterface) getValues(s godevmandb.Subinterface) {
	r.SifID = s.SifID
	r.IfID = s.IfID
	r.Ifindex = s.Ifindex
	r.Descr = s.Descr
	r.Alias = s.Alias
	r.Oper = s.Oper
	r.Adm = s.Adm
	r.Speed = s.Speed
	r.TypeEnum = s.TypeEnum
	r.Notes = s.Notes
	r.UpdatedOn = s.UpdatedOn
	r.CreatedOn = s.CreatedOn
	r.Mac = pgMacaddrToPtr(s.Mac)
}

// Return corresponding godevmandb create parameters
func (r *subinterface) createParams() godevmandb.CreateSubinterfaceParams {
	s := godevmandb.CreateSubinterfaceParams{}

	s.IfID = r.IfID
	s.Ifindex = r.Ifindex
	s.Descr = r.Descr
	s.Alias = r.Alias
	s.Oper = r.Oper
	s.Adm = r.Adm
	s.Speed = r.Speed
	s.TypeEnum = r.TypeEnum
	s.Notes = r.Notes
	s.Mac = strToPgMacaddr(r.Mac)

	return s
}

// Return corresponding godevmandb update parameters
func (r *subinterface) updateParams() godevmandb.UpdateSubinterfaceParams {
	s := godevmandb.UpdateSubinterfaceParams{}

	s.IfID = r.IfID
	s.Ifindex = r.Ifindex
	s.Descr = r.Descr
	s.Alias = r.Alias
	s.Oper = r.Oper
	s.Adm = r.Adm
	s.Speed = r.Speed
	s.TypeEnum = r.TypeEnum
	s.Notes = r.Notes
	s.Mac = strToPgMacaddr(r.Mac)

	return s
}
