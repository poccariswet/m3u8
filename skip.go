package m3u8

type SkipSegment struct {
	SkippedSegments uint64
}

func NewSkip(line string) (*SkipSegment, error) {
	return &SkipSegment{}, nil
}

func (ss *SkipSegment) String() string {
	return "SkipSegment"
}
