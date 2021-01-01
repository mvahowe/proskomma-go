package succinct

import (
	"testing"
)

func TestConstructor(t *testing.T) {
	ba := NewByteArray(32)
	if ba.usedBytes != 0 {
		t.Errorf("usedBytes for new ByteArray is %d not 0", ba.usedBytes)
	}
	if len(ba.bytes) != 32 {
		t.Errorf("bytes length for new ByteArray is %d not 32", len(ba.bytes))
	}
}

func TestReadPushByte(t *testing.T) {
	ba := NewByteArray(1)
	var v uint8
	var err error
	_, err = ba.Byte(0)
	if err == nil {
		t.Errorf("Accessing unset byte for ByteArray did not throw error")
	}
	err = ba.PushByte(99)
	if err != nil {
		t.Errorf("PushByte for ByteArray threw error: '%s'", err)
	}
	v, err = ba.Byte(0)
	if err != nil {
		t.Errorf("Accessing set byte for ByteArray threw error: %s", err)
	}
	if v != 99 {
		t.Errorf("Value of 0th byte should be 99, not %d", v)
	}
}

func TestWriteByte(t *testing.T) {
	ba := NewByteArray(1)
	var v uint8
	_ = ba.PushByte(93)
	v, _ = ba.Byte(0)
	if v != 93 {
		t.Errorf("0th byte for ByteArray should be 93 after PushByte")
	}
	err := ba.SetByte(0, 27)
	if err != nil {
		t.Errorf("accessing set byte for ByteArray threw error: %s", err)
	}
	v, _ = ba.Byte(0)
	if v != 27 {
		t.Errorf("Oth byte for ByteArray should be 27 after SetByte")
	}
	err = ba.SetByte(3, 27)
	if err == nil {
		t.Errorf("writing unset byte for ByteArray did not throw error")
	}
}

func TestPushReadBytes(t *testing.T) {
	ba := NewByteArray(10)
	err := ba.PushBytes([]uint8{2, 4, 6, 8})
	if err != nil {
		t.Errorf("PushBytes for ByteArray threw error: %s", err)
	}
	v, _ := ba.Byte(0)
	if v != 2 {
		t.Errorf("0th byte after PushBytes should be 2, not %d", v)
	}
	v, _ = ba.Byte(2)
	if v != 6 {
		t.Errorf("2nd byte after PushBytes should be 6, not %d", v)
	}
	vl, err := ba.Bytes(1, 3)
	if err != nil {
		t.Errorf("Bytes for ByteArray threw error: %s", err)
	}
	if vl[2] != 8 {
		t.Errorf("4th byte via Bytes should be 8, not %d", v)
	}
}

func TestSetBytes(t *testing.T) {
	ba := NewByteArray(10)
	_ = ba.PushBytes([]uint8{2, 4, 6, 8})
	v, _ := ba.Byte(2)
	if v != 6 {
		t.Errorf("2nd byte after PushBytes should be 6, not %d", v)
	}
	err := ba.SetBytes(1, []uint8{3, 5, 7})
	if err != nil {
		t.Errorf("SetBytes for ByteArray threw error: %s", err)
	}
	v, _ = ba.Byte(2)
	if v != 5 {
		t.Errorf("2nd byte after SetBytes should be 5, not %d", v)
	}
}

func TestGrow(t *testing.T) {
	ba := NewByteArray(5)
	_ = ba.PushBytes([]uint8{2, 4, 6, 8, 10})
	if ba.usedBytes != 5 {
		t.Errorf("usedBytes after initial push should be 5, not %d", ba.usedBytes)
	}
	if len(ba.bytes) != 5 {
		t.Errorf("Length after initial push should be 5, not %d", len(ba.bytes))
	}
	_ = ba.PushByte(12)
	if ba.usedBytes != 6 {
		t.Errorf("usedBytes after 2nd push should be 6, not %d", ba.usedBytes)
	}
	if len(ba.bytes) != 10 {
		t.Errorf("Length after 2nd push should be 10, not %d", len(ba.bytes))
	}
}