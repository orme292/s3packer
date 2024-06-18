package provider_v2

type Stats struct {
	Objects int
	Bytes   int64
}

func (s *Stats) Increment(n int, b int64) {
	s.Objects += n
	s.Bytes += b
}
