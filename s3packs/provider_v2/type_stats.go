package provider_v2

type Stats struct {
	Objects int
	Failed  int
	Skipped int

	Bytes int64
}

func (s *Stats) IncObjects(n int, b int64) {
	s.Objects += n
	s.Bytes += b
}

func (s *Stats) IncFailed(n int, b int64) {
	s.Failed += n
	s.Bytes += b
}

func (s *Stats) IncSkipped(n int, b int64) {
	s.Skipped += n
	s.Bytes += b
}
