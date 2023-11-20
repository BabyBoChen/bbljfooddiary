package utils

func MapContainsKey[T comparable](dict map[string]T, key string) bool {
	hasKey := false
	for k := range dict {
		if k == key {
			hasKey = true
			break
		}
	}
	return hasKey
}
