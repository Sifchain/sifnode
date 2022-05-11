package slicex

// ContainsString returns true if the given "sl" slice contains the given value "val".
func ContainsString(sl []string, val string) bool {
	for _, sv := range sl {
		if sv == val {
			return true
		}
	}
	return false
}
