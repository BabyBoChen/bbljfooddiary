package utils

func GetValueFromFormData(form map[string][]string, key string) string {
	var val string
	vals, ok := form[key]
	if ok {
		if len(vals) >= 1 {
			val = vals[0]
		}
	}
	return val
}
