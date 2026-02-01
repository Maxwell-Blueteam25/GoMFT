package engine

import (
	"GoMFT/internal/models"
	"time"
)

type Correlator struct {
	Cache map[uint64]PendingRename
}

type PendingRename struct {
	OldName   string
	Timestamp time.Time
}

func NewCorrelator() *Correlator {
	return &Correlator{
		Cache: make(map[uint64]PendingRename),
	}
}

func (c *Correlator) AddPending(frn uint64, oldName string, timestamp time.Time) {
	r := PendingRename{
		OldName:   oldName,
		Timestamp: timestamp,
	}
	c.Cache[frn] = r
}

func (c *Correlator) ResolveRename(frn uint64, newName string, timestamp time.Time) (models.RenameEvent, bool) {
	pending, ok := c.Cache[frn]
	if !ok {
		return models.RenameEvent{}, false
	}
	event := models.RenameEvent{
		OldName:   pending.OldName,
		NewName:   newName,
		Timestamp: timestamp,
	}

	delete(c.Cache, frn)

	return event, true
}
