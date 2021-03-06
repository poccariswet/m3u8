package m3u8

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

// decode parses a playlist
func decode(buf *bytes.Buffer) (*Playlist, error) {
	playlist := NewPlaylist()
	var end bool
	states := new(States)

	for !end {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			end = true
		} else if err != nil {
			return nil, err
		}

		if len(line) < 1 || line == "\r" {
			continue
		}

		line = strings.TrimSpace(line)
		if err := decodeLine(playlist, line, states); err != nil {
			return playlist, err
		}
	}

	return playlist, nil
}

// DecodeFrom read a playlist passed from the io.Reader
func DecodeFrom(r io.Reader) (*Playlist, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	return decode(buf)
}

// ReadFile reads contents from filepath and return Playlist
func ReadFile(path string) (*Playlist, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "ReadFile err")
	}

	return decode(bytes.NewBuffer(file))
}

// decodeLine decodes a line of playlist and parses
func decodeLine(p *Playlist, line string, s *States) error {
	if !s.m3u8 && line != EXTM3U {
		return errors.New("invalid playlist, not exist #EXTM3U")
	}

	switch {
	case line == EXTM3U:
		s.m3u8 = true
	case strings.HasPrefix(line, ExtENDList):
		p.live = false
	case strings.HasPrefix(line, ExtVersion):
		p.hasVersion = true
		_, err := fmt.Sscanf(line, ExtVersion+":%d", &p.Version)
		if err != nil {
			return errors.Wrap(err, "invalid scan version")
		}
	case strings.HasPrefix(line, EXTINF):
		inf, err := NewExtInf(line)
		if err != nil {
			return errors.Wrap(err, "new extinf err")
		}
		p.master = false
		s.segment = inf
		s.segmentTag = true
	case strings.HasPrefix(line, ExtMedia):
		m, err := NewMedia(line)
		if err != nil {
			return errors.Wrap(err, "new media err")
		}
		p.AppendSegment(m)
	case strings.HasPrefix(line, ExtStreamInf):
		p.master = true
		s.segmentTag = true
		line = line[len(ExtStreamInf+":"):]
		v, err := NewVariant(line)
		if err != nil {
			return errors.Wrap(err, "new variant err")
		}
		s.segment = v
	case strings.HasPrefix(line, ExtFrameStreamInf):
		p.master = true
		s.segmentTag = false
		line = line[len(ExtFrameStreamInf+":"):]
		v, err := NewVariant(line)
		if err != nil {
			return errors.Wrap(err, "new variant err")
		}
		v.IFrame = true
		s.segment = v
		p.AppendSegment(v)
	case strings.HasPrefix(line, ExtByteRange):
		br, err := NewByteRange(line)
		if err != nil {
			return errors.Wrap(err, "new byte range err")
		}
		br.Extflag = true
		if m, has := s.segment.(*MapSegment); has {
			m.ByteRange = br
			s.segment = m
			br.Extflag = false
		} else if inf, has := s.segment.(*InfSegment); has {
			inf.ByteRange = br
			s.segment = inf
		}
	case strings.HasPrefix(line, ExtMap):
		m, err := NewMap(line)
		if err != nil {
			return errors.Wrap(err, "new map err")
		}
		p.AppendSegment(m)
	case strings.HasPrefix(line, ExtKey):
		key, err := NewKey(line)
		if err != nil {
			return errors.Wrap(err, "new key err")
		}
		p.AppendSegment(key)
	case strings.HasPrefix(line, ExtProgramDateTime):
		dt, err := NewProgramDateTime(line)
		if err != nil {
			return errors.Wrap(err, "new program date time err")
		}
		p.AppendSegment(dt)
	case strings.HasPrefix(line, ExtDateRange):
		dr, err := NewDateRange(line)
		if err != nil {
			return errors.Wrap(err, "new date range err")
		}
		p.AppendSegment(dr)

	/* low-latency tags */
	case strings.HasPrefix(line, ExtServerControl):
		sc, err := NewServerControl(line)
		if err != nil {
			return errors.Wrap(err, "new server control err")
		}
		p.AppendSegment(sc)
	case strings.HasPrefix(line, ExtPartInf):
		pi, err := NewPartInf(line)
		if err != nil {
			return errors.Wrap(err, "new part inf err")
		}
		p.AppendSegment(pi)
	case strings.HasPrefix(line, ExtRenditionReport):
		report, err := NewReport(line)
		if err != nil {
			return errors.Wrap(err, "new rendition report err")
		}
		p.AppendSegment(report)
	case strings.HasPrefix(line, ExtSkip):
		skip, err := NewSkip(line)
		if err != nil {
			return errors.Wrap(err, "new skip err")
		}
		p.AppendSegment(skip)
	case strings.HasPrefix(line, ExtPart):
		part, err := NewPart(line)
		if err != nil {
			return errors.Wrap(err, "new part err")
		}
		p.AppendSegment(part)

	/* session tags */
	case strings.HasPrefix(line, ExtSessionKey):
		sk, err := NewSessionKey(line)
		if err != nil {
			return errors.Wrap(err, "new session key err")
		}
		p.AppendSegment(sk)
	case strings.HasPrefix(line, ExtSessionData):
		sd, err := NewSessionData(line)
		if err != nil {
			return errors.Wrap(err, "new session data err")
		}
		p.AppendSegment(sd)

	case strings.HasPrefix(line, ExtStart):
		start, err := NewStart(line)
		if err != nil {
			return errors.Wrap(err, "new start err")
		}
		p.AppendSegment(start)
	case strings.HasPrefix(line, ExtIndependentSegments):
		p.IndependentSegments = true

		/* playlist tags */
	case strings.HasPrefix(line, ExtPlaylistType):
		_, err := fmt.Sscanf(line, ExtPlaylistType+":%s", &p.PlaylistType)
		if err != nil {
			return errors.Wrap(err, "invalid scan playlist type")
		}
	case strings.HasPrefix(line, ExtIFramesOnly):
		p.IFrameOnly = true
	case strings.HasPrefix(line, ExtTargetDutation):
		_, err := fmt.Sscanf(line, ExtTargetDutation+":%f", &p.TargetDuration)
		if err != nil {
			return errors.Wrap(err, "invalid scan TargetDuration")
		}
	case strings.HasPrefix(line, ExtDiscontinuitySequence):
		_, err := fmt.Sscanf(line, ExtTargetDutation+":%d", &p.DiscontinuitySequence)
		if err != nil {
			return errors.Wrap(err, "invalid scan DiscontinuitySequence")
		}
	case strings.HasPrefix(line, ExtAllowCache):
		p.AllowCache = parseBool(line[len(ExtAllowCache+":"):])
	case strings.HasPrefix(line, ExtMediaSequence):
		_, err := fmt.Sscanf(line, ExtTargetDutation+":%d", &p.MediaSequence)
		if err != nil {
			return errors.Wrap(err, "invalid scan MediaSequence")
		}
	default:
		line = strings.Trim(line, "\n")
		uri := strings.TrimSpace(line)
		if s.segment != nil && s.segmentTag {
			if p.master {
				v, has := s.segment.(*VariantSegment)
				if !has {
					return errors.New("invalid variant playlist")
				}
				v.URI = uri
				p.AppendSegment(v)
			} else {
				i, has := s.segment.(*InfSegment)
				if !has {
					return errors.New("invalid EXTINF segment")
				}
				i.URI = uri
				p.AppendSegment(i)
			}
			s.segmentTag = false

			return nil
		}
	}
	return nil
}
