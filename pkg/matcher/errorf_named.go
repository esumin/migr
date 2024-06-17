package matcher

import (
	"fmt"
	"regexp"
	"strings"
)

// example
// return nil, errors.Errorf("Found an unexpected number of volumes: volume_id=%s result_count=%d", id, len(vols))
func matchErrorfWithTwoNamedParams(line string) string {
	wrap := regexp.MustCompile(`(.*)errors.Errorf\("(.+)(?:[:.]\s*)([^=]+)=%[sd]{1}[,\s]+([^=]+)=%[sd]{1}\s*",\s*([^,]+),\s*([^,]+)\)(.*)`)
	match := wrap.FindStringSubmatch(line)
	if match == nil || strings.Contains(match[4], "%") || strings.Contains(match[4], "fmt.Sprintf") {
		return ""
	}

	return fmt.Sprintf("%serrkit.New(\"%s\", \"%s\", %s, \"%s\", %s)%s", match[1], match[2], match[3], match[5], match[4], match[6], match[7])
}

// example
// return nil, errors.Errorf("Required volume fields not available, volumeType: %s, Az: %s, VolumeTags: %v", snapshot.Volume.VolumeType, snapshot.Volume.Az, snapshot.Volume.Tags)",
func matchErrorfWithThreeNamedParams(line string) string {
	wrap := regexp.MustCompile(`(.*)errors.Errorf\("(.+)\s+([^:=]+)[:=]{1}\s*%[sdv]{1}[,\s]+([^:=]+)[:=]{1}\s*%[sdv]{1}[,\s]+([^:=]+)[:=]{1}\s*%[sdv]{1}\s*",\s*([^,]+),\s*([^,]+),\s*([^,]+)\)(.*)`)
	match := wrap.FindStringSubmatch(line)
	if match == nil || strings.Contains(match[4], "%") {
		return ""
	}

	msg := maybeTrimComma(match[2])

	return fmt.Sprintf("%serrkit.New(\"%s\", \"%s\", %s, \"%s\", %s, \"%s\", %s)%s", match[1], msg, match[3], match[6], match[4], match[7], match[5], match[8], match[9])
}

func MatchErrorfWithNamedParams(line string) string {
	return MatchSequentially([]Matcher{
		matchErrorfWithThreeNamedParams,
		matchErrorfWithTwoNamedParams,
	}, line)
}
