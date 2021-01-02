package succinct

func (ba *ByteArray) headerBytes(pos int) (int, int, int, error) {
	headerByte, err := ba.Byte(pos)
	if err != nil {
		return 0, 0, 0, err
	}
	itemType := headerByte >> 6
	itemLength := headerByte & 0x3F
	itemSubtype, err := ba.Byte(pos + 1)
	if err != nil {
		return int(itemLength), int(itemType), 0, err
	}
	return int(itemLength), int(itemType), int(itemSubtype), nil
}
