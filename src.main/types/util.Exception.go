package types

func Mismatched_Format_Exception(expected string, got string) *Exception {
	return NewException(NewMask(
		"Mismatched_FORMAT_ERROR",
		"格式错误：预期：{{expected}}，实际：{{got}}",
	), map[string]string{
		"expected": expected,
		"got":      got,
	}, nil)
}