package dns_record

import (
	"fmt"
	"github.com/wiktor-mazur/dns-go/src/buffer"
	"strings"
)

type MX struct {
	AbstractDnsRecord
	priority uint16
	host     string
}

func (v *MX) ReadData(buf *buffer.BytePacketBuffer) error {
	priority, err := buf.ReadUint16()
	if err != nil {
		return err
	}

	host, err := buf.ReadLabel()
	if err != nil {
		return err
	}

	v.priority = priority
	v.host = host

	return nil
}

func (v *MX) WriteData(buf *buffer.BytePacketBuffer) error {
	err := buf.PrependDataLength(func() error {
		err := buf.WriteUint16(v.priority)
		if err != nil {
			return err
		}

		err = buf.WriteLabel(v.host)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (v *MX) String() string {
	r := new(strings.Builder)

	fmt.Fprintf(r, v.AbstractDnsRecord.String())
	fmt.Fprintf(r, "Priority: %d\n", v.priority)
	fmt.Fprintf(r, "Host: %s", v.host)

	return r.String()
}

func (v *MX) CompactString() string {
	return fmt.Sprintf("MX { Domain: %s, Priority: %d, Host: %s, TTL: %d }", v.Name, v.priority, v.host, v.TTL)
}
