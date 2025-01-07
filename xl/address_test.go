package xl

import (
	"testing"
)

func TestParseAddress(t *testing.T) {
	badAddresses := []string{
		"Sheet1!",
		"Sheet1!$",
		"'Sheet1!$B",
		"'Sheet1'!$B",
		"'Sheet1'!B",
	}
	for _, strAddress := range badAddresses {
		address, convErr := ParseAddress(strAddress)
		if convErr == nil {
			t.Errorf("ParseAddress(%s) succeeded with '%s' but should error", strAddress, address)
		}
	}
}
