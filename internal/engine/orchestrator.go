package engine

import (
	"GoMFT/internal/models"
	"GoMFT/internal/output"
	"GoMFT/internal/parser"
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

type Orchestrator struct {
	MftPath     string
	JournalPath string
	PathMap     map[uint64]string
	Correlator  *Correlator
	Writer      *output.JSONWriter
	Timeline    *output.Timeline
}

type tempNode struct {
	ParentFRN uint64
	Name      string
}

func NewOrchestrator(mftPath, journalPath, outputPath string) *Orchestrator {
	writer, _ := output.NewJSONWriter(outputPath)
	return &Orchestrator{
		MftPath:     mftPath,
		JournalPath: journalPath,
		PathMap:     make(map[uint64]string),
		Correlator:  NewCorrelator(),
		Writer:      writer,
		Timeline:    output.NewTimeline(),
	}
}

func windowsTimeToGo(ft int64) time.Time {
	const epochDiff = 116444736000000000
	nano := (ft - epochDiff) * 100
	return time.Unix(0, nano)
}

func (o *Orchestrator) BuildPathMap() error {

	reader, err := parser.NewMftReader(o.MftPath)
	if err != nil {
		return fmt.Errorf("failed to open MFT: %v", err)
	}
	defer reader.FileHandler.Close()

	relMap := make(map[uint64]tempNode)
	var mftIndex uint64 = 0

	for {
		data, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			mftIndex++
			continue
		}

		header := parser.ParseMftHeader(data)
		if header.Signature == 0x454C4946 {
			if header.Flags&0x01 != 0 {
				parent, name := parser.GetFileNameAttribute(data, header.FirstAttributeOffset)
				if name != "" {
					relMap[mftIndex] = tempNode{ParentFRN: parent, Name: name}
				}
			}
		}
		mftIndex++
	}
	for frn := range relMap {
		fullPath := o.resolvePath(frn, relMap)
		o.PathMap[frn] = fullPath
	}
	return nil
}

func (o *Orchestrator) resolvePath(frn uint64, relMap map[uint64]tempNode) string {
	index := frn & 0x0000FFFFFFFFFFFF
	if index == 5 {
		return "."
	}
	node, exists := relMap[index]
	if !exists {
		return "$Orphan"
	}
	parentPath := o.resolvePath(node.ParentFRN, relMap)
	return parentPath + "\\" + node.Name
}

func (o *Orchestrator) Run() error {
	jReader, err := parser.NewJournalReader(o.JournalPath)
	if err != nil {
		return err
	}
	defer o.Writer.Close()

	fmt.Println("[*] Phase 2: Streaming USN Journal...")
	lastRender := time.Now()

	renderInterval := 100 * time.Millisecond

	for {
		n, err := jReader.ReadChunk()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		buffer := jReader.Buffer[:n]
		offset := 0

		for offset < len(buffer) {

			if time.Since(lastRender) > renderInterval {
				o.Timeline.RenderLive()
				lastRender = time.Now()
			}

			if offset+4 > len(buffer) {
				break
			}

			recordLen := binary.LittleEndian.Uint32(buffer[offset : offset+4])
			if recordLen == 0 {
				offset += 4
				continue
			}

			if uint64(offset)+uint64(recordLen) > uint64(len(buffer)) {
				break
			}

			if recordLen < 60 {
				offset += int(recordLen)
				continue
			}

			record, nameStr := parser.ParseUsnRecord(buffer[offset : offset+int(recordLen)])

			path := o.PathMap[record.FileReference&0x0000FFFFFFFFFFFF]
			if path == "" {
				path = "$Unknown"
			}

			timestamp := windowsTimeToGo(record.TimeStamp)

			isRename := (record.Reason&models.USN_REASON_RENAME_NEW_NAME != 0)
			isDelete := (record.Reason&models.USN_REASON_FILE_DELETE != 0)

			o.Timeline.AddEvent(timestamp.UnixNano(), isRename, isDelete, false)

			if record.Reason&models.USN_REASON_RENAME_OLD_NAME != 0 {
				o.Correlator.AddPending(record.FileReference, nameStr, timestamp)
			}

			if record.Reason&models.USN_REASON_RENAME_NEW_NAME != 0 {
				event, found := o.Correlator.ResolveRename(record.FileReference, nameStr, timestamp)
				if found {
					fullEvent := models.FileLifecycle{
						FRN:       record.FileReference,
						FullPath:  path,
						Renames:   []models.RenameEvent{event},
						IsActive:  true,
						IsRenamed: true,
					}
					o.Writer.Write(fullEvent)
				}
			}

			offset += int(record.RecordLength)
		}
	}

	o.Timeline.RenderLive()
	return nil
}
