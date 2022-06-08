package protocol

import (
	"fmt"
	"github.com/wiktor-mazur/dns-go/src/buffer"
	"github.com/wiktor-mazur/dns-go/src/common"
	"github.com/wiktor-mazur/dns-go/src/protocol/dns_record"
	"math/rand"
	"strings"
)

type DnsPacket struct {
	Header      DnsHeader
	Questions   []DnsQuestion
	Answers     []DnsRecord
	Authorities []DnsRecord
	Resources   []DnsRecord
}

func NewDnsPacket() *DnsPacket {
	result := &DnsPacket{}

	result.Header.ID = uint16(rand.Uint32())

	return result
}

func DnsPacketFromRawBuffer(rawBuf []byte) (*DnsPacket, error) {
	buf := buffer.BytePacketBufferFromRawBuffer(rawBuf)

	packet := &DnsPacket{}

	err := packet.Header.Read(buf)
	if err != nil {
		return packet, err
	}

	for i := uint16(0); i < packet.Header.QuestionsCount; i++ {
		question := NewDnsQuestion()

		err = question.Read(buf)
		if err != nil {
			return packet, err
		}

		packet.Questions = append(packet.Questions, *question)
	}

	for i := uint16(0); i < packet.Header.AnswersCount; i++ {
		record, err := ReadDnsRecord(buf)
		if err != nil {
			return packet, err
		}

		packet.Answers = append(packet.Answers, record)
	}

	for i := uint16(0); i < packet.Header.AuthoritiesCount; i++ {
		record, err := ReadDnsRecord(buf)
		if err != nil {
			return packet, err
		}

		packet.Authorities = append(packet.Authorities, record)
	}

	for i := uint16(0); i < packet.Header.ResourcesCount; i++ {
		record, err := ReadDnsRecord(buf)
		if err != nil {
			return packet, err
		}

		packet.Resources = append(packet.Resources, record)
	}

	return packet, nil
}

func (v *DnsPacket) ToBuffer() (*buffer.BytePacketBuffer, error) {
	buf := buffer.NewBytePacketBuffer()

	err := v.Write(buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (v *DnsPacket) ToRawBuffer() ([]byte, error) {
	buf := buffer.NewBytePacketBuffer()

	err := v.Write(buf)
	if err != nil {
		return nil, err
	}

	return buf.GetBytes(), nil
}

func (v *DnsPacket) Write(buf *buffer.BytePacketBuffer) error {
	v.Header.QuestionsCount = uint16(len(v.Questions))
	v.Header.AnswersCount = uint16(len(v.Answers))
	v.Header.AuthoritiesCount = uint16(len(v.Authorities))
	v.Header.ResourcesCount = uint16(len(v.Resources))

	err := v.Header.Write(buf)
	if err != nil {
		return err
	}

	for _, question := range v.Questions {
		err = question.Write(buf)
		if err != nil {
			return err
		}
	}

	for _, record := range v.Answers {
		err = record.WritePreamble(buf)
		if err != nil {
			return err
		}

		err = record.WriteData(buf)
		if err != nil {
			return err
		}
	}

	for _, record := range v.Authorities {
		err = record.WritePreamble(buf)
		if err != nil {
			return err
		}

		err = record.WriteData(buf)
		if err != nil {
			return err
		}
	}

	for _, record := range v.Resources {
		err = record.WritePreamble(buf)
		if err != nil {
			return err
		}

		err = record.WriteData(buf)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *DnsPacket) AddQuestion(question DnsQuestion) {
	v.Questions = append(v.Questions, question)
	v.Header.QuestionsCount += 1
}

func (v *DnsPacket) AddAnswer(record DnsRecord) {
	v.Answers = append(v.Answers, record)
	v.Header.AnswersCount += 1
}

func (v *DnsPacket) AddAuthority(record DnsRecord) {
	v.Authorities = append(v.Authorities, record)
	v.Header.AuthoritiesCount += 1
}

func (v *DnsPacket) AddResource(record DnsRecord) {
	v.Resources = append(v.Resources, record)
	v.Header.ResourcesCount += 1
}

// GetAuthorityNameServers returns all NS records from authority section with matching qName
func (v *DnsPacket) GetAuthorityNameServers(qName string) []*dns_record.NS {
	nameServers := make([]*dns_record.NS, 0)

	// extract name servers from authority section
	for _, record := range v.Authorities {
		if record.GetType() == common.NS && strings.HasSuffix(qName, record.GetName()) {
			ns, ok := record.(*dns_record.NS)

			if ok {
				nameServers = append(nameServers, ns)
			}
		}
	}

	return nameServers
}

func (v *DnsPacket) GetResolvedNS(qName string) *dns_record.A {
	// try to find corresponding A record in the additional section
	for _, ns := range v.GetAuthorityNameServers(qName) {
		for _, record := range v.Resources {
			if record.GetType() == common.A && record.GetName() == ns.GetHost() {
				aRecord, ok := record.(*dns_record.A)

				if ok {
					return aRecord
				}
			}
		}
	}

	return nil
}

func (v *DnsPacket) GetUnresolvedNS(qName string) *dns_record.NS {
	ns := v.GetAuthorityNameServers(qName)

	if len(ns) <= 0 {
		return nil
	}

	return ns[0]
}

func (v *DnsPacket) GetFirstARecord() *dns_record.A {
	for _, record := range v.Answers {
		if record.GetType() == common.A {
			return record.(*dns_record.A)
		}
	}

	return nil
}

func (v *DnsPacket) String() string {
	result := "################################################################\n"

	result += "Header:\n---------------------\n" + v.Header.String() + "\n\n"

	result += "Questions:\n---------------------\n"
	for idx, v := range v.Questions {
		result += fmt.Sprintf("%d.\n%s\n", idx+1, v.String())
	}

	result += "\nRecords:\n---------------------\n"
	for idx, v := range v.Answers {
		result += fmt.Sprintf("%d.\n%s\n", idx+1, v.String())
	}

	result += "\nAuthority section:\n---------------------\n"
	for idx, v := range v.Authorities {
		result += fmt.Sprintf("%d.\n%s\n", idx+1, v.String())
	}

	result += "\nAdditional section:\n---------------------\n"
	for idx, v := range v.Resources {
		result += fmt.Sprintf("%d.\n%s\n", idx+1, v.String())
	}

	result += "\n################################################################\n"

	return result
}
