package output

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type Timeline struct {
	Events []int64
	Start  int64
	End    int64

	Renames    int
	Deletions  int
	Timestomps int
}

func NewTimeline() *Timeline {
	return &Timeline{
		Events: make([]int64, 0, 10000),
		Start:  math.MaxInt64,
		End:    0,
	}
}

func (t *Timeline) AddEvent(timestamp int64, isRename, isDelete, isTimestomp bool) {
	if timestamp < t.Start {
		t.Start = timestamp
	}
	if timestamp > t.End {
		t.End = timestamp
	}

	t.Events = append(t.Events, timestamp)

	if isRename {
		t.Renames++
	}
	if isDelete {
		t.Deletions++
	}
	if isTimestomp {
		t.Timestomps++
	}
}

func (t *Timeline) RenderLive() {

	fmt.Print("\033[H\033[2J")
	t.Render()
}

func (t *Timeline) Render() {
	if len(t.Events) == 0 {
		fmt.Println("[!] Waiting for events...")
		return
	}

	fmt.Println("\n[=== LIVE ACTIVITY DASHBOARD ===]")
	fmt.Printf("Range: %s  <-->  %s\n",
		time.Unix(0, t.Start).Format("15:04:05"),
		time.Unix(0, t.End).Format("15:04:05"))
	fmt.Printf("Total Events: %d | Renames: %d | Deletions: %d\n",
		len(t.Events), t.Renames, t.Deletions)
	fmt.Println("---------------------------------------------------------------")

	numBuckets := 20
	duration := t.End - t.Start
	if duration == 0 {
		duration = 1
	}
	bucketSize := duration / int64(numBuckets)
	buckets := make([]int, numBuckets)

	maxCount := 0
	for _, ts := range t.Events {
		offset := ts - t.Start
		idx := int(offset / bucketSize)
		if idx >= numBuckets {
			idx = numBuckets - 1
		}
		if idx < 0 {
			idx = 0
		}
		buckets[idx]++
		if buckets[idx] > maxCount {
			maxCount = buckets[idx]
		}
	}

	for i, count := range buckets {
		bucketTime := t.Start + (int64(i) * bucketSize)
		timeLabel := time.Unix(0, bucketTime).Format("15:04:05")

		barLen := 0
		if maxCount > 0 {
			barLen = (count * 40) / maxCount
		}

		bar := strings.Repeat("█", barLen)
		if barLen == 0 && count > 0 {
			bar = "▒"
		}

		anomaly := ""

		if len(t.Events) > 100 && count > (len(t.Events)/4) {
			anomaly = " <!- SPIKE"
		}

		fmt.Printf("[%s] %s (%d)%s\n", timeLabel, bar, count, anomaly)
	}
	fmt.Println("---------------------------------------------------------------")
	fmt.Println("Processing... (Ctrl+C to stop)")
}
