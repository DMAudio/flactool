package types

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type Severity uint8

const (
	RsDebug   Severity = 1 << (8 - 1 - iota) //128
	RsInfo                                   //64
	RsNotify                                 //32
	RsWarning                                //16
	RsError                                  //8
	RsPanic                                  //4
	RsUnknown                                //2
)

type Throwable interface {
	GetMessage(bool) string
}

var ThrowableTraceLevel *uint

func InitThrowableExt() {
	ThrowableTraceLevel = flag.Uint("trace", 48, "over which debug level shall stack tracing be enabled")
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
	case RsUnknown:
		fallthrough
	default:
		return "U"
	}
}

func Throw(t Throwable, severity Severity) {
	message := WrappedMessage(
		fmt.Sprintf("[%d][%s]", time.Now().Unix(), severity2Str(severity)),
		t.GetMessage(uint(severity)&*ThrowableTraceLevel > 0),
		false,
	)
	if severity != RsPanic {
		if severity&(RsWarning|RsError|RsPanic|RsUnknown) > 0 {
			_, _ = fmt.Fprintln(os.Stderr, message)
		}else {
			_, _ = fmt.Fprintln(os.Stdout, message)
		}
	} else {
		log.Panicln(message)
	}

}
