package types

import (
	"encoding/hex"
	"fmt"
	"hash"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type LogEntries struct {
	Metadata Metadata   `json:"metadata" mapstructure:"metadata"`
	Entries  []LogEntry `json:"entries" mapstructure:"entries"`
	Events   []Event    `json:"events" mapstructure:"events"`
	// Status used as a callback for pullers/retrievers to be able to ack/noack based on status.
	// Should be able to support multiple calls without causing errors or issues, similar to context.CancelFunc
	// to support chunking of large entries structures.
	Status func(Status) `json:"-"`
}

func (entries *LogEntries) ReportStatus(status Status) {
	if entries.Status == nil {
		return
	}

	entries.Status(status)
}

func (entries *LogEntries) Append(logEntries ...LogEntry) {
	entries.Entries = append(entries.Entries, logEntries...)
}

func (entries *LogEntries) EntryCount() int {
	return len(entries.Entries)
}

func (entries *LogEntries) SizeInBytes() int {
	return sizeOfEntries(entries.Entries)
}

func (entries *LogEntries) ChunkBySize(size int) []*LogEntries {
	if size <= 0 {
		return []*LogEntries{entries}
	}

	return entries.chunkBySize(entries.Entries, size)
}

func (entries *LogEntries) ChunkByCount(count int) []*LogEntries {
	if count <= 0 {
		return []*LogEntries{entries}
	}

	return entries.chunkByCount(entries.Entries, count)
}

func sizeOfEntries(entries []LogEntry) int {
	var size int

	for _, entry := range entries {
		size += len(entry.Log)
	}

	return size
}

func (entries *LogEntries) chunkBySize(chunk []LogEntry, size int) []*LogEntries {
	if sizeOfEntries(chunk) < size {
		return []*LogEntries{{
			Status:   entries.Status,
			Metadata: entries.Metadata,
			Entries:  chunk,
		}}
	}

	middle := len(chunk) / 2 // nolint: gomnd
	left := chunk[:middle]
	right := chunk[middle:]

	return append(entries.chunkBySize(left, size), entries.chunkBySize(right, size)...)
}

func (entries *LogEntries) chunkByCount(chunk []LogEntry, count int) []*LogEntries {
	if len(chunk) < count {
		return []*LogEntries{{
			Status:   entries.Status,
			Metadata: entries.Metadata,
			Entries:  chunk,
		}}
	}

	middle := len(chunk) / 2 // nolint: gomnd
	left := chunk[:middle]
	right := chunk[middle:]

	return append(entries.chunkBySize(left, count), entries.chunkBySize(right, count)...)
}

type Event map[string]any

type LogEntry struct {
	Log       string `json:"log,omitempty" mapstructure:"log"`
	Signature string `json:"signature,omitempty" mapstructure:"signature"`
	// Deprecated
	Timestamp time.Time `json:"timestamp,omitempty" mapstructure:"timestamp"`
}

func CreateLogEntries() *LogEntries {
	return &LogEntries{
		Metadata: CreateMetadata(),
		Entries:  make([]LogEntry, 0),
	}
}

// CreateStructuredLogEntry is a constructor/helper that automatically hashes the log.
func CreateStructuredLogEntry(log any, timestamp time.Time, hash hash.Hash) (LogEntry, error) {
	jsonLog, err := jsoniter.Marshal(log)
	if err != nil {
		return LogEntry{}, fmt.Errorf("failed to serialize log entry: %w", err)
	}

	return CreateUnstructuredLogEntry(string(jsonLog), timestamp, hash), nil
}

// CreateUnstructuredLogEntry is a constructor/helper that automatically hashes the log.
func CreateUnstructuredLogEntry(log string, timestamp time.Time, hash hash.Hash) LogEntry {
	entry := LogEntry{
		Log:       log,
		Timestamp: timestamp,
	}

	entry.Signature = hashLog([]byte(log), hash)

	return entry
}

func hashLog(log []byte, hash hash.Hash) string {
	// Implementations of hash cannot return errors, so we skip checking for one
	_, _ = hash.Write(log)

	return hex.EncodeToString(hash.Sum(nil))
}

func (logEntry *LogEntry) IsEmpty() bool {
	return len(logEntry.Log) <= 0
}
