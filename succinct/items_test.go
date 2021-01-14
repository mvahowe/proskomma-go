package succinct

import (
	"testing"
)

func checkHeaderBytes(t *testing.T, ba *ByteArray) {
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
	if pos != len(ba.bytes) {
		t.Errorf("last itemLength should point one past byteArray (%d/%d)", pos, len(ba.bytes))
	}
}

func TestHeaderBytes(t *testing.T) {
	succinctString := "AwCvAwKAAwCJAwKABABqgQMCgAMA9QMBgQMCgAMAgQMCgAMAqQMCgAQAdYEDAoADAMUDAoAEAAOCAwKABAAGggMBgQMCgAMAgQMCgAQACIIDAYIDAoDDBYjDBoiDBomDBYkDALoDAoADAIcDAoADAIMDAoADAIUDAYEDAoADAYMEAAmCAwGBAwKAAwDeAwKABABBgwMCgAMAmgMCgAMA+QMCgAMAlQMCgAMAhAMCgAMAgwMCgAMAiwMCgAMAtQMBggMBhMMFicMGiQ=="
	ba, _ := NewByteArrayFromBase64(succinctString)
	checkHeaderBytes(t, &ba)
}
