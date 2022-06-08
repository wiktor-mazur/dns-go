package protocol

import (
	"github.com/wiktor-mazur/dns-go/src/buffer"
	"github.com/wiktor-mazur/dns-go/src/common"
	"github.com/wiktor-mazur/dns-go/src/protocol/dns_record"
)

type DnsRecord interface {
	GetName() string
	GetType() common.QueryType
	ReadPreamble(buf *buffer.BytePacketBuffer) error
	ReadData(buf *buffer.BytePacketBuffer) error
	WritePreamble(buf *buffer.BytePacketBuffer) error
	WriteData(buf *buffer.BytePacketBuffer) error
	String() string
	CompactString() string
}

func ReadDnsRecord(buf *buffer.BytePacketBuffer) (DnsRecord, error) {
	var record DnsRecord

	abstract := dns_record.NewAbstractRecord()
	err := abstract.ReadPreamble(buf)
	if err != nil {
		return nil, err
	}

	switch abstract.QueryType {
	case common.A:
		record = &dns_record.A{AbstractDnsRecord: abstract}
		break
	case common.NS:
		record = &dns_record.NS{AbstractDnsRecord: abstract}
		break
	case common.CNAME:
		record = &dns_record.CNAME{AbstractDnsRecord: abstract}
		break
	case common.SOA:
		record = &dns_record.SOA{AbstractDnsRecord: abstract}
		break
	case common.MX:
		record = &dns_record.MX{AbstractDnsRecord: abstract}
		break
	case common.AAAA:
		record = &dns_record.AAAA{AbstractDnsRecord: abstract}
		break
	default:
		record = &abstract
	}

	err = record.ReadData(buf)
	if err != nil {
		return nil, err
	}

	return record, nil
}
