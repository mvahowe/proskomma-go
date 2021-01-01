package succinct

import (
	b64 "encoding/base64"
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
	newBa := make([]uint8, ba.usedBytes)
	copy(newBa, ba.bytes)
	ba.bytes = newBa
	return nil
}

func (ba *ByteArray) PushNByte(v uint32) error {
	if v < 128 {
		err := ba.PushByte(uint8(v + 128))
		if err != nil {
			return err
		} else {
			return nil
		}
	} else {
		mod := v % 128
		err := ba.PushByte(uint8(mod))
		if err != nil {
			return err
		}
		return ba.PushNByte(v >> 7)
	}
}

func (ba *ByteArray) NByte(n int) (uint32, error) {
	if n >= ba.usedBytes {
		return 0, fmt.Errorf(
			"attempt to read byte %d for NByte(length %d)",
			n,
			ba.usedBytes,
		)
	}
	v, err := ba.Byte(n)
	if err != nil {
		return 0, err
	}
	if v > 127 {
		return uint32(v - 128), nil
	} else {
		v2, err := ba.NByte(n + 1)
		if err != nil {
			return 0, err
		} else {
			return uint32(v) + (v2 * 128), nil
		}
	}
}

func (ba *ByteArray) PushCountedString(s string) error {
	sA := []byte(s)
	err := ba.PushByte(uint8(len(sA)))
	if err != nil {
		return err
	}
	err = ba.PushBytes(sA)
	if err != nil {
		return err
	}
	return nil
}

func (ba *ByteArray) CountedString(n int) (string, error) {
	sLength, err := ba.Byte(n)
	if err != nil {
		return "", err
	}
	sA, err := ba.Bytes(n + 1, int(sLength))
	if err != nil {
		return "", err
	}
	return string(sA), nil
}

func (ba *ByteArray) Clear() {
	for n := 0; n < ba.usedBytes; n++ {
		ba.bytes[n] = 0
	}
	ba.usedBytes = 0
}

func (ba *ByteArray) NByteLength(v int) int {
	ret := 1
	for v > 127 {
		v = v >> 7
		ret++
	}
	return ret
}

func (ba *ByteArray) base64() string {
	return b64.StdEncoding.EncodeToString(ba.bytes)
}

func (ba *ByteArray) fromBase64(s string) {
	bytes, _ := b64.StdEncoding.DecodeString(s)
	ba.bytes = bytes
	ba.usedBytes = len(bytes)
}

func SayHello(str string) {
	fmt.Println(str)
}
