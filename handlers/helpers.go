package handlers

import (
	"database/sql"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

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

// IP/CIDR string to pgtype.Inet parser
func strToPgInet(s string) pgtype.Inet {
	r := pgtype.Inet{
		Status: pgtype.Null,
	}

	if s != "" {
		// Add CIDR network part if missing
		if ok := strings.Contains(s, "/"); !ok {
			if strings.Contains(s, ".") {
				s += "/32"
			} else {
				s += "/128"
			}
		}

		if _, ip, err := net.ParseCIDR(s); err == nil {
			r.IPNet = ip
			r.Status = pgtype.Present
		}
	}
	return r
}

// MAC string to pgtype.Macaddr parser
func strToPgMacaddr(s string) pgtype.Macaddr {
	r := pgtype.Macaddr{
		Status: pgtype.Null,
	}

	if s != "" {
		if mac, err := net.ParseMAC(s); err == nil {
			r.Addr = mac
			r.Status = pgtype.Present
		}
	}
	return r
}

// String to sql.NullString parser
func strToNullString(s string) sql.NullString {
	r := sql.NullString{}

	if s != "" {
		r = sql.NullString{
			String: s,
			Valid:  true,
		}
	}
	return r
}
