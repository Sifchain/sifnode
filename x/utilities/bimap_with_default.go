package utilities

import "github.com/vishalkuo/bimap"

func GetWithDefault(bimap *bimap.BiMap, k interface{}, defaultValue interface{}) interface{} {
	result, found := bimap.Get(k)
	if found {
		return result
	} else {
		return defaultValue
	}
}

func GetInverseWithDefault(bimap *bimap.BiMap, k interface{}, defaultValue interface{}) interface{} {
	result, found := bimap.GetInverse(k)
	if found {
		return result
	} else {
		return defaultValue
	}
}
