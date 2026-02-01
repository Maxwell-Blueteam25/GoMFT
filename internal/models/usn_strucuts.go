package models

type UsnRecordV2 struct {
	RecordLength    uint32
	MajorVersion    uint16
	MinorVersion    uint16
	FileReference   uint64
	ParentReference uint64
	Usn             int64
	TimeStamp       int64
	Reason          uint32
	SourceInfo      uint32
	SecurityId      uint32
	FileAttributes  uint32
	FileNameLength  uint16
	FileNameOffset  uint16
}

const (
	USN_REASON_DATA_OVERWRITE    = 0x00000001
	USN_REASON_DATA_EXTEND       = 0x00000002
	USN_REASON_DATA_TRUNCATION   = 0x00000004
	USN_REASON_FILE_CREATE       = 0x00000100
	USN_REASON_FILE_DELETE       = 0x00000200
	USN_REASON_RENAME_OLD_NAME   = 0x00001000
	USN_REASON_RENAME_NEW_NAME   = 0x00002000
	USN_REASON_BASIC_INFO_CHANGE = 0x00008000
	USN_REASON_CLOSE             = 0x80000000
)
