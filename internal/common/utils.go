package common

func MapIsNulOrEmpty(target map[string]string) bool {
	return target == nil || len(target) == 0
}

func IsStringInSlice(target string, list []string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}
