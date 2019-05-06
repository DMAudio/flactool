package types

import (
	"fmt"
	"runtime/debug"
	"strings"
)

type Exception struct {
	mask  Mask
	param *SSMap
	stack []byte
	cause interface{}
}

func (t *Exception) GetMessage(withStack bool) string {
	if t == nil {
		return ""
	}
	message := NewBuffer()

	message.WriteStringsIE(
		"#", t.mask.GetMask(), "# ",
		t.mask.GetMessage(t.param.Dump()), "\n",
	)

	if t.cause != nil {
		switch t.cause.(type) {
		case map[string]*Exception:
			for title, throwableObj := range t.cause.(map[string]*Exception) {
				message.WriteStringsIE(WrappedMessage(
					fmt.Sprintf("^ %s: ", title),
					throwableObj.GetMessage(withStack),
					true,
				), "\n")
			}
		case []*Exception:
			for _, throwableObj := range t.cause.([]*Exception) {
				message.WriteStringsIE(WrappedMessage(
					"^ ",
					throwableObj.GetMessage(withStack),
					true,
				), "\n")
			}
		case *Exception:
			message.WriteStringsIE(WrappedMessage(
				"^ ",
				t.cause.(*Exception).GetMessage(withStack),
				true,
			), "\n")
		default:
			message.WriteStringsIE(WrappedMessage(
				"^ ",
				fmt.Sprintf("%v", t.cause),
				false,
			), "\n")
		}
	}

	if withStack {
		message.WriteStringsIE(WrappedMessage(
			"> ",
			string(t.stack),
			true,
		), "\n")
	}
	return strings.TrimSpace(message.String())
}
func (t *Exception) SetCause(cause interface{}) *Exception {
	t.cause = cause
	return t
}
func (t *Exception) TraceStack() *Exception {
	t.stack = debug.Stack()
	return t
}

func (t *Exception) SetParam(key, value string) *Exception {
	t.param.Set(key, value)
	return t
}

func NewException(mask Mask, param map[string]string, cause interface{}) *Exception {
	t := &Exception{
		mask:  mask,
		param: (&SSMap{}).From(param),
		stack: debug.Stack(),
		cause: cause,
	}

	return t
}

