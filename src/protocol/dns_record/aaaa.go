package dns_record

import (
	"fmt"
	"github.com/wiktor-mazur/dns-go/src/buffer"
	"github.com/wiktor-mazur/dns-go/src/utils"
	"strings"
)

type AAAA struct {
	AbstractDnsRecord
	ip utils.IPv6
}

func (v *AAAA) ReadData(buf *buffer.BytePacketBuffer) error {
	if v.DataLength != uint16(16) {
		return fmt.Errorf("invalid data length in AAAA record")
	}

	data := make([]byte, 0)
	dataRead := uint16(0)

	for dataRead < v.DataLength {
		curByte, err := buf.ReadByte()
		if err != nil {
			return err
		}

		data = append(data, curByte)
		dataRead += 1
	}

	ip, err := utils.NewIPv6(data)
	if err != nil {
		return err
	}

	v.ip = *ip

	return nil
}

func (v *AAAA) WriteData(buf *buffer.BytePacketBuffer) error {
	err := buf.WriteUint16(uint16(16))
	if err != nil {
		return err
	}

	for i := uint16(0); i < 16; i++ {
		err = buf.WriteByte(v.ip.Data[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *AAAA) String() string {
	r := new(strings.Builder)

	fmt.Fprintf(r, v.AbstractDnsRecord.String())
	fmt.Fprintf(r, "IPv6: %s", v.ip.String())

	return r.String()
}

func (v *AAAA) CompactString() string {
	return fmt.Sprintf("A { Domain: %s, IP: %s, TTL: %d }", v.Name, v.ip.String(), v.TTL)
}
