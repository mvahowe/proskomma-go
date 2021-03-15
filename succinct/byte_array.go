package succinct

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
)

type ByteArray struct {
	bytes []uint8
}

func NewByteArray(size uint) ByteArray {
	return ByteArray{
		bytes: make([]uint8, 0, size),
	}
}

func NewByteArrayFromBase64(s string) (ByteArray, error) {
	bytes, err := b64.StdEncoding.DecodeString(s)
	if err != nil {
		return ByteArray{}, err
	}
	ba := NewByteArray(uint(len(bytes)))
	ba.bytes = bytes
	return ba, nil
}

func (ba *ByteArray) Byte(n int) (uint8, error) {
	if n >= len(ba.bytes) {
		return 0, fmt.Errorf(
			"attempt to read byte %d (length %d)",
			n,
			len(ba.bytes),
		)
	}
	return ba.bytes[n], nil
}

func (ba *ByteArray) Bytes(n int, l int) ([]uint8, error) {
	if (n + l) > len(ba.bytes) {
		return []uint8{}, fmt.Errorf(
			"attempt to read %d bytes from %d (length %d)",
			l,
			n,
			len(ba.bytes),
		)
	}
	return ba.bytes[n : n+l], nil
}

func (ba *ByteArray) SetByte(n int, v uint8) error {
	if n >= len(ba.bytes) {
		return fmt.Errorf(
			"attempt to set byte %d (length %d)",
			n,
			len(ba.bytes),
		)
	}
	ba.bytes[n] = v
	return nil
}

func (ba *ByteArray) SetBytes(n int, vl []uint8) error {
	if (n + len(vl)) > len(ba.bytes) {
		return fmt.Errorf(
			"attempt to set %d bytes from %d (length %d)",
			len(vl),
			n,
			len(ba.bytes),
		)
	}
	for i := range vl {
		ba.bytes[n+i] = vl[i]
	}
	return nil
}

func (ba *ByteArray) PushByte(v uint8) {
	ba.bytes = append(ba.bytes, v)
}

func (ba *ByteArray) PushBytes(vl []uint8) {
	ba.bytes = append(ba.bytes, vl...)
}

func (ba *ByteArray) Trim() error {
	newBa := make([]uint8, len(ba.bytes))
	copy(newBa, ba.bytes)
	ba.bytes = newBa
	return nil
}

func (ba *ByteArray) PushNByte(v uint32) {
	if v < 128 {
		ba.PushByte(uint8(v + 128))
	} else {
		mod := v % 128
		ba.PushByte(uint8(mod))
		ba.PushNByte(v >> 7)
	}
}

func (ba *ByteArray) PushNBytes(vl []uint32) {
	for i := range vl {
		ba.PushNByte(vl[i])
	}
}

func (ba *ByteArray) NByte(n int) (uint32, error) {
	if n >= len(ba.bytes) {
		return 0, fmt.Errorf(
			"attempt to read byte %d for NByte(length %d)",
			n,
			len(ba.bytes),
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

func (ba *ByteArray) PushCountedString(s string) {
	sA := []byte(s)
	ba.PushByte(uint8(len(sA)))
	ba.PushBytes(sA)
}

func (ba *ByteArray) CountedString(n int) (string, error) {
	sLength, err := ba.Byte(n)
	if err != nil {
		return "", err
	}
	sA, err := ba.Bytes(n+1, int(sLength))
	if err != nil {
		return "", err
	}
	return string(sA), nil
}

func (ba *ByteArray) CountedStrings() ([]string, error) {
	var strs []string
	var sLength uint8 = 0
	var err error
	for n := 0; (n + int(sLength)) < len(ba.bytes); n += int(sLength) + 1 {
		sLength, err = ba.Byte(n)
		if err != nil {
			return nil, err
		}
		sA, err := ba.Bytes(n+1, int(sLength))
		if err != nil {
			return nil, err
		}
		strs = append(strs, string(sA))
	}
	return strs, nil
}

func (ba *ByteArray) Clear() {
	ba.bytes = nil
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

func (ba *ByteArray) UnmarshalJSON(b []byte) error {
	var base64Str string
	err := json.Unmarshal(b, &base64Str)
	if err != nil {
		return err
	}
	*ba, err = NewByteArrayFromBase64(base64Str)
	if err != nil {
		return err
	}
	err = ba.Trim()
	return err
}

func (ba *ByteArray) MarshalJSON() ([]byte, error) {
	var dst []byte
	b64.StdEncoding.Encode(dst, ba.bytes)
	return dst, nil
}

func (ba *ByteArray) DeleteItem(n int) error {
	sLength, err := ba.Byte(n)
	if err != nil {
		return err
	}
	itemLength := int(sLength & 0x0000003F)
	if len(ba.bytes) > n {
		remainingBytes := make([]uint8, len(ba.bytes)-(n+itemLength))
		copy(remainingBytes, ba.bytes[n+itemLength:])
		err = ba.SetBytes(n, remainingBytes)
		if err != nil {
			return err
		}
		ba.bytes = ba.bytes[:len(ba.bytes)-n]
	}
	return nil
}
