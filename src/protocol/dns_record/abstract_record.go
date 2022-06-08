package dns_record

import (
	"fmt"
	"github.com/wiktor-mazur/dns-go/src/buffer"
	"github.com/wiktor-mazur/dns-go/src/common"
	"strings"
)

type AbstractDnsRecord struct {
	Name       string
	QueryType  common.QueryType
	Class      common.Class
	TTL        uint32
	DataLength uint16
	data       []byte
}

func NewAbstractRecord() AbstractDnsRecord {
	result := AbstractDnsRecord{}
	result.Class = common.IN

	return result
}

func (v *AbstractDnsRecord) GetName() string {
	return v.Name
}

func (v *AbstractDnsRecord) GetType() common.QueryType {
	return v.QueryType
}

func (v *AbstractDnsRecord) ReadPreamble(buf *buffer.BytePacketBuffer) error {
	name, err := buf.ReadLabel()
	if err != nil {
		return err
	}

	queryType, err := buf.ReadUint16()
	if err != nil {
		return err
	}

	class, err := buf.ReadUint16()
	if err != nil {
		return err
	}

	ttl, err := buf.ReadUint32()
	if err != nil {
		return err
	}

	dataLength, err := buf.ReadUint16()
	if err != nil {
		return err
	}

	v.Name = name
	v.QueryType = common.QueryType(queryType)
	v.Class = common.Class(class)
	v.TTL = ttl
	v.DataLength = dataLength

	return nil
}

func (v *AbstractDnsRecord) ReadData(buf *buffer.BytePacketBuffer) error {
	dataRead := uint16(0)

	for dataRead < v.DataLength {
		curByte, err := buf.ReadByte()
		if err != nil {
			return err
		}

		v.data = append(v.data, curByte)
		dataRead += 1
	}

	return nil
}

func (v *AbstractDnsRecord) WritePreamble(buf *buffer.BytePacketBuffer) error {
	err := buf.WriteLabel(v.Name)
	if err != nil {
		return err
	}

	err = buf.WriteUint16(uint16(v.QueryType))
	if err != nil {
		return err
	}

	err = buf.WriteUint16(uint16(v.Class))
	if err != nil {
		return err
	}

	err = buf.WriteUint32(v.TTL)
	if err != nil {
		return err
	}

	return nil
}

func (v *AbstractDnsRecord) WriteData(buf *buffer.BytePacketBuffer) error {
	err := buf.WriteUint16(uint16(len(v.data)))
	if err != nil {
		return err
	}

	for i := uint16(0); i < v.DataLength; i++ {
		err = buf.WriteByte(v.data[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *AbstractDnsRecord) String() string {
	r := new(strings.Builder)

	fmt.Fprintf(r, "Name: %s\n", v.Name)
	fmt.Fprintf(r, "Type: %d (%s)\n", v.QueryType, v.QueryType.String())
	fmt.Fprintf(r, "Class: %d (%s)\n", v.Class, v.Class.String())
	fmt.Fprintf(r, "TTL: %d\n", v.TTL)
	fmt.Fprintf(r, "Data length: %d byte(s)\n", v.DataLength)

	return r.String()
}

func (v *AbstractDnsRecord) CompactString() string {
	return fmt.Sprintf("DnsRecord { Name: %s, Type: %d, Class: %s, TTL: %d }", v.Name, v.QueryType, v.Class.String(), v.TTL)
}
