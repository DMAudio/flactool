package types

import (
	"flag"
	"fmt"
	"log"
	"time"
)

type Severity uint

const (
	RsDebug   Severity = 1
	RsInfo    Severity = 2
	RsNotify  Severity = 4
	RsWarning Severity = 8
	RsError   Severity = 16
	RsPanic   Severity = 32
)

type Throwable interface {
	GetMessage(bool) string
}

var ThrowableTraceLevel *uint

func InitThrowableExt() {
	ThrowableTraceLevel = flag.Uint("trace", 48, "")
}

func severity2Str(s Severity) string {
	switch s {
	case RsDebug:
		return "D"
	case RsInfo:
		return "I"
	case RsNotify:
		return "N"
	case RsWarning:
		return "W"
	case RsError:
		return "E"
	case RsPanic:
		return "P"
	}
	return "(unknown)"
}

func Throw(t Throwable, severity Severity) {
	message := WrappedMessage(
		fmt.Sprintf("[%d][%s]", time.Now().Unix(), severity2Str(severity)),
		t.GetMessage(uint(severity)&*ThrowableTraceLevel > 0),
		false,
	)
	if severity != RsPanic {
		fmt.Println(message)
	} else {
		log.Panicln(message)
	}
}
