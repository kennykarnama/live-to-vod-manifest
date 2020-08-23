package function

import (
	"fmt"
	"github.com/kennykarnama/convert-live-to-vod/entity"
	"github.com/kennykarnama/convert-live-to-vod/errtype"
	"time"
)

//LiveToVod converts a manifest
//from live type to be a vod
func LiveToVod(m *entity.Manifest) (mutated bool, err error) {
	if m.GetType() == entity.Dash {
		// do
		c := m.GetContent()
		el := c.SelectElement("MPD")
		if el == nil {
			return false, nil
		}
		attr := el.SelectAttr("type")
		if attr == nil {
			return false, nil
		}
		if attr.Value == "dynamic" {
			attr.Value = "static"
		}
		el.CreateAttr("mediaPresentationDuration", fmt.Sprintf("PT13.8S"))
		return true, nil
	}
	return false, nil
}
//CalculateVideoDurationFromSegmentsModifiedTime
//returns video duration
//by subtracting last segment modified time and first segment
func CalculateVideoDurationFromSegmentsModifiedTime(segments []*entity.ManifestSegment) (time.Duration, error) {
	if len(segments) < 2 {
		return time.Duration(0), errtype.EmptySegments
	}
	firstSegment := segments[0]
	lastSegment := segments[len(segments)-1]
	d := lastSegment.ModifiedTime.Sub(firstSegment.ModifiedTime)
	return d, nil
}

func GetStartSegment(segments []*entity.ManifestSegment, targetDuration time.Duration) int {
	totalSegment := len(segments)
	for idx := 1; idx < totalSegment-1; idx++ {
		t := segments[idx].ModifiedTime.UnixNano() - segments[0].ModifiedTime.UnixNano()
		if t >= targetDuration.Nanoseconds() {
			return idx+1
		}
	}
	return -1
}


