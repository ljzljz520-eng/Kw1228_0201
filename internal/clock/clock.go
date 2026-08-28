package clock

type Clock interface {
	Next() int64
}

type Sequence struct {
	value int64
}

func New(start int64) *Sequence {
	if start < 0 {
		start = 0
	}
	return &Sequence{value: start}
}

func (s *Sequence) Next() int64 {
	s.value++
	return s.value
}

func (s *Sequence) Current() int64 {
	return s.value
}

func (s *Sequence) Advance(count int64) int64 {
	if count > 0 {
		s.value += count
	}
	return s.value
}

type Fixed struct {
	value int64
}

func NewFixed(value int64) *Fixed {
	return &Fixed{value: value}
}

func (f *Fixed) Next() int64 {
	return f.value
}

func (f *Fixed) Set(value int64) {
	f.value = value
}
