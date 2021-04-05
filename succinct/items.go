package succinct

import "strings"

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

func scopeLabel(enums *Enums, succinct *ByteArray, itemSubType int, pos int) (string, error) {
	nScopeBits := nComponentsForScope[itemSubType]
	offset := 2
	scopeBits := []string{scopeStrings[itemSubType]}
	for nScopeBits > 1 {
		scopeBitIndex, err := succinct.NByte(pos + 2)
		if err != nil {
			return "", err
		}
		scopeBit := enums.ScopeBits[scopeBitIndex]
		scopeBits = append(scopeBits, scopeBit)
		offset += succinct.NByteLength(int(scopeBitIndex))
		nScopeBits--
	}
	return strings.Join(scopeBits, "/"), nil
}

func (ba *ByteArray) pushSuccinctTokenBytes(tokenEnumIndex uint8, charsEnumIndex uint32) {
	lengthPos := len(ba.bytes)
	ba.PushByte(0)
	ba.PushByte(tokenEnumIndex)
	ba.PushNByte(charsEnumIndex)
	ba.SetByte(lengthPos, uint8((len(ba.bytes)-lengthPos)|itemString2Int["token"]<<6))
}

func (ba *ByteArray) pushSuccinctGraftBytes(graftTypeEnumIndex uint8, seqEnumIndex uint32) {
	lengthPos := len(ba.bytes)
	ba.PushByte(0)
	ba.PushByte(graftTypeEnumIndex)
	ba.PushNByte(seqEnumIndex)
	ba.SetByte(lengthPos, uint8((len(ba.bytes)-lengthPos)|itemString2Int["graft"]<<6))
}

func (ba *ByteArray) pushSuccinctScopeBytes(itemTypeByte uint8, scopeTypeByte uint8, scopeBitBytes []uint32) {
	lengthPos := len(ba.bytes)
	ba.PushByte(0)
	ba.PushByte(scopeTypeByte)
	ba.PushNBytes(scopeBitBytes)
	ba.SetByte(lengthPos, uint8((len(ba.bytes)-lengthPos)|int(itemTypeByte)<<6))
}

func (ba *ByteArray) unpackEnum() ([]string, error) {
	pos, count := 0, 0
	var s []string
	for pos < len(ba.bytes) {
		strLen, err := ba.Byte(pos)
		if err != nil {
			return nil, err
		}
		unpacked, err := ba.CountedString((pos))
		if err != nil {
			return nil, err
		}
		s = append(s, unpacked)
		pos += int(strLen + 1)
		count++
	}
	return s, nil
}
