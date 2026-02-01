package models

import "time"

type FileLifecycle struct {
	FRN         uint64
	FullPath    string
	IsActive    bool
	IsDirectory bool
	IsRenamed   bool

	Birth time.Time
	Death time.Time

	Renames []RenameEvent
	Events  []UsnRecordV2

	Timestomped bool
	Phantom     bool
}

type RenameEvent struct {
	OldName   string
	NewName   string
	Timestamp time.Time
}
