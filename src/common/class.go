package common

type Class uint16

/*
 * All classes except IN are now obsolete and not used
 */

const (
	IN Class = 1 // Internet
)

func (v *Class) String() string {
	switch *v {
	case IN:
		return "IN"
	default:
		return "UNKNOWN"
	}
}
