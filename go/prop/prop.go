// Package prop adds methods for handling SGF properties, including validation.
package prop

import "errors"

// Prop is used to store a label parsed from SGF
type Prop string

// ValidProperties lists all valid SGF properties
var ValidProperties = map[Prop]bool{"AB": true, "AE": true, "AN": true, "AP": true, "AR": true, "AW": true, "B": true, "BL": true, "BM": true, "BR": true,
	"BT": true, "C": true, "CA": true, "CP": true, "CR": true, "DD": true, "DM": true, "DO": true, "DT": true, "EV": true, "FF": true, "FG": true, "GB": true,
	"GC": true, "GM": true, "GN": true, "GW": true, "HA": true, "HO": true, "IT": true, "KM": true, "KO": true, "LB": true, "LN": true, "MA": true, "MN": true,
	"N": true, "OB": true, "ON": true, "OT": true, "OW": true, "PB": true, "PC": true, "PL": true, "PM": true, "PW": true, "RE": true, "RO": true, "RU": true,
	"SL": true, "SO": true, "SQ": true, "ST": true, "SZ": true, "TB": true, "TE": true, "TM": true, "TR": true, "TW": true, "UC": true, "US": true, "V": true,
	"VW": true, "W": true, "WL": true, "WR": true, "WT": true}

// Validate returns true if Prop is an accepted SGF property.
func Validate(prop Prop) bool {
	_, ok := ValidProperties[prop]
	return ok
}

var ErrConvertingProp = errors.New("error converting property")
