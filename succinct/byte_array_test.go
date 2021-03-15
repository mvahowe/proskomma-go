package succinct

import (
	"testing"
)

func TestConstructor(t *testing.T) {
	ba := NewByteArray(32)
	if cap(ba.bytes) != 32 {
		t.Errorf("bytes capacity for new ByteArray is %d not 32", len(ba.bytes))
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
	ba.PushByte(99)
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
	ba.PushByte(93)
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
	ba.PushBytes([]uint8{2, 4, 6, 8})
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
	ba.PushBytes([]uint8{2, 4, 6, 8})
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

func TestNByte(t *testing.T) {
	ba := NewByteArray(5)
	ba.PushNByte(127)
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
	ba.PushNByte(130)
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

func TestNBytes(t *testing.T) {
	ba := NewByteArray(1)
	tValues := []uint32{127, 17000, 130}
	for i := range tValues {
		ba.PushNByte(tValues[i])
	}
	nBytes, err := ba.NBytes(0, len(tValues))
	if err != nil {
		t.Errorf("NBytes threw error: %s", err)
	}
	if len(nBytes) != len(tValues) {
		t.Errorf("nBytes expected to be length %d but was %d", len(tValues), len(nBytes))
	}
	for i := range tValues {
		if nBytes[i] != tValues[i] {
			t.Errorf("nBytes[%d] expected to be %d but was %d", i, tValues[i], nBytes[i])
		}
	}

	nBytesEmpty, err := ba.NBytes(0, 0)
	if err != nil {
		t.Errorf("NBytes threw error: %s", err)
	}
	if len(nBytesEmpty) != 0 {
		t.Errorf("nBytesEmpty expected to be length 0 but was %d", len(nBytesEmpty))
	}

	_, err = ba.NBytes(0, len(tValues)+1)
	if err == nil {
		t.Errorf("NBytes was expected to throw error with invalid input but did not")
	}
}

func testCountedString(t *testing.T, testString string) {
	ba := NewByteArray(32)
	ba.PushCountedString(testString)
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

func testCountedStrings(t *testing.T, testStrings []string) {
	ba := NewByteArray(32)
	for _, testString := range testStrings {
		ba.PushCountedString(testString)
	}
	v, err := ba.Byte(0)
	if err != nil {
		t.Errorf("reading first Count Byte threw error: %s", err)
	}
	if v != uint8(len(testStrings[0])) {
		t.Errorf(
			"String length of first Counted String should be %d, not %d",
			len(testStrings),
			v,
		)
	}
	strs, err := ba.CountedStrings()
	if err != nil {
		t.Errorf("CountedStrings after 1st PushCountedString threw error: %s", err)
	}
	for i, s := range strs {
		if s != testStrings[i] {
			t.Errorf(
				"expected first string to be '%s', not '%s'",
				testStrings[i],
				s,
			)
		}
	}
}

func TestCountedString(t *testing.T) {
	testCountedString(t, "abc")
	testCountedString(t, "égale")
	testCountedString(t, "וּ⁠בְ⁠דֶ֣רֶך")
}

func TestCountedStrings(t *testing.T) {
	testCountedStrings(t, []string{"abcd", "efg", "hijklmnop"})
	testCountedStrings(t, []string{"égale"})
	testCountedStrings(t, []string{"וּ⁠בְ⁠דֶ֣רֶך", "égale", "abcd", "efg", "hijklmnop"})
}

func TestClear(t *testing.T) {
	ba := NewByteArray(32)
	ba.PushCountedString("abcde")
	ba.Clear()
	if len(ba.bytes) != 0 {
		t.Errorf(
			"usedBytes after Clear should be 0, not %d",
			len(ba.bytes),
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

func pow2(y int) int {
	ret := 1
	for y > 0 {
		ret *= 2
		y--
	}
	return ret
}

func TestNByteLength(t *testing.T) {
	ba := NewByteArray(32)
	testNByteLength(t, &ba, pow2(7)-1, 1)
	testNByteLength(t, &ba, pow2(7), 2)
	testNByteLength(t, &ba, pow2(14), 3)
	testNByteLength(t, &ba, pow2(21), 4)
}

func TestTrim(t *testing.T) {
	ba := NewByteArray(32)
	ba.PushCountedString("abcdef")
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
	ba.PushCountedString("abcde")
	ba2, err := NewByteArrayFromBase64(ba.base64())
	if err != nil {
		t.Errorf("NewByteArrayFromBase64 threw error: %s", err)
	}
	s, _ := ba2.CountedString(0)
	if s != "abcde" {
		t.Errorf(
			"Base64'd string should be '%s', not '%s'",
			"abcde",
			s,
		)
	}
}

func TestDeleteItem(t *testing.T) {
	ba := NewByteArray(1)
	lengthPos := len(ba.bytes)
	ba.PushByte(0)
	ba.PushByte(1)
	ba.PushNByte(299)
	ba.SetByte(lengthPos, uint8((len(ba.bytes)-lengthPos)|itemString2Int["token"]<<6))

	lengthPos = len(ba.bytes)
	ba.PushByte(0)
	ba.PushByte(10)
	ba.PushNByte(143)
	ba.SetByte(lengthPos, uint8((len(ba.bytes)-lengthPos)|itemString2Int["graft"]<<6))

	lengthPos = len(ba.bytes)
	ba.PushByte(0)
	ba.PushByte(2)
	ba.PushNByte(567)
	ba.SetByte(lengthPos, uint8((len(ba.bytes)-lengthPos)|3<<6))

	firstLength, err := ba.Byte(0)
	if err != nil {
		t.Errorf("error: %s", err)
	}
	firstLength = firstLength & 0x0000003F

	secondLength, err := ba.Byte(int(firstLength))
	if err != nil {
		t.Errorf("error: %s", err)
	}
	secondLength = secondLength & 0x0000003F

	thirdLength, err := ba.Byte(int(firstLength + secondLength))
	if err != nil {
		t.Errorf("error: %s", err)
	}
	thirdLength = thirdLength & 0x0000003F

	fullLength := int(firstLength + secondLength + thirdLength)
	if fullLength != (len(ba.bytes)) {
		t.Errorf("ByteArray expected to be length %d but was %d", fullLength, len(ba.bytes))
	}

	ba.DeleteItem(int(firstLength))

	newFirstLength, err := ba.Byte(0)
	if err != nil {
		t.Errorf("error: %s", err)
	}
	newFirstLength = newFirstLength & 0x0000003F

	newSecondLength, err := ba.Byte(int(newFirstLength))
	if err != nil {
		t.Errorf("error: %s", err)
	}
	newSecondLength = newSecondLength & 0x0000003F

	if secondLength != newFirstLength {
		t.Errorf("secondLength expected to be %d but was %d", newFirstLength, firstLength)
	}

	if thirdLength != newSecondLength {
		t.Errorf("thirdLength expected to be %d but was %d", newSecondLength, thirdLength)
	}

	newFullLength := int(newFirstLength + newSecondLength)
	if newFullLength != len(ba.bytes) {
		t.Errorf("newFullLength expected to be %d but was %d", len(ba.bytes), newFullLength)
	}
}
