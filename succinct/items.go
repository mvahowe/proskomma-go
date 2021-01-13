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

/*
succinctScopeLabel
 */

func tokenChars(enums *Enums, succinct *ByteArray, itemSubType int, pos int) (string, error) {
	enumCategory := tokenCategory[itemSubType]
	var enumForToken EnumList
	if enumCategory == "wordLike" {
		enumForToken = enums.WordLike
	} else {
		enumForToken = enums.NotWordLike
	}
	itemIndex, err := succinct.NByte(pos + 2)
	if err != nil {
		return "", err
	}
	return enumForToken[itemIndex], nil
}

func graftName(enums *Enums, itemSubType int) (string, error) {
	return enums.GraftTypes[itemSubType], nil
}

func graftSeqId(enums *Enums, succinct *ByteArray, pos int) (string, error) {
	itemIndex, err := succinct.NByte(pos + 2)
	if err != nil {
		return "", err
	}
	return enums.IDs[itemIndex], nil
}
