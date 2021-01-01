package succinct

import (
	"fmt"
)

type ByteArray struct {
	bytes []uint8
	usedBytes int
}

func NewByteArray(size uint) ByteArray {
	return ByteArray{
		bytes: make([]uint8, size),
		usedBytes: 0,
	}
}

func (ba *ByteArray) Byte(n int) (uint8, error) {
	if n >= ba.usedBytes {
		return 0, fmt.Errorf(
			"attempt to read byte %d (length %d)",
			n,
			ba.usedBytes,
		)
	}
	return ba.bytes[n], nil
}

func (ba *ByteArray) Bytes(n int, l int) ([]uint8, error) {
	if (n + l) > ba.usedBytes {
		return []uint8{}, fmt.Errorf(
			"attempt to read %d bytes from %d (length %d)",
			l,
			n,
			ba.usedBytes,
			)
	}
	return ba.bytes[n: n + l], nil
}

func (ba *ByteArray) SetByte(n int, v uint8) error {
	if n >= ba.usedBytes {
		return fmt.Errorf(
			"attempt to set byte %d (length %d)",
			n,
			ba.usedBytes,
		)
	}
	ba.bytes[n] = v
	return nil
}

func (ba *ByteArray) SetBytes(n int, vl []uint8) error {
	if (n + len(vl)) > ba.usedBytes {
		return fmt.Errorf(
			"attempt to set %d bytes from %d (length %d)",
			len(vl),
			n,
			ba.usedBytes,
		)
	}
	for i := range vl {
		ba.bytes[n + i] = vl[i]
	}
	return nil
}

func (ba *ByteArray) PushByte(v uint8) error {
	if ba.usedBytes == len(ba.bytes) {
		err := ba.grow()
		if err != nil {
			return err
		}
	}
	ba.bytes[ba.usedBytes] = v
	ba.usedBytes += 1
	return nil
}

func (ba *ByteArray) PushBytes(vl []uint8) error {
	for v := range vl {
		var err = ba.PushByte(vl[v])
		if err != nil {
			return err
		}
	}
	return nil
}

func (ba *ByteArray) grow() error {
	newBa := make([]uint8, len(ba.bytes) * 2)
	copy(newBa, ba.bytes)
	ba.bytes = newBa
	return nil
}

func (ba *ByteArray) Trim() error {
	newBa := make([]uint8, 0, ba.usedBytes)
	copy(newBa, ba.bytes)
	ba.bytes = newBa
	return nil
}

func SayHello(str string) {
	fmt.Println(str)
}
