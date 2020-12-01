// Package prop adds methods for handling SGF properties, including validation.
package prop

import (
	"strings"

	"github.com/otrego/clamshell/core/movetree"
)

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

// FromSGF converts an SGF Property to node property
type FromSGF func(node *movetree.Node, values []string) error

// ToSGF converts an Node property to an SGF property list.
type ToSGF func(node *movetree.Node) ([]string, error)

// A Converter converts SGF properties to / from node properties.
type Converter struct {
	// Prop is the name of the SGF property. Ex: "AW"
	Prop Prop
	// From converts from SGF data
	From FromSGF
	// To converts to SGF data
	To ToSGF
}

// BuildSGFString adds SGF content to a string builder.
func (c *Converter) BuildSGFString(node *movetree.Node, sb *strings.Builder) error {
	data, err := c.To(node)
	if err != nil {
		return err
	}
	sb.WriteString(string(c.Prop))
	if data != nil && len(data) == 0 {
		// Special case (mostly for passes)
		sb.WriteString("[]")
	}
	for _, d := range data {
		sb.WriteString("[" + d + "]")
	}
	return nil
}
