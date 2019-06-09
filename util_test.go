package m3u8_test

import (
	"fmt"
	"strings"
	"testing"

	m3u8 "github.com/poccariswet/m3u8-decoder"
)

func parseLine(line string) map[string]string {
	m := map[string]string{}
	lines := strings.Split(line, ",")

	// if val has multiple items value, map's tmp key put in the value
	var tmp string
	for _, v := range lines {
		val := strings.Split(v, "=")
		if len(val) != 2 {
			str := m[tmp]
			m[tmp] = fmt.Sprintf("%s,%s", str, strings.Trim(val[0], `"`))
		} else {
			tmp = val[0]
			m[val[0]] = strings.Trim(val[1], `"`)
		}
	}

	return m
}

func TestUtil(t *testing.T) {
	line := `#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=195023,CODECS="avc1.42e00a,mp4a.40.2",AUDIO="audio"`

	values := parseLine(line[len(m3u8.ExtStreamInf+":"):])
	fmt.Println(values)
}
