package common

type ResultCode uint8

const (
	NOERROR  ResultCode = 0
	FORMERR  ResultCode = 1
	SERVFAIL ResultCode = 2
	NXDOMAIN ResultCode = 3
	NOTIMP   ResultCode = 4
	REFUSED  ResultCode = 5
)

func (v *ResultCode) String() string {
	switch *v {
	case NOERROR:
		return "NOERROR"
	case FORMERR:
		return "FORMERR"
	case SERVFAIL:
		return "SERVFAIL"
	case NXDOMAIN:
		return "NXDOMAIN"
	case NOTIMP:
		return "NOTIMP"
	case REFUSED:
		return "REFUSED"
	default:
		return "UNKNOWN"
	}
}
