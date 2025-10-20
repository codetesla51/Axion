package units

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert_Length(t *testing.T) {
	tests := []struct {
		from, to string
		value    float64
		expected float64
	}{
		{"m", "cm", 1, 100},
		{"cm", "m", 100, 1},
		{"km", "m", 2, 2000},
		{"mi", "km", 1, 1.60934},
		{"ft", "in", 3, 36},
		{"yd", "m", 1, 0.9144},
		{"mm", "cm", 10, 1},
	}

	for _, tt := range tests {
		t.Run(tt.from+"→"+tt.to, func(t *testing.T) {
			got, err := Convert(tt.value, tt.from, tt.to)
			assert.NoError(t, err)
			assert.InDelta(t, tt.expected, got, 1e-6)
		})
	}
}

func TestConvert_Weight(t *testing.T) {
	tests := []struct {
		from, to string
		value    float64
		expected float64
	}{
		{"kg", "g", 1, 1000},
		{"g", "kg", 1000, 1},
		{"mg", "g", 1000, 1},
		{"lb", "kg", 1, 0.453592},
		{"oz", "g", 1, 28.3495},
		{"ton", "kg", 1, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.from+"→"+tt.to, func(t *testing.T) {
			got, err := Convert(tt.value, tt.from, tt.to)
			assert.NoError(t, err)
			assert.InDelta(t, tt.expected, got, 1e-6)
		})
	}
}

func TestConvert_Time(t *testing.T) {
	tests := []struct {
		from, to string
		value    float64
		expected float64
	}{
		{"s", "ms", 1, 1000},
		{"ms", "s", 1000, 1},
		{"min", "s", 1, 60},
		{"h", "s", 1, 3600},
		{"d", "h", 1, 24},
	}

	for _, tt := range tests {
		t.Run(tt.from+"→"+tt.to, func(t *testing.T) {
			got, err := Convert(tt.value, tt.from, tt.to)
			assert.NoError(t, err)
			assert.InDelta(t, tt.expected, got, 1e-6)
		})
	}
}

func TestConvert_Invalid(t *testing.T) {
	tests := []struct {
		from, to string
	}{
		{"m", "kg"},      // cross-category
		{"s", "lb"},      // cross-category
		{"unknown", "m"}, // unknown unit
		{"kg", "xyz"},    // unknown unit
	}

	for _, tt := range tests {
		t.Run(tt.from+"→"+tt.to, func(t *testing.T) {
			got, err := Convert(1, tt.from, tt.to)
			assert.Error(t, err)
			assert.Equal(t, 0.0, got)
		})
	}
}

func TestConvert_Identity(t *testing.T) {
	// Conversion to same unit should return same value
	units := []string{"m", "cm", "kg", "g", "s", "ms"}
	for _, u := range units {
		t.Run(u+"→"+u, func(t *testing.T) {
			got, err := Convert(123.456, u, u)
			assert.NoError(t, err)
			assert.Equal(t, 123.456, got)
		})
	}
}

func TestConvert_ZeroValue(t *testing.T) {
	got, err := Convert(0, "m", "cm")
	assert.NoError(t, err)
	assert.Equal(t, 0.0, got)
}

func TestConvert_LargeNumbers(t *testing.T) {
	got, err := Convert(1e6, "km", "m")
	assert.NoError(t, err)
	assert.Equal(t, 1e9, got)
}
