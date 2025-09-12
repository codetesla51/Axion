package constants

import (
    "encoding/json"
    "fmt"
    "os"
)

var Table map[string]float64

func Load(file string) error {
    f, err := os.ReadFile(file)
    if err != nil {
        return fmt.Errorf("failed to read constants file: %w", err)
    }
    Table = make(map[string]float64)
    if err := json.Unmarshal(f, &Table); err != nil {
        return fmt.Errorf("failed to parse constants: %w", err)
    }
    return nil
}

func Get(name string) (float64, bool) {
    val, ok := Table[name]
    return val, ok
}