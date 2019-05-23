package types

import (
	"fmt"
	"strings"
)

type Mask struct {
	mask     string
	template string
}

func NewMask(mask string, template string) Mask {
	return Mask{mask: mask, template: template}
}

func (m *Mask) GetMask() string {
	return m.mask
}

func (m *Mask) GetMessage(args map[string]string) string {
	var message = m.template

	if message == "" {
		return fmt.Sprintf("%s|%v", m.mask, args)
	}

	for arg, content := range args {
		message = strings.Replace(message, "{{"+arg+"}}", content, -1)
	}
	return message
}
