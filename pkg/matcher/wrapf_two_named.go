package matcher

import (
	"fmt"
	"regexp"
	"strings"

	"mig/pkg/util"
)

// example
// return nil, errors.Wrapf(err, "Failed to create snapshot, volume_id: %s", *csi.VolumeId)
func matchWrapfWithNamedParam(line string) string {
	wrap := regexp.MustCompile(`(.*)errors.Wrapf\((.*?), "(.*?)\s*([^.\s=]+)[:=]{1}\s*\%s\s*",\s*([^,)]+)\)(.*)`)
	match := wrap.FindStringSubmatch(line)
	if match == nil || strings.Contains(match[3], "%") {
		return ""
	}

	msg := maybeTrimComma(match[3])

	return fmt.Sprintf("%serrkit.Wrap(%s, \"%s\", \"%s\", %s)%s", match[1], match[2], msg, match[4], match[5], match[6])
}

// example
// return errors.Wrapf(err, "Failed to delete subnet group. You may need to delete it manually. app=%s name=%s", a.name, a.dbSubnetGroup)
// return errors.Wrapf(err, "DiskCLient.CreateOrUpdate in VolumeCreateFromSnapshot, diskName: %s, snapshotID: %s", diskName, snapshot.ID)
func matchWrapfWithTwoNamedParams(line string) string {
	wrap := regexp.MustCompile(`(.*)errors.Wrapf\(((.*?), "(.*?)\s*([^.\s=]+)=\%s\s*([^.\s=]+)=\%s\s*",\s*([^,)]+),\s*([^,)]+))\)(.*)`)
	match := wrap.FindStringSubmatch(line)
	if match == nil || strings.Contains(match[4], "%") || strings.Contains(match[4], "fmt.Sprintf") {
		return ""
	}

	return fmt.Sprintf("%serrkit.Wrap(%s, \"%s\", \"%s\", %s, \"%s\", %s)%s", match[1], match[3], match[4], match[5], match[7], match[6], match[8], match[9])
}

// example
// return errors.Wrapf(err, "Failed to install helm chart. app=%s chart=%s release=%s", c.name, c.chart.Chart, c.chart.Release)
func matchWrapfWithThreeNamedParams(line string) string {
	wrap := regexp.MustCompile(`(.*)errors.Wrapf\(((.*?), "(.*?)\s*([^.\s=]+)=\%s\s*([^.\s=]+)=\%s\s*([^.\s=]+)=\%s\s*",\s*([^,)]+),\s*([^,)]+),\s*([^,)]+))\)(.*)`)
	match := wrap.FindStringSubmatch(line)
	if match == nil || strings.Contains(match[4], "%") || strings.Contains(match[4], "fmt.Sprintf") {
		return ""
	}

	return fmt.Sprintf("%serrkit.Wrap(%s, \"%s\", \"%s\", %s, \"%s\", %s, \"%s\", %s)%s", match[1], match[3], match[4], match[5], match[8], match[6], match[9], match[7], match[10], match[11])
}

// example
// return errors.Wrapf(err, "Waiting on snapshot %v", snap)
// return nil, errors.Wrapf(err, "Snapshot %s did not complete", snapID)
func matchWrapfWithImplicitParam(line string) string {
	wrap := regexp.MustCompile(`(.*)errors.Wrapf\((.*),\s*"([^%]*?)(:?\s\%[a-z]{1})([^%]*)",\s*([a-zA-Z0-9.]+)\)(.*)`)
	match := wrap.FindStringSubmatch(line)
	if match == nil {
		return ""
	}

	paramName := util.GuessParamName(match[3], match[6])

	return fmt.Sprintf("%serrkit.Wrap(%s, \"%s%s\", \"%s\", %s)%s", match[1], match[2], match[3], match[5], paramName, match[6], match[7])
}

func MatchWrapfWithNamedParams(line string) string {
	return MatchSequentially([]Matcher{
		matchWrapfWithThreeNamedParams,
		matchWrapfWithTwoNamedParams,
		matchWrapfWithNamedParam,
		matchWrapfWithImplicitParam,
	}, line)
}
