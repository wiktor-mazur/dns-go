package dns_record

import (
	"fmt"
	"github.com/wiktor-mazur/dns-go/src/buffer"
	"strings"
)

type NS struct {
	AbstractDnsRecord
	host string
}

func (v *NS) GetHost() string {
	return v.host
}

func (v *NS) ReadData(buf *buffer.BytePacketBuffer) error {
	host, err := buf.ReadLabel()
	if err != nil {
		return err
	}

	v.host = host

	return nil
}

func (v *NS) WriteData(buf *buffer.BytePacketBuffer) error {
	err := buf.PrependDataLength(func() error {
		return buf.WriteLabel(v.host)
	})
	if err != nil {
		return err
	}

	return nil
}

func (v *NS) String() string {
	r := new(strings.Builder)

	fmt.Fprintf(r, v.AbstractDnsRecord.String())
	fmt.Fprintf(r, "Host: %s", v.host)

	return r.String()
}

func (v *NS) CompactString() string {
	return fmt.Sprintf("NS { Domain: %s, Host: %s, TTL: %d }", v.Name, v.host, v.TTL)
}
