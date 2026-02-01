package parser

import (
	"GoMFT/internal/models"
	"encoding/binary"
	"unicode/utf16"
)

func ParseMftHeader(data []byte) models.MftRecordHeader {
	return models.MftRecordHeader{
		Signature:            binary.LittleEndian.Uint32(data[0:4]),
		SequenceNumber:       binary.LittleEndian.Uint16(data[16:18]),
		FirstAttributeOffset: binary.LittleEndian.Uint16(data[20:22]),
		Flags:                binary.LittleEndian.Uint16(data[22:24]),
	}
}

func GetFileNameAttribute(data []byte, offset uint16) (uint64, string) {
	currentOffset := uint32(offset)

	for currentOffset+8 < uint32(len(data)) {
		attrType := binary.LittleEndian.Uint32(data[currentOffset : currentOffset+4])
		attrLen := binary.LittleEndian.Uint32(data[currentOffset+4 : currentOffset+8])

		if attrType == 0xFFFFFFFF {
			break
		}
		if attrLen == 0 {
			break
		}

		if attrType == 0x30 {
			nonResident := data[currentOffset+8]
			if nonResident == 0 {
				infoOffset := binary.LittleEndian.Uint16(data[currentOffset+20 : currentOffset+22])
				bodyOffset := currentOffset + uint32(infoOffset)

				if bodyOffset+66 > uint32(len(data)) {
					break
				}

				rawParent := binary.LittleEndian.Uint64(data[bodyOffset : bodyOffset+8])
				parentFrn := rawParent & 0x0000FFFFFFFFFFFF

				nameLen := uint32(data[bodyOffset+64])
				nameStart := bodyOffset + 66

				if nameStart+(nameLen*2) > uint32(len(data)) {
					break
				}

				u16s := make([]uint16, nameLen)
				for i := uint32(0); i < nameLen; i++ {
					u16s[i] = binary.LittleEndian.Uint16(data[nameStart+(i*2) : nameStart+(i*2)+2])
				}

				return parentFrn, string(utf16.Decode(u16s))
			}
		}

		currentOffset += attrLen
	}

	return 0, ""
}

func ParseUsnRecord(data []byte) (models.UsnRecordV2, string) {

	rec := models.UsnRecordV2{
		RecordLength:    binary.LittleEndian.Uint32(data[0:4]),
		MajorVersion:    binary.LittleEndian.Uint16(data[4:6]),
		FileReference:   binary.LittleEndian.Uint64(data[8:16]),
		ParentReference: binary.LittleEndian.Uint64(data[16:24]),
		Usn:             int64(binary.LittleEndian.Uint64(data[24:32])),
		TimeStamp:       int64(binary.LittleEndian.Uint64(data[32:40])),
		Reason:          binary.LittleEndian.Uint32(data[40:44]),
	}

	nameLen := binary.LittleEndian.Uint16(data[56:58])
	nameOffset := binary.LittleEndian.Uint16(data[58:60])

	if uint32(nameOffset)+uint32(nameLen) > uint32(len(data)) {
		return rec, ""
	}

	nameBytes := data[nameOffset : nameOffset+nameLen]
	u16s := make([]uint16, len(nameBytes)/2)
	for i := 0; i < len(u16s); i++ {
		u16s[i] = binary.LittleEndian.Uint16(nameBytes[i*2 : i*2+2])
	}

	return rec, string(utf16.Decode(u16s))
}
