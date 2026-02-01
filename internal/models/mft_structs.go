package models

type MftRecordHeader struct {
	Signature            uint32
	UpdateSequenceOffset uint16
	UpdateSequenceSize   uint16
	LogSequenceNumber    uint64
	SequenceNumber       uint16
	HardLinkCount        uint16
	FirstAttributeOffset uint16
	Flags                uint16
	RealSize             uint32
	AllocatedSize        uint32
	BaseRecordReference  uint64
	NextAttributeID      uint16
}

type AttributeHeader struct {
	TypeCode        uint32
	Length          uint32
	NonResidentFlag uint8
	NameLength      uint8
	NameOffset      uint16
	Flags           uint16
	AttributeID     uint16

	InfoLength uint32
	InfoOffset uint16
	IndexFlag  uint8
	Padding    uint8

	StartingVCN     uint64
	LastVCN         uint64
	RunArrayOffset  uint16
	CompressionUnit uint16
	PaddingNR       [4]byte
	AllocatedSize   uint64
	RealSize        uint64
	InitializedSize uint64
}

type StandardInformation struct {
	CreationTime         uint64
	ModificationTime     uint64
	MftModifiedTime      uint64
	AccessTime           uint64
	FileAttributes       uint32
	MaxVersions          uint32
	VersionNumber        uint32
	ClassId              uint32
	OwnerId              uint32
	SecurityId           uint32
	QuotaCharge          uint64
	UpdateSequenceNumber uint64
}

type FileNameAttribute struct {
	ParentDirectoryReference uint64
	CreationTime             uint64
	ModificationTime         uint64
	MftModifiedTime          uint64
	AccessTime               uint64
	AllocatedSize            uint64
	RealSize                 uint64
	Flags                    uint32
	ReparseTag               uint32
	NameLength               uint8
	Namespace                uint8
	FileName                 [255]uint16
}
