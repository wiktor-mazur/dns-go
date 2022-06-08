package common

type QueryType uint16

const (
	A     QueryType = 1
	NS    QueryType = 2
	CNAME QueryType = 5
	SOA   QueryType = 6
	MX    QueryType = 15
	AAAA  QueryType = 28
)

func (v *QueryType) String() string {
	switch *v {
	case A:
		return "A"
	case NS:
		return "NS"
	case CNAME:
		return "CNAME"
	case SOA:
		return "SOA"
	case MX:
		return "MX"
	case AAAA:
		return "AAAA"
	default:
		return "UNKNOWN"
	}
}
