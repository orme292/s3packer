package provider

import (
	"fmt"
)

type Stats struct {
	Objects      int
	ObjectsBytes int64
	Failed       int
	Skipped      int
	SkippedBytes int64
}

func (s *Stats) IncObjects(n int, b int64) {
	s.Objects += n
	s.ObjectsBytes += b
}

func (s *Stats) IncFailed(n int, b int64) {
	s.Failed += n
}

func (s *Stats) IncSkipped(n int, b int64) {
	s.Skipped += n
	s.SkippedBytes += b
}

func (s *Stats) String() string {
	return fmt.Sprintf("%d objects uploaded, %d skipped, and %d failed.",
		s.Objects, s.Skipped, s.Failed)
}

func (s *Stats) Merge(other *Stats) {
	if other != nil {
		s.Objects += other.Objects
		s.ObjectsBytes += other.ObjectsBytes
		s.Failed += other.Failed
		s.Skipped += other.Skipped
		s.SkippedBytes += other.SkippedBytes
	}
}

func (s *Stats) ReadableString() map[int64]string {
	humanStr := make(map[int64]string)
	humanStr[s.ObjectsBytes] = humanReadableByteCount(s.ObjectsBytes)
	humanStr[s.SkippedBytes] = humanReadableByteCount(s.SkippedBytes)
	return humanStr
}

func humanReadableByteCount(b int64) string {

	unit := int64(1024)
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}

	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.2f %ciB", float64(b)/float64(div), "KMGTPE"[exp])

}
