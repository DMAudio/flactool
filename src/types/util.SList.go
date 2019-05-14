package types

func SListDeleteByIndex(input []string, elIndex int) []string {
	return append(input[:elIndex], input[elIndex+1:]...)
}

func SListDeleteByElement(input []string, el string) []string {
	target := make([]string, 0)
	for _, item := range input {
		if item != el {
			target = append(target, item)
		}
	}
	return target
}

func SListInsertBefore(input []string, pos int, els ...string) []string {
	if pos > len(input) {
		pos = len(input)
	}
	rear := append([]string{}, input[pos:]...)
	result := append(input[:pos], els...)
	result = append(result, rear...)
	return result
}

func SListInsertAfter(input []string, pos int, els ...string) []string {
	if pos > len(input)-1 {
		pos = len(input) - 1
	}
	rear := append([]string{}, input[pos+1:]...)
	result := append(input[:pos+1], els...)
	result = append(result, rear...)
	return result
}