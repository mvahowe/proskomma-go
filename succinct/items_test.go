package succinct

import (
	"testing"
)

func TestHeaderBytes(t *testing.T) {
	succinctString := "AwCvAwKAAwCJAwKABABqgQMCgAMA9QMBgQMCgAMAgQMCgAMAqQMCgAQAdYEDAoADAMUDAoAEAAOCAwKABAAGggMBgQMCgAMAgQMCgAQACIIDAYIDAoDDBYjDBoiDBomDBYkDALoDAoADAIcDAoADAIMDAoADAIUDAYEDAoADAYMEAAmCAwGBAwKAAwDeAwKABABBgwMCgAMAmgMCgAMA"
	ba := NewByteArray(256)
	ba, err := NewByteArrayFromBase64(succinctString)
	if err != nil {
		t.Errorf("NewByteArrayFromBase64 threw error: %s", err)
	}
	pos := 0
	for pos < len(ba.bytes) {
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
	if pos != len(ba.bytes)+1 {
		t.Errorf("last itemLength should point one past usedBytes")
	}
}
