package handlers

import (
	"time"

	"github.com/aretaja/godevmandb"
)

// JSON friendly local type to use in web api. Replaces sql.Null*/pgtype fields
type iface struct {
	UpdatedOn  time.Time `json:"updated_on"`
	CreatedOn  time.Time `json:"created_on"`
	Adm        *int16    `json:"adm"`
	Mac        *string   `json:"mac"`
	ConID      *int64    `json:"con_id"`
	EntID      *int64    `json:"ent_id"`
	Ifindex    *int64    `json:"ifindex"`
	OtnIfID    *int64    `json:"otn_if_id"`
	Alias      *string   `json:"alias"`
	Oper       *int16    `json:"oper"`
	Parent     *int64    `json:"parent"`
	Speed      *int64    `json:"speed"`
	Minspeed   *int64    `json:"minspeed"`
	TypeEnum   *int16    `json:"type_enum"`
	Descr      string    `json:"descr"`
	IfID       int64     `json:"if_id"`
	DevID      int64     `json:"dev_id"`
	Monstatus  int16     `json:"monstatus"`
	Monerrors  int16     `json:"monerrors"`
	Monload    int16     `json:"monload"`
	Montraffic int16     `json:"montraffic"`
}

// Import values from corresponding godevmandb struct
func (r *iface) getValues(s godevmandb.Interface) {
	r.IfID = s.IfID
	r.ConID = s.ConID
	r.OtnIfID = s.OtnIfID
	r.DevID = s.DevID
	r.EntID = s.EntID
	r.Ifindex = s.Ifindex
	r.Descr = s.Descr
	r.Alias = s.Alias
	r.Oper = s.Oper
	r.Adm = s.Adm
	r.Speed = s.Speed
	r.Minspeed = s.Minspeed
	r.TypeEnum = s.TypeEnum
	r.Monstatus = s.Monstatus
	r.Monerrors = s.Monerrors
	r.Monload = s.Monload
	r.Montraffic = s.Montraffic
	r.UpdatedOn = s.UpdatedOn
	r.CreatedOn = s.CreatedOn
	r.Mac = pgMacaddrToPtr(s.Mac)
}

// Return corresponding godevmandb create parameters
func (r *iface) createParams() godevmandb.CreateInterfaceParams {
	s := godevmandb.CreateInterfaceParams{}

	s.ConID = r.ConID
	s.OtnIfID = r.OtnIfID
	s.DevID = r.DevID
	s.EntID = r.EntID
	s.Ifindex = r.Ifindex
	s.Descr = r.Descr
	s.Alias = r.Alias
	s.Oper = r.Oper
	s.Adm = r.Adm
	s.Speed = r.Speed
	s.Minspeed = r.Minspeed
	s.TypeEnum = r.TypeEnum
	s.Monstatus = r.Monstatus
	s.Monerrors = r.Monerrors
	s.Monload = r.Monload
	s.Montraffic = r.Montraffic
	s.Mac = strToPgMacaddr(r.Mac)

	return s
}

// Return corresponding godevmandb update parameters
func (r *iface) updateParams() godevmandb.UpdateInterfaceParams {
	s := godevmandb.UpdateInterfaceParams{}

	s.ConID = r.ConID
	s.OtnIfID = r.OtnIfID
	s.DevID = r.DevID
	s.EntID = r.EntID
	s.Ifindex = r.Ifindex
	s.Descr = r.Descr
	s.Alias = r.Alias
	s.Oper = r.Oper
	s.Adm = r.Adm
	s.Speed = r.Speed
	s.Minspeed = r.Minspeed
	s.TypeEnum = r.TypeEnum
	s.Monstatus = r.Monstatus
	s.Monerrors = r.Monerrors
	s.Monload = r.Monload
	s.Montraffic = r.Montraffic
	s.Mac = strToPgMacaddr(r.Mac)

	return s
}
