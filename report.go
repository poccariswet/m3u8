package m3u8

type RenditionReportSegment struct {
	URI      string
	LastMSN  int
	LastPART int
}

func NewReport(line string) (*RenditionReportSegment, error) {
	return &RenditionReportSegment{}, nil
}

func (rs *RenditionReportSegment) String() string {
	return "RenditionReportSegment"
}
