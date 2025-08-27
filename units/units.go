/*
Units Module - Unit Conversion System
=====================================

This module implements unit conversion functionality for three measurement
categories: length, weight, and time. All conversions are performed through
a base unit system where each unit is defined by its conversion factor to
the base unit of its category.

Conversion methodology:
1. Convert input value from source unit to base unit
2. Convert from base unit to target unit
3. Formula: result = value * (source_factor / target_factor)

Base units:
- Length: meters (m)
- Weight: kilograms (kg)
- Time: seconds (s)

The system prevents cross-category conversions (e.g., meters to kilograms)
by checking unit category membership before performing calculations.
*/

package units

import (
	"fmt"
)

// lengthUnits defines conversion factors to meters (base unit)
var lengthUnits = map[string]float64{
	"m":  1,       // meters (base unit)
	"cm": 0.01,    // centimeters
	"mm": 0.001,   // millimeters
	"km": 1000,    // kilometers
	"in": 0.0254,  // inches
	"ft": 0.3048,  // feet
	"yd": 0.9144,  // yards
	"mi": 1609.34, // miles
}

// weightUnits defines conversion factors to kilograms (base unit)
var weightUnits = map[string]float64{
	"kg":  1,         // kilograms (base unit)
	"g":   0.001,     // grams
	"mg":  0.000001,  // milligrams
	"lb":  0.453592,  // pounds
	"oz":  0.0283495, // ounces
	"ton": 1000,      // metric tons
}

// timeUnits defines conversion factors to seconds (base unit)
var timeUnits = map[string]float64{
	"s":   1,     // seconds (base unit)
	"ms":  0.001, // milliseconds
	"min": 60,    // minutes
	"h":   3600,  // hours
	"d":   86400, // days
}

// Convert performs unit conversion between compatible units
// Returns converted value or error for unsupported/incompatible units
func Convert(value float64, from, to string) (float64, error) {
	// Attempt length unit conversion
	if sourceFactor, sourceExists := lengthUnits[from]; sourceExists {
		if targetFactor, targetExists := lengthUnits[to]; targetExists {
			// Both units are length units - perform conversion
			// Formula: value * (source_to_base_ratio / target_to_base_ratio)
			return value * sourceFactor / targetFactor, nil
		}
	}

	// Attempt weight unit conversion
	if sourceFactor, sourceExists := weightUnits[from]; sourceExists {
		if targetFactor, targetExists := weightUnits[to]; targetExists {
			// Both units are weight units - perform conversion
			return value * sourceFactor / targetFactor, nil
		}
	}

	// Attempt time unit conversion
	if sourceFactor, sourceExists := timeUnits[from]; sourceExists {
		if targetFactor, targetExists := timeUnits[to]; targetExists {
			// Both units are time units - perform conversion
			return value * sourceFactor / targetFactor, nil
		}
	}

	// No valid conversion found - units are either unknown or incompatible
	return 0, fmt.Errorf("unsupported conversion: %s → %s", from, to)
}
