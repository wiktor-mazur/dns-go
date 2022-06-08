package dns_record

import (
	"fmt"
	"github.com/wiktor-mazur/dns-go/src/buffer"
	"strings"
)

type SOA struct {
	AbstractDnsRecord
	mName   string
	rName   string
	serial  uint32
	refresh uint32
	retry   uint32
	expire  uint32
	minimum uint32
}

func (v *SOA) ReadData(buf *buffer.BytePacketBuffer) error {
	mName, err := buf.ReadLabel()
	if err != nil {
		return err
	}

	rName, err := buf.ReadLabel()
	if err != nil {
		return err
	}

	serial, err := buf.ReadUint32()
	if err != nil {
		return err
	}

	refresh, err := buf.ReadUint32()
	if err != nil {
		return err
	}

	retry, err := buf.ReadUint32()
	if err != nil {
		return err
	}

	expire, err := buf.ReadUint32()
	if err != nil {
		return err
	}

	minimum, err := buf.ReadUint32()
	if err != nil {
		return err
	}

	v.mName = mName
	v.rName = rName
	v.serial = serial
	v.refresh = refresh
	v.retry = retry
	v.expire = expire
	v.minimum = minimum

	return nil
}

func (v *SOA) WriteData(buf *buffer.BytePacketBuffer) error {
	err := buf.PrependDataLength(func() error {
		err := buf.WriteLabel(v.mName)
		if err != nil {
			return err
		}

		err = buf.WriteLabel(v.rName)
		if err != nil {
			return err
		}

		err = buf.WriteUint32(v.serial)
		if err != nil {
			return err
		}

		err = buf.WriteUint32(v.refresh)
		if err != nil {
			return err
		}

		err = buf.WriteUint32(v.retry)
		if err != nil {
			return err
		}

		err = buf.WriteUint32(v.expire)
		if err != nil {
			return err
		}

		err = buf.WriteUint32(v.minimum)
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

func (v *SOA) String() string {
	r := new(strings.Builder)

	fmt.Fprintf(r, v.AbstractDnsRecord.String())
	fmt.Fprintf(r, "m_name: %s\n", v.mName)
	fmt.Fprintf(r, "r_name: %s\n", v.rName)
	fmt.Fprintf(r, "Serial: %d\n", v.serial)
	fmt.Fprintf(r, "Refresh: %d\n", v.refresh)
	fmt.Fprintf(r, "Retry: %d\n", v.retry)
	fmt.Fprintf(r, "Expire: %d\n", v.expire)
	fmt.Fprintf(r, "Minimum: %d", v.minimum)

	return r.String()
}

func (v *SOA) CompactString() string {
	return fmt.Sprintf("SOA { Domain: %s, MName: %s, RName: %s, Serial: %d, Refresh: %d, Retry: %d, Expire: %d, Minimum: %d, TTL: %d }", v.Name, v.mName, v.rName, v.serial, v.refresh, v.retry, v.expire, v.minimum, v.TTL)
}
