package buffer

import (
	"fmt"
	"strings"
)

const DNSBufferSize = 512

var PosTooLargeErr = fmt.Errorf("pos must not be larger than %d", DNSBufferSize)

func bufferOverflowErr(pos uint) error {
	return fmt.Errorf("buffer overflow (tried to access %d byte of %d total)", pos, DNSBufferSize)
}

type BytePacketBuffer struct {
	buffer []byte
	pos    uint // @TODO uint
}

func NewBytePacketBuffer() *BytePacketBuffer {
	return &BytePacketBuffer{
		buffer: make([]byte, DNSBufferSize),
		pos:    0,
	}
}

func BytePacketBufferFromRawBuffer(buf []byte) *BytePacketBuffer {
	return &BytePacketBuffer{
		buffer: buf,
		pos:    0,
	}
}

func (v *BytePacketBuffer) GetPos() uint {
	return v.pos
}

func (v *BytePacketBuffer) GetBytes() []byte {
	return v.buffer[:v.pos]
}

func (v *BytePacketBuffer) Seek(pos uint) error {
	if pos > DNSBufferSize-1 {
		return PosTooLargeErr
	}

	v.pos = pos

	return nil
}

func (v *BytePacketBuffer) ReadByteAt(pos uint) (byte, error) {
	if pos > DNSBufferSize-1 {
		return 0, PosTooLargeErr
	}

	return v.buffer[pos], nil
}

func (v *BytePacketBuffer) ReadAtRange(pos uint, len uint) ([]byte, error) {
	if pos+len > DNSBufferSize {
		return nil, bufferOverflowErr(pos + len)
	}

	result := make([]byte, len)

	for i := uint(0); i < len; i++ {
		currentByte, err := v.ReadByteAt(pos + i)
		if err != nil {
			return nil, err
		}

		result[i] = currentByte
	}

	return result, nil
}

func (v *BytePacketBuffer) ReadByte() (byte, error) {
	result := v.buffer[v.pos]

	if v.pos+1 < DNSBufferSize {
		v.pos += 1
	}

	return result, nil
}

func (v *BytePacketBuffer) ReadUint16() (uint16, error) {
	firstByte, err := v.ReadByte()
	if err != nil {
		return 0, err
	}

	secondByte, err := v.ReadByte()
	if err != nil {
		return 0, err
	}

	return uint16(firstByte)<<8 | uint16(secondByte), nil
}

func (v *BytePacketBuffer) ReadUint32() (uint32, error) {
	firstUint16, err := v.ReadUint16()
	if err != nil {
		return 0, err
	}

	secondUint16, err := v.ReadUint16()
	if err != nil {
		return 0, err
	}

	return uint32(firstUint16)<<16 | uint32(secondUint16), nil
}

// @TODO: refactor/comment
func (v *BytePacketBuffer) ReadLabel() (string, error) {
	result := ""
	localPos := v.pos

	delim := ""

	jumpsCount := 0
	maxJumps := 5

	for {
		if jumpsCount > maxJumps {
			return "", fmt.Errorf("exceeded the limit of %d jumps while reading label", maxJumps)
		}

		lengthByte, err := v.ReadByteAt(localPos)
		if err != nil {
			return "", err
		}

		if lengthByte&0xC0 == 0xC0 {
			if jumpsCount == 0 {
				err := v.Seek(localPos + 2)
				if err != nil {
					return "", err
				}
			}

			jumpByte, err := v.ReadByteAt(localPos + 1)
			if err != nil {
				return "", err
			}

			jumpOffset := ((uint16(lengthByte) ^ 0xC0) << 8) | uint16(jumpByte)

			localPos = uint(jumpOffset)

			jumpsCount += 1
		} else {
			localPos += 1

			if lengthByte == 0 {
				break
			}

			labelChunk, err := v.ReadAtRange(localPos, uint(lengthByte))
			if err != nil {
				return "", err
			}

			result += delim + string(labelChunk)

			localPos += uint(lengthByte)

			delim = "."
		}
	}

	if jumpsCount == 0 {
		err := v.Seek(localPos)
		if err != nil {
			return "", err
		}
	}

	return result, nil
}

func (v *BytePacketBuffer) SetByte(pos uint, byte byte) error {
	if pos > DNSBufferSize-1 {
		return PosTooLargeErr
	}

	v.buffer[pos] = byte

	return nil
}

func (v *BytePacketBuffer) SetUint16(pos uint, val uint16) error {
	err := v.SetByte(pos, byte(val>>8))
	if err != nil {
		return err
	}

	err = v.SetByte(pos+1, byte(val<<8>>8))
	if err != nil {
		return err
	}

	return nil
}

func (v *BytePacketBuffer) WriteByte(byte byte) error {
	if v.pos >= DNSBufferSize-1 {
		return bufferOverflowErr(v.pos + 1)
	}

	v.buffer[v.pos] = byte
	v.pos++

	return nil
}

func (v *BytePacketBuffer) WriteUint16(uint16 uint16) error {
	err := v.WriteByte(byte(uint16 >> 8))
	if err != nil {
		return err
	}

	err = v.WriteByte(byte(uint16 << 8 >> 8))
	if err != nil {
		return err
	}

	return nil
}

func (v *BytePacketBuffer) WriteUint32(uint32 uint32) error {
	err := v.WriteUint16(uint16(uint32 >> 16))
	if err != nil {
		return err
	}

	err = v.WriteUint16(uint16(uint32 << 16 >> 16))
	if err != nil {
		return err
	}

	return nil
}

func (v *BytePacketBuffer) WriteLabel(label string) error {
	for _, label := range strings.Split(label, ".") {
		if len(label) > 0x3f {
			return fmt.Errorf("given label is too long")
		}

		err := v.WriteByte(byte(len(label)))
		if err != nil {
			return err
		}

		for _, b := range []byte(label) {
			err := v.WriteByte(b)
			if err != nil {
				return err
			}
		}
	}

	err := v.WriteByte(0)
	if err != nil {
		return err
	}

	return nil
}

func (v *BytePacketBuffer) PrependDataLength(writeData func() error) error {
	dataLengthPos := v.pos
	err := v.WriteUint16(0)
	if err != nil {
		return err
	}

	err = writeData()
	if err != nil {
		return err
	}

	size := v.pos - (dataLengthPos + 2)
	err = v.SetUint16(dataLengthPos, uint16(size))
	if err != nil {
		return err
	}

	return nil
}
