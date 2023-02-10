package handlers

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aretaja/godevmandb"
	"github.com/go-chi/httplog"
	"github.com/jackc/pgtype"
)

// Pagination values
func paginateValues(r *http.Request) []*int32 {
	res := make([]*int32, 2)
	hlog := httplog.LogEntry(r.Context())

	l, err := strconv.ParseInt(r.FormValue("limit"), 10, 32)
	if err != nil {
		hlog.Debug().Msg("Parse limit - " + err.Error())
		hlog.Info().Msg("Invalid limit value. Using default")
	} else {
		if l < 1000 || l > 0 {
			lo := int32(l)
			res[0] = &lo
		} else {
			hlog.Info().Msg("Value of limit value out of range 0 - 1000")
		}
	}

	o, err := strconv.ParseInt(r.FormValue("offset"), 10, 32)
	if err != nil {
		hlog.Debug().Msg("Parse offset - " + err.Error())
		hlog.Info().Msg("Invalid offset value. Using default")
	} else {
		if o > 0 {
			oo := int32(o)
			res[1] = &oo
		} else {
			hlog.Info().Msg("Value of offset value out of range > 0")
		}
	}

	return res
}

// Time filter
func parseTimeFilter(r *http.Request) []time.Time {
	res := make([]time.Time, 4)
	hlog := httplog.LogEntry(r.Context())
	keys := []string{"updated_ge", "updated_le", "created_ge", "created_le"}

	for i := 0; i < 4; i++ {
		v := r.FormValue(keys[i])
		if v == "" {
			continue
		}
		uts, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			hlog.Debug().Msg("Parse " + keys[i] + " - " + err.Error())
		} else {
			res[i] = time.UnixMilli(uts)
		}
	}

	return res
}

// IP/CIDR string pointer to pgtype.Inet converter
func strToPgInet(p *string) pgtype.Inet {
	n := net.IPNet{}
	r := pgtype.Inet{IPNet: &n, Status: pgtype.Null}

	if p != nil {
		// Add CIDR network part if missing
		if ok := strings.Contains(*p, "/"); !ok {
			if strings.Contains(*p, ".") {
				*p += "/32"
			} else {
				*p += "/128"
			}
		}

		if _, ip, err := net.ParseCIDR(*p); err == nil {
			r.IPNet = ip
			r.Status = pgtype.Present
		}
	}
	return r
}

// MAC string pointer to pgtype.Macaddr converter
func strToPgMacaddr(p *string) pgtype.Macaddr {
	r := pgtype.Macaddr{Status: pgtype.Null}

	if p != nil {
		if mac, err := net.ParseMAC(*p); err == nil {
			r.Addr = mac
			r.Status = pgtype.Present
		}
	}
	return r
}

// String pointer to sql.NullString converter
func strToNullString(p *string) sql.NullString {
	r := sql.NullString{}

	if p != nil {
		r = sql.NullString{
			String: *p,
			Valid:  true,
		}
	}
	return r
}

// Int64 to sql.NullInt16 converter
func int64ToNullInt16(p *int64) sql.NullInt16 {
	r := sql.NullInt16{}

	if p != nil {
		r = sql.NullInt16{
			Int16: int16(*p),
			Valid: true,
		}
	}
	return r
}

// Int64 to sql.NullInt64 converter
func int64ToNullInt64(p *int64) sql.NullInt64 {
	r := sql.NullInt64{}

	if p != nil {
		r = sql.NullInt64{
			Int64: *p,
			Valid: true,
		}
	}
	return r
}

// godevmandb.SnmpAuthProto pointer to godevmandb.NullSnmpAuthProto converter
func snmpAuthProtoToNullSnmpAuthProto(p *godevmandb.SnmpAuthProto) godevmandb.NullSnmpAuthProto {
	r := godevmandb.NullSnmpAuthProto{}

	if p != nil {
		r = godevmandb.NullSnmpAuthProto{
			SnmpAuthProto: *p,
			Valid:         true,
		}
	}
	return r
}

// godevmandb.SnmpPrivProto pointer to godevmandb.NullSnmpPrivProto converter
func snmpPrivProtoToNullSnmpPrivProto(p *godevmandb.SnmpPrivProto) godevmandb.NullSnmpPrivProto {
	r := godevmandb.NullSnmpPrivProto{}

	if p != nil {
		r = godevmandb.NullSnmpPrivProto{
			SnmpPrivProto: *p,
			Valid:         true,
		}
	}
	return r
}

// godevmandb.SnmpSecLevel pointer to godevmandb.NullSnmpSecLevel converter
func snmpSecLevelToNullSnmpSecLevel(p *godevmandb.SnmpSecLevel) godevmandb.NullSnmpSecLevel {
	r := godevmandb.NullSnmpSecLevel{}

	if p != nil {
		r = godevmandb.NullSnmpSecLevel{
			SnmpSecLevel: *p,
			Valid:        true,
		}
	}
	return r
}

// sql.NullString to string pointer converter
func nullStringToPtr(n sql.NullString) *string {
	if n.Valid {
		if v, err := n.Value(); err == nil {
			if res, ok := v.(string); ok {
				return &res
			}
		}
	}
	return nil
}

// sql.NullInt16 to int64 pointer converter
func nullInt16ToPtr(n sql.NullInt16) *int64 {
	if n.Valid {
		if v, err := n.Value(); err == nil {
			if res, ok := v.(int64); ok {
				return &res
			}
		}
	}
	return nil
}

// sql.NullInt64 to int64 pointer converter
func nullInt64ToPtr(n sql.NullInt64) *int64 {
	if n.Valid {
		if v, err := n.Value(); err == nil {
			if res, ok := v.(int64); ok {
				return &res
			}
		}
	}
	return nil
}

// pgtype.Inet to string pointer converter
func pgInetToPtr(n pgtype.Inet) *string {
	if n.Status == pgtype.Present {
		if v, err := n.Value(); err == nil {
			res := fmt.Sprintf("%s", v)
			return &res
		}
	}
	return nil
}

// pgtype.Macaddr to string pointer converter
func pgMacaddrToPtr(n pgtype.Macaddr) *string {
	if n.Status == pgtype.Present {
		if v, err := n.Value(); err == nil {
			res := fmt.Sprintf("%s", v)
			return &res
		}
	}
	return nil
}

// godevmandb.NullSnmpAuthProto to SnmpAuthProto pointer converter
func nullSnmpAuthProtoToPtr(n godevmandb.NullSnmpAuthProto) *godevmandb.SnmpAuthProto {
	if n.Valid {
		return &n.SnmpAuthProto
	}
	return nil
}

// godevmandb.NullSnmpPrivProto to SnmpPrivProto pointer converter
func nullSnmpPrivProtoToPtr(n godevmandb.NullSnmpPrivProto) *godevmandb.SnmpPrivProto {
	if n.Valid {
		return &n.SnmpPrivProto
	}
	return nil
}

// godevmandb.NullSnmpSecLevel to SnmpSecLevel pointer converter
func nullSnmpSecLevelToPtr(n godevmandb.NullSnmpSecLevel) *godevmandb.SnmpSecLevel {
	if n.Valid {
		return &n.SnmpSecLevel
	}
	return nil
}
