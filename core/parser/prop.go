package prop

import (
	"fmt"
)


type Prop struct{
	prop map["string"][]string
}


func Validate(properties *Properties) (int, error) {
	var valid = map[string]bool{"AB": true, "AE": true, "AN": true, "AP": true, "AR": true, "AW": true, "B": true, "BL": true, "BM": true, "BR": true, "BT": true, "C": true, "CA": true, "CP": true, "CR": true, "DD": true, "DM": true, "DO": true, "DT": true, "EV": true, "FF": true, "FG": true, "GB": true, "GC": true, "GM": true, "GN": true, "GW": true, "HA": true, "HO": true, "IT": true, "KM": true, "KO": true, "LB": true, "LN": true, "MA": true, "MN": true, "N": true, "OB": true, "ON": true, "OT": true, "OW": true, "PB": true, "PC": true, "PL": true, "PM": true, "PW": true, "RE": true, "RO": true, "RU": true, "SL": true, "SO": true, "SQ": true, "ST": true, "SZ": true, "TB": true, "TE": true, "TM": true, "TR": true, "TW": true, "UC": true, "US": true, "V": true, "VW": true, "W": true, "WL": true, "WR": true, "WT": true,}
	for key, _ := range properties.Prop {
		_, ok := valid[key]
		if !ok {
	 		return 0, fmt.Errorf("Properties: %s is an invalid property.", key)
		}
	}
	return 1, nil
}

