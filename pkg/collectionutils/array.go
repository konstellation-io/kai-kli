package collectionutils

func ArrayContains(array []string, value string) bool {
	for _, el := range array {
		if el == value {
			return true
		}
	}

	return false
}
