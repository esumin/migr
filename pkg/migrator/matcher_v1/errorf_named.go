package matcher_v1

import (
	"fmt"
	"regexp"
	"strings"

	"mig/pkg/util"
)

// example
// return nil, errors.Errorf("Not implemented")
func matchErrorfWithoutParams(line string) string {
	regex := regexp.MustCompile(`(.*)errors.Errorf\("([^%]*)"\)(.*)`)
	match := regex.FindStringSubmatch(line)
	if match == nil {
		return ""
	}

	return fmt.Sprintf("%serrkit.New(\"%s\")%s", match[1], match[2], match[3])
}

// example
// return nil, errors.Errorf("Found an unexpected number of volumes: volume_id=%s result_count=%d", id, len(vols))
// return nil, errors.Errorf("Required volume fields not available, volumeType: %s, Az: %s", snapshot.Volume.VolumeType, snapshot.Volume.Az)
func matchErrorfWithTwoNamedParams(line string) string {
	res := matchErrorfWithTwoNamedParams1(line)
	if res != "" {
		return res
	}

	return matchErrorfWithTwoNamedParams2(line)
}

// example
// return nil, errors.Errorf("Found an unexpected number of volumes: volume_id=%s result_count=%d", id, len(vols))
func matchErrorfWithTwoNamedParams1(line string) string {
	regex := regexp.MustCompile(`(.*)errors.Errorf\("(.+?)(?:[,:.]\s+|\s+)(\w+)(?:=|:)\s*%[sd]{1}[,\s]+(\w+)(?:=|:)\s*%[sd]{1}\s*",\s*([^,]+),\s*([^,]+)\)(.*)`)
	match := regex.FindStringSubmatch(line)
	if match == nil || strings.Contains(match[4], "%") || strings.Contains(match[4], "fmt.Sprintf") {
		return ""
	}

	return fmt.Sprintf("%serrkit.New(\"%s\", \"%s\", %s, \"%s\", %s)%s", match[1], match[2], match[3], match[5], match[4], match[6], match[7])
}

// example
// return nil, errors.Errorf("Found an unexpected number of volumes: volume_id=%s result_count=%d", id, len(vols))
// return nil, errors.Errorf("Required volume fields not available, volumeType: %s, Az: %s", snapshot.Volume.VolumeType, snapshot.Volume.Az)
func matchErrorfWithTwoNamedParams2(line string) string {
	regex := regexp.MustCompile(`(.*)errors.Errorf\("(.+?)(?:[,:.]\s+|\s+)(\w+)(?:=|:)\s*%[sd]{1}[,\s]+(\w+)(?:=|:)\s*%[sd]{1}\s*",\s*([^,]+),\s*([^,]+)\)(.*)`)
	match := regex.FindStringSubmatch(line)
	if match == nil || strings.Contains(match[4], "%") || strings.Contains(match[4], "fmt.Sprintf") {
		return ""
	}

	return fmt.Sprintf("%serrkit.New(\"%s\", \"%s\", %s, \"%s\", %s)%s", match[1], match[2], match[3], match[5], match[4], match[6], match[7])
}

// example
// return nil, errors.Errorf("Required volume fields not available, volumeType: %s, Az: %s, VolumeTags: %v", snapshot.Volume.VolumeType, snapshot.Volume.Az, snapshot.Volume.Tags)",
func matchErrorfWithThreeNamedParams(line string) string {
	regex := regexp.MustCompile(`(.*)errors.Errorf\("(.+)\s+([^:=]+)[:=]{1}\s*%[sdv]{1}[,\s]+([^:=]+)[:=]{1}\s*%[sdv]{1}[,\s]+([^:=]+)[:=]{1}\s*%[sdv]{1}\s*",\s*([^,]+),\s*([^,]+),\s*([^,]+)\)(.*)`)
	match := regex.FindStringSubmatch(line)
	if match == nil || strings.Contains(match[4], "%") {
		return ""
	}

	msg := maybeTrimComma(match[2])

	return fmt.Sprintf("%serrkit.New(\"%s\", \"%s\", %s, \"%s\", %s, \"%s\", %s)%s", match[1], msg, match[3], match[6], match[4], match[7], match[5], match[8], match[9])
}

var a = 0

// example
// return "", errors.Errorf("no zones specified, zone: %s", az)
func matchErrorfWithImplicitParam(line string) string {
	//fmt.Println("11111")
	//a++
	//if a == 274 {
	//	fmt.Println("22222")
	//}
	wrap := regexp.MustCompile(`(.*)errors.Errorf\("([^%]*?)(:?\s\%[a-z]{1})([^%]*)",\s*([a-zA-Z0-9.]+)\)(.*)`)
	match := wrap.FindStringSubmatch(line)
	if match == nil {
		return ""
	}

	paramName := util.GuessParamName(match[2], match[5])

	return fmt.Sprintf("%serrkit.New(\"%s%s\", \"%s\", %s)%s", match[1], match[2], match[4], paramName, match[5], match[6])
}

// example
// return nil, errors.Errorf("Required volume fields not available, volumeType: %s, Az: %s", snapshot.Volume.VolumeType, snapshot.Volume.Az)

func MatchErrorfWithNamedParams(line string) string {
	return MatchSequentially([]Matcher{
		matchErrorfWithThreeNamedParams,
		matchErrorfWithTwoNamedParams,
		matchErrorfWithImplicitParam,
		matchErrorfWithoutParams,
	}, line)
}
