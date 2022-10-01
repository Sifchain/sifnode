package bimap_with_default //nolint:stylecheck

import "github.com/vishalkuo/bimap"

func GetWithDefault(bimap *bimap.BiMap, k interface{}, defaultValue interface{}) interface{} {
	result, found := bimap.Get(k)
	if found {
		return result
	}
	return defaultValue
}

func GetInverseWithDefault(bimap *bimap.BiMap, k interface{}, defaultValue interface{}) interface{} {
	result, found := bimap.GetInverse(k)
	if found {
		return result
	}
	return defaultValue
}
