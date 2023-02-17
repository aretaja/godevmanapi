package handlers

import (
	"time"

	"github.com/aretaja/godevmandb"
)

// JSON friendly local type to use in web api. Replaces sql.Null*/pgtype fields
type xconnect struct {
	UpdatedOn   time.Time `json:"updated_on"`
	CreatedOn   time.Time `json:"created_on"`
	PeerIp      *string   `json:"peer_ip"`
	IfID        *int64    `json:"if_id"`
	PeerIfalias *string   `json:"peer_ifalias"`
	Xname       *string   `json:"xname"`
	Descr       *string   `json:"descr"`
	OpStat      *string   `json:"op_stat"`
	OpStatIn    *string   `json:"op_stat_in"`
	OpStatOut   *string   `json:"op_stat_out"`
	PeerDevID   *int64    `json:"peer_dev_id"`
	VcIdx       int64     `json:"vc_idx"`
	VcID        int64     `json:"vc_id"`
	XcID        int64     `json:"xc_id"`
	DevID       int64     `json:"dev_id"`
}

// Import values from corresponding godevmandb struct
func (r *xconnect) getValues(s godevmandb.Xconnect) {
	r.XcID = s.XcID
	r.DevID = s.DevID
	r.PeerDevID = s.PeerDevID
	r.IfID = s.IfID
	r.VcIdx = s.VcIdx
	r.VcID = s.VcID
	r.PeerIfalias = s.PeerIfalias
	r.Xname = s.Xname
	r.Descr = s.Descr
	r.OpStat = s.OpStat
	r.OpStatIn = s.OpStatIn
	r.OpStatOut = s.OpStatOut
	r.UpdatedOn = s.UpdatedOn
	r.CreatedOn = s.CreatedOn
	r.PeerIp = pgInetToPtr(s.PeerIp)
}

// Return corresponding godevmandb create parameters
func (r *xconnect) createParams() godevmandb.CreateXconnectParams {
	s := godevmandb.CreateXconnectParams{}

	s.DevID = r.DevID
	s.PeerDevID = r.PeerDevID
	s.IfID = r.IfID
	s.VcIdx = r.VcIdx
	s.VcID = r.VcID
	s.PeerIfalias = r.PeerIfalias
	s.Xname = r.Xname
	s.Descr = r.Descr
	s.OpStat = r.OpStat
	s.OpStatIn = r.OpStatIn
	s.OpStatOut = r.OpStatOut
	s.PeerIp = strToPgInet(r.PeerIp)

	return s
}

// Return corresponding godevmandb update parameters
func (r *xconnect) updateParams() godevmandb.UpdateXconnectParams {
	s := godevmandb.UpdateXconnectParams{}

	s.DevID = r.DevID
	s.PeerDevID = r.PeerDevID
	s.IfID = r.IfID
	s.VcIdx = r.VcIdx
	s.VcID = r.VcID
	s.PeerIfalias = r.PeerIfalias
	s.Xname = r.Xname
	s.Descr = r.Descr
	s.OpStat = r.OpStat
	s.OpStatIn = r.OpStatIn
	s.OpStatOut = r.OpStatOut
	s.PeerIp = strToPgInet(r.PeerIp)

	return s
}
