package common

type OPCODE uint8

const (
	QUERY                 OPCODE = 0
	INVERSE_QUERY         OPCODE = 1
	SERVER_STATUS_REQUEST OPCODE = 2
)

func (v *OPCODE) String() string {
	switch *v {
	case QUERY:
		return "QUERY"
	case INVERSE_QUERY:
		return "IQUERY"
	case SERVER_STATUS_REQUEST:
		return "STATUS"
	default:
		return "UNKNOWN"
	}
}
