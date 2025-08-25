package units

import (
	"fmt"
)

var lengthUnits = map[string]float64{
	"m":  1,
	"cm": 0.01,
	"mm": 0.001,
	"km": 1000,
	"in": 0.0254,
	"ft": 0.3048,
	"yd": 0.9144,
	"mi": 1609.34,
}

var weightUnits = map[string]float64{
	"kg":  1,
	"g":   0.001,
	"mg":  0.000001,
	"lb":  0.453592,
	"oz":  0.0283495,
	"ton": 1000,
}

var timeUnits = map[string]float64{
	"s":   1,
	"ms":  0.001,
	"min": 60,
	"h":   3600,
	"d":   86400,
}

func Convert(value float64, from, to string) (float64, error) {
	if factor, ok := lengthUnits[from]; ok {
		if target, ok := lengthUnits[to]; ok {
			return value * factor / target, nil
		}
	}
	if factor, ok := weightUnits[from]; ok {
		if target, ok := weightUnits[to]; ok {
			return value * factor / target, nil
		}
	}
	if factor, ok := timeUnits[from]; ok {
		if target, ok := timeUnits[to]; ok {
			return value * factor / target, nil
		}
	}
	return 0, fmt.Errorf("unsupported conversion: %s → %s", from, to)
}
