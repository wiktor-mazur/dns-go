package dns_record

import (
	"fmt"
	"github.com/wiktor-mazur/dns-go/src/buffer"
	"github.com/wiktor-mazur/dns-go/src/utils"
	"strings"
)

type A struct {
	AbstractDnsRecord
	ip utils.IPv4
}

func (v *A) GetIP() utils.IPv4 {
	return v.ip
}

func (v *A) ReadData(buf *buffer.BytePacketBuffer) error {
	if v.DataLength != uint16(4) {
		return fmt.Errorf("invalid data length in A record")
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

	ip, err := utils.NewIPv4(data)
	if err != nil {
		return err
	}

	v.ip = *ip

	return nil
}

func (v *A) WriteData(buf *buffer.BytePacketBuffer) error {
	err := buf.WriteUint16(uint16(4))
	if err != nil {
		return err
	}

	for i := uint16(0); i < 4; i++ {
		err = buf.WriteByte(v.ip.Octets[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *A) String() string {
	r := new(strings.Builder)

	fmt.Fprintf(r, v.AbstractDnsRecord.String())
	fmt.Fprintf(r, "IPv4: %s", v.ip.String())

	return r.String()
}

func (v *A) CompactString() string {
	return fmt.Sprintf("A { Domain: %s, IP: %s, TTL: %d }", v.Name, v.ip.String(), v.TTL)
}
