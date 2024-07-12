package provider_v2

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
	return fmt.Sprintf("Total %d objects [%d Bytes], %d skipped objects, %d failed objects.",
		s.Objects, s.ObjectsBytes, s.Skipped, s.Failed)
}

func (s *Stats) Merge(other *Stats) {
	s.Objects += other.Objects
	s.ObjectsBytes += other.ObjectsBytes
	s.Failed += other.Failed
	s.Skipped += other.Skipped
	s.SkippedBytes += other.SkippedBytes
}
