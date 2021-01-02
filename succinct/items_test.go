package succinct

import (
	"testing"
)

func TestHeaderBytes(t *testing.T) {
	succinctString := "AwCvAwKAAwCJAwKABABqgQMCgAMA9QMBgQMCgAMAgQMCgAMAqQMCgAQAdYEDAoADAMUDAoAEAAOCAwKABAAGggMBgQMCgAMAgQMCgAQACIIDAYIDAoDDBYjDBoiDBomDBYkDALoDAoADAIcDAoADAIMDAoADAIUDAYEDAoADAYMEAAmCAwGBAwKAAwDeAwKABABBgwMCgAMAmgMCgAMA"
	ba := NewByteArray(256)
	ba.fromBase64(succinctString)
	pos := 0
	for pos < ba.usedBytes {
		itemLength, itemType, itemSubtype, err := ba.headerBytes(pos)
		if err != nil {
			t.Errorf("headerBytes threw error: %s", err)
		}
		if itemType > len(itemStrings) {
			t.Errorf("Unexpected itemType %d", itemType)
		}
		if (itemType == 0) && (itemSubtype > len(tokenStrings)) {
			t.Errorf("Unexpected token subtype %d", itemSubtype)
		} else if (itemType >= 2) && (itemSubtype > len(scopeStrings)) {
			t.Errorf("Unexpected scope subtype %d", itemSubtype)
		}
		pos += itemLength
	}
	if pos != ba.usedBytes + 1 {
		t.Errorf("last itemLength should point one past usedBytes")
	}
}