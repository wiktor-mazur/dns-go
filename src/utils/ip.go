package utils

import "fmt"

type IPv4 struct {
	Octets []byte
}

func NewIPv4(data []byte) (*IPv4, error) {
	if len(data) != 4 {
		return nil, fmt.Errorf("invalid IPv4")
	}

	return &IPv4{Octets: data}, nil
}

func (v *IPv4) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", v.Octets[0], v.Octets[1], v.Octets[2], v.Octets[3])
}

type IPv6 struct {
	Data []byte
}

func NewIPv6(data []byte) (*IPv6, error) {
	if len(data) != 16 {
		return nil, fmt.Errorf("invalid IPv6")
	}

	return &IPv6{Data: data}, nil
}

func (v *IPv6) String() string {
	return fmt.Sprintf(
		"%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x",
		v.Data[0], v.Data[1], v.Data[2], v.Data[3], v.Data[4], v.Data[5], v.Data[6], v.Data[7], v.Data[8], v.Data[9], v.Data[10], v.Data[11], v.Data[12], v.Data[13], v.Data[14], v.Data[15],
	)
}
