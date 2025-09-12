  package settings

import "fmt"

var Precision = 6 

func Set(p int) error {
	if p < 0 || p > 20 {
		return fmt.Errorf("precision must be between 0 and 20")
	}
	Precision = p
	return nil
}