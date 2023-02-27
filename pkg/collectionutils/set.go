package collectionutils

func ToSet(array []string) []string {
	auxMap := make(map[string]bool)

	for _, el := range array {
		auxMap[el] = true
	}

	set := make([]string, 0, len(auxMap))
	for key := range auxMap {
		set = append(set, key)
	}

	return set
}
