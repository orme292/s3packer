package objectify

type Stats struct {
	Bytes    int64
	Uploaded int
	Ignored  int
	Objects  int
}

func (s *Stats) Add(s2 Stats) {
	s.Bytes += s2.Bytes
	s.Uploaded += s2.Uploaded
	s.Ignored += s2.Ignored
	s.Objects += s2.Objects
}

const (
	EmptyString = ""
)
