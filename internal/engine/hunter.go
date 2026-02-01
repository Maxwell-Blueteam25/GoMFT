package engine

import (
	"GoMFT/internal/models"
	"time"
)

func IsTimeStomped(siCreation time.Time, fnCreation time.Time) bool {

	if siCreation.Add(2 * time.Second).Before(fnCreation) {
		return true
	}

	if siCreation.Nanosecond() == 0 {
		return true
	}
	return false
}

func IsPhantom(lifecycle models.FileLifecycle) bool {

	if lifecycle.IsActive {
		return false
	}

	if lifecycle.Death.IsZero() {
		return false
	}

	duration := lifecycle.Death.Sub(lifecycle.Birth)

	if duration < 2*time.Second {
		return true
	}

	return false
}
