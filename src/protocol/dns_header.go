package protocol

import (
	"fmt"
	"github.com/wiktor-mazur/dns-go/src/buffer"
	"github.com/wiktor-mazur/dns-go/src/common"
	"github.com/wiktor-mazur/dns-go/src/utils"
	"strconv"
	"strings"
)

type DnsHeader struct {
	ID                  uint16
	IsResponse          bool
	OPCODE              common.OPCODE
	AuthoritativeAnswer bool
	TruncatedMessage    bool
	RecursionDesired    bool
	RecursionAvailable  bool
	DNSSECAvailable     bool
	AuthedData          bool
	CheckingDisabled    bool
	ResultCode          common.ResultCode
	QuestionsCount      uint16
	AnswersCount        uint16
	AuthoritiesCount    uint16
	ResourcesCount      uint16
}

func (v *DnsHeader) Read(buf *buffer.BytePacketBuffer) error {
	id, err := buf.ReadUint16()
	if err != nil {
		return err
	}

	flagsPartOne, err := buf.ReadByte()
	if err != nil {
		return err
	}

	flagsPartTwo, err := buf.ReadByte()
	if err != nil {
		return err
	}

	questionsCount, err := buf.ReadUint16()
	if err != nil {
		return err
	}

	answersCount, err := buf.ReadUint16()
	if err != nil {
		return err
	}

	authoritiesCount, err := buf.ReadUint16()
	if err != nil {
		return err
	}

	resourcesCount, err := buf.ReadUint16()
	if err != nil {
		return err
	}

	v.ID = id

	v.IsResponse = flagsPartOne>>7 > 0
	v.OPCODE = common.OPCODE(flagsPartOne << 1 >> 4)
	v.AuthoritativeAnswer = flagsPartOne<<5>>7 > 0
	v.TruncatedMessage = flagsPartOne<<6>>7 > 0
	v.RecursionDesired = flagsPartOne<<7>>7 > 0

	v.RecursionAvailable = flagsPartTwo>>7 > 0
	v.DNSSECAvailable = flagsPartTwo<<1>>7 > 0
	v.AuthedData = flagsPartTwo<<2>>7 > 0
	v.CheckingDisabled = flagsPartTwo<<3>>7 > 0
	v.ResultCode = common.ResultCode(flagsPartTwo << 4 >> 4)

	v.QuestionsCount = questionsCount
	v.AnswersCount = answersCount
	v.AuthoritiesCount = authoritiesCount
	v.ResourcesCount = resourcesCount

	return nil
}

func (v *DnsHeader) Write(buf *buffer.BytePacketBuffer) error {
	err := buf.WriteUint16(v.ID)
	if err != nil {
		return err
	}

	flagsPartOne := byte(0)
	flagsPartOne |= utils.BoolToByte(v.RecursionDesired)
	flagsPartOne |= utils.BoolToByte(v.TruncatedMessage) << 1
	flagsPartOne |= utils.BoolToByte(v.AuthoritativeAnswer) << 2
	flagsPartOne |= byte(v.OPCODE) << 3
	flagsPartOne |= utils.BoolToByte(v.IsResponse) << 7
	err = buf.WriteByte(flagsPartOne)
	if err != nil {
		return err
	}

	flagsPartTwo := byte(v.ResultCode)
	flagsPartTwo |= utils.BoolToByte(v.CheckingDisabled) << 4
	flagsPartTwo |= utils.BoolToByte(v.AuthedData) << 5
	flagsPartTwo |= utils.BoolToByte(v.DNSSECAvailable) << 6
	flagsPartTwo |= utils.BoolToByte(v.RecursionAvailable) << 7
	err = buf.WriteByte(flagsPartTwo)
	if err != nil {
		return err
	}

	err = buf.WriteUint16(v.QuestionsCount)
	if err != nil {
		return err
	}

	err = buf.WriteUint16(v.AnswersCount)
	if err != nil {
		return err
	}

	err = buf.WriteUint16(v.AuthoritiesCount)
	if err != nil {
		return err
	}

	err = buf.WriteUint16(v.ResourcesCount)
	if err != nil {
		return err
	}

	return nil
}

func (v *DnsHeader) String() string {
	r := new(strings.Builder)

	hType := "question"

	if v.IsResponse == true {
		hType = "answer"
	}

	fmt.Fprintf(r, "ID: %d\n", v.ID)
	fmt.Fprintf(r, "Type: %s\n", hType)
	fmt.Fprintf(r, "OPCODE: %d (%s)\n", v.OPCODE, v.OPCODE.String())
	fmt.Fprintf(r, "Authoritative: %s\n", strconv.FormatBool(v.AuthoritativeAnswer))
	fmt.Fprintf(r, "Truncated (over 512 bytes): %s\n", strconv.FormatBool(v.TruncatedMessage))
	fmt.Fprintf(r, "Recursion desired: %s\n", strconv.FormatBool(v.RecursionDesired))
	fmt.Fprintf(r, "Recursion available: %s\n", strconv.FormatBool(v.RecursionAvailable))
	fmt.Fprintf(r, "DNSSEC Available: %s\n", strconv.FormatBool(v.DNSSECAvailable))
	fmt.Fprintf(r, "Is data authenticated (DNSSEC): %s\n", strconv.FormatBool(v.AuthedData))
	fmt.Fprintf(r, "Is checking disabled (DNSSEC): %s\n", strconv.FormatBool(v.CheckingDisabled))
	fmt.Fprintf(r, "Result code: %d (%s)\n", v.ResultCode, v.ResultCode.String())
	fmt.Fprintf(r, "Questions count: %d\n", v.QuestionsCount)
	fmt.Fprintf(r, "Answers count: %d\n", v.AnswersCount)
	fmt.Fprintf(r, "Authorities count: %d\n", v.AuthoritiesCount)
	fmt.Fprintf(r, "Resources count: %d\n", v.ResourcesCount)

	return r.String()
}
