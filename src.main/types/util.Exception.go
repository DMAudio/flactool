package types

func Exception_Mismatched_Format(expected string, got string) *Exception {
	return NewException(NewMask(
		"Mismatched_FORMAT",
		"格式错误：预期：{{expected}}，实际：{{got}}",
	), map[string]string{
		"expected": expected,
		"got":      got,
	}, nil)
}

func Exception_Mismatched_ArgumentList(expected string, got string) *Exception {
	return NewException(NewMask(
		"Mismatched_ArgumentList",
		"传入参数不符合要求：预期：{{expected}}，实际：{{got}}",
	), map[string]string{
		"expected": expected,
		"got":      got,
	}, nil)
}