package protocol

import (
	"fmt"
	"github.com/wiktor-mazur/dns-go/src/buffer"
	"github.com/wiktor-mazur/dns-go/src/common"
	"strings"
)

type DnsQuestion struct {
	Name      string
	QueryType common.QueryType
	Class     common.Class
}

func NewDnsQuestion() *DnsQuestion {
	return &DnsQuestion{Class: common.IN}
}

func (v *DnsQuestion) Read(buf *buffer.BytePacketBuffer) error {
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

	v.Name = name
	v.QueryType = common.QueryType(queryType)
	v.Class = common.Class(class)

	return nil
}

func (v *DnsQuestion) Write(buf *buffer.BytePacketBuffer) error {
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

	return nil
}

func (v *DnsQuestion) String() string {
	r := new(strings.Builder)

	fmt.Fprintf(r, "Name: %s\n", v.Name)
	fmt.Fprintf(r, "Type: %d (%s)\n", v.QueryType, v.QueryType.String())
	fmt.Fprintf(r, "Class: %d (%s)\n", v.Class, v.Class.String())

	return r.String()
}

func (v *DnsQuestion) CompactString() string {
	return fmt.Sprintf("DnsQuestion { Name: %s, Type: %s, Class: %s }", v.Name, v.QueryType.String(), v.Class.String())
}
