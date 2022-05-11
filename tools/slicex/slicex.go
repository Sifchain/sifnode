package slicex

import "strings"

// StringsContain returns true if the given value contains at least one value from the given slice
func StringsContain(val string, sl []string) bool {
	for _, sv := range sl {
		if strings.Contains(val, sv) {
			return true
		}
	}
	return false
}
