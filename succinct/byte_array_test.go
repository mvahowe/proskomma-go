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

func TestNByte(t *testing.T) {
	ba := NewByteArray(5)
	err := ba.PushNByte(127)
	if err != nil {
		t.Errorf("PushNByte threw error: %s", err)
	}
	v, err := ba.Byte(0)
	if err != nil {
		t.Errorf("Byte for ByteArray after 1st PushNByte threw error: %s", err)
	}
	if v != (127 + 128) {
		t.Errorf("Oth byte after first NBytes should be 255, not %d", v)
	}
	bv, err := ba.NByte(0)
	if err != nil {
		t.Errorf("NByte threw error: %s", err)
	}
	if bv != 127 {
		t.Errorf("First NByte should be 127, not %d", bv)
	}
	err = ba.PushNByte(130)
	if err != nil {
		t.Errorf("Byte for ByteArray after 2nd PushNByte threw error: %s", err)
	}
	v, err = ba.Byte(1)
	if err != nil {
		t.Errorf("Byte for ByteArray after 2nd PushNByte threw error: %s", err)
	}
	vl, err := ba.Bytes(1, 2)
	if err != nil {
		t.Errorf("Bytes for ByteArray after 2nd PushNByte threw error: %s", err)
	}
	if vl[0] != 2 {
		t.Errorf("First byte after 2nd PushNByte should be 2, not %d", vl[0])
	}
	if vl[1] != (1 + 128) {
		t.Errorf("2nd byte after 2nd PushNByte should be 129, not %d", vl[1])
	}
	bv, err = ba.NByte(1)
	if err != nil {
		t.Errorf("2nd NByte threw error: %s", err)
	}
	if bv != 130 {
		t.Errorf("2nd NByte should be 130, not %d", bv)
	}
}

func testCountedString(t *testing.T, testString string) {
	ba := NewByteArray(32)
	err:= ba.PushCountedString(testString)
	if err != nil {
		t.Errorf("1st PushCountedString threw error: %s", err)
	}
	v, err := ba.Byte(0)
	if err != nil {
		t.Errorf("Byte after 1st PushCountedString threw error: %s", err)
	}
	if v != uint8(len(testString)) {
		t.Errorf(
			"String length after 1st PushCountedString should be %d, not %d",
			len(testString),
			v,
			)
	}
	s, err := ba.CountedString(0)
	if err != nil {
		t.Errorf("CountedString after 1st PushCountedString threw error: %s", err)
	}
	if s != testString {
		t.Errorf(
			"expected first string to be '%s', not '%s'",
			testString,
			s,
			)
	}
}

func TestCountedString(t *testing.T) {
	testCountedString(t, "abc")
	testCountedString(t, "égale")
	testCountedString(t, "וּ⁠בְ⁠דֶ֣רֶך")
}

func TestClear(t *testing.T) {
	ba := NewByteArray(32)
	_ = ba.PushCountedString("abcde")
	ba.Clear()
	if ba.usedBytes != 0 {
		t.Errorf(
			"usedBytes after Clear should be 0, not %d",
			ba.usedBytes,
		)
	}
}

func testNByteLength(t *testing.T, ba *ByteArray, v int, l int) {
	if ba.NByteLength(v) != l {
		t.Errorf(
			"NByteLength(%d) should be %d, not %d",
			v,
			l,
			ba.NByteLength(v),
			)
	}
}

func pow2 (y int) int {
	ret := 1
	for y > 0 {
		ret *= 2
		y--
	}
	return ret
}

func TestNByteLength(t *testing.T) {
	ba := NewByteArray(32)
	testNByteLength(t, &ba, pow2(7) - 1, 1)
	testNByteLength(t, &ba, pow2(7), 2)
	testNByteLength(t, &ba, pow2(14), 3)
	testNByteLength(t, &ba, pow2(21), 4)
}

func TestTrim(t *testing.T) {
	ba := NewByteArray(32)
	_ = ba.PushCountedString("abcdef")
	err := ba.Trim()
	if err != nil {
		t.Errorf("Trim threw error: %s", err)
	}
	if len(ba.bytes) != 7 {
		t.Errorf(
			"bytes length after Trim should be 7, not %d",
			len(ba.bytes),
			)
	}
}

func TestBase64(t *testing.T) {
	ba := NewByteArray(32)
	_ = ba.PushCountedString("abcde")
	ba2 := NewByteArray(32)
	ba2.fromBase64(ba.base64())
	s, _ := ba2.CountedString(0)
	if s != "abcde" {
		t.Errorf(
			"Base64'd string should be '%s', not '%s'",
			"abcde",
			s,
			)
	}
}