package helpers

func MergeStringMapLeft(left map[string]string, right map[string]string) map[string]string {
	if left == nil && right == nil {
		return nil
	}

	newMap := make(map[string]string)

	for key, value := range right {
		newMap[key] = value
	}

	for key, value := range left {
		newMap[key] = value
	}

	return newMap
}

func MergeStringMapRight(left map[string]string, right map[string]string) map[string]string {
	if left == nil && right == nil {
		return nil
	}

	newMap := make(map[string]string)

	for key, value := range left {
		newMap[key] = value
	}

	for key, value := range right {
		newMap[key] = value
	}

	return newMap
}
