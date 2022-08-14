package prop

import (
	"fmt"
	"sort"
	"strings"

	"github.com/otrego/clamshell/go/movetree"
)

// ProcessPropertyData uses converters to process property data.
func ProcessPropertyData(n *movetree.Node, p string, propData []string) error {
	if !HasConverter(p) {
		// For properties without an explicit converter, add to unprocessed
		// Properties.
		n.SGFProperties[p] = propData
		return nil
	}
	conv := Converter(p)
	if conv.Scope == RootScope && (n.MoveNum() != 0 || n.VarNum() != 0) {
		return fmt.Errorf("%w: for property %s: property is a root-node only property, but was found at {move:%d, variation: %d}",
			ErrConvertingProp, p, n.MoveNum(), n.VarNum())
	}

	if err := conv.From(n, p, propData); err != nil {
		return err
	}
	return nil
}

// Scope indicates the scope for a property.
type Scope string

const (
	// RootScope indicates a property that only applies to the root node.
	RootScope Scope = "RootScope"

	// AllScope indicates a property that applies to all nodes.
	AllScope Scope = "AllScope"
)

// FromSGF converts an SGF Property to node property
type FromSGF func(node *movetree.Node, prop string, values []string) error

// ToSGF converts an Node property to an SGF property list.
type ToSGF func(node *movetree.Node) (string, error)

// A SGFConverter converts SGF properties to / from node properties.
type SGFConverter struct {
	// Props is the name of the SGF properties that apply to this
	// property-converter.
	// Ex: {"AW", "AB"}
	Props []Prop
	// Scope indicates what the property-converter applies to (Root, All)
	Scope Scope
	// From converts from SGF data
	From FromSGF
	// To converts to SGF data
	To ToSGF
}

// HasConverter indicates whether there's a known SGF Property converter.
func HasConverter(prop string) bool {
	_, ok := propToConv[Prop(prop)]
	return ok
}

// Converter gets a property converter for converting to/from SGf, returning nil
// if no property converter can be found.
func Converter(prop string) *SGFConverter {
	return propToConv[Prop(prop)]
}

// ConvertNode converts all the properties in a node
func ConvertNode(n *movetree.Node) (string, error) {
	var sb strings.Builder
	for _, c := range converters {
		if c.Scope == RootScope && n.MoveNum() != 0 {
			// skip non-root-scoped properties for non-root nodes.
			continue
		}
		s, err := c.To(n)
		if err != nil {
			return "", err
		}
		sb.WriteString(s)
	}

	var keys []string
	for key := range n.SGFProperties {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		sb.WriteString(key)
		for _, value := range n.SGFProperties[key] {
			sb.WriteString("[" + value + "]")
		}
	}
	return sb.String(), nil
}
