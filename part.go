package m3u8

import "github.com/pkg/errors"

type PartInfSegment struct {
	PartTartget float64
}

func NewPartInf(line string) (*PartInfSegment, error) {
	item := parseLine(line[len(ExtPartInf+":"):])

	target, err := extractFloat64(item, PARTTARGET)
	if err != nil {
		return nil, errors.Wrap(err, "extractFloat64 err")
	}

	return &PartInfSegment{
		PartTartget: target,
	}, nil
}

func (ps *PartInfSegment) String() string {
	return "PartInfSegment"
}

type PartSegment struct {
	Duration    float64
	URI         string
	Independent bool
	ByteRange   *ByteRangeSegment
	Gap         bool
}

func NewPart(line string) (*PartSegment, error) {
	item := parseLine(line[len(ExtPart+":"):])

	duration, err := extractFloat64(item, DURATION)
	if err != nil {
		return nil, errors.Wrap(err, "extractFloat64 err")
	}

	independent, err := extractBool(item, INDEPENDENT)
	if err != nil {
		return nil, errors.Wrap(err, "extractBool err")
	}

	br, err := NewByteRange(item[ByteRange])
	if err != nil {
		return nil, errors.Wrap(err, "new byte range")
	}

	gap, err := extractBool(item, GAP)
	if err != nil {
		return nil, errors.Wrap(err, "extractBool err")
	}

	return &PartSegment{
		Duration:    duration,
		URI:         item[URI],
		Independent: independent,
		ByteRange:   br,
		Gap:         gap,
	}, nil
}

func (ps *PartSegment) String() string {
	return "PartSegment"
}
