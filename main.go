package main

import (
	"context"
	"fmt"
	"github.com/kennykarnama/convert-live-to-vod/builder"
	"github.com/kennykarnama/convert-live-to-vod/config"
	"github.com/kennykarnama/convert-live-to-vod/entity"
	"github.com/kennykarnama/convert-live-to-vod/errtype"
	"github.com/kennykarnama/convert-live-to-vod/function"
	"github.com/kennykarnama/convert-live-to-vod/s3"
	"log"
	"os"
	"time"
)

var (
	ctx context.Context
)

//Testing on several videos
//For VOD, it is more accurate to diff between
//time of the last and first segment.
func main() {
	ctx = context.Background()
	env := config.Get()
	infile := env.InputFile
	var dashCollection []*entity.Manifest

	s3collection, err := builder.ManifestPathsFromCsv(infile)
	errtype.OnErrorPanic(err)

	for _, s3item := range s3collection {
		fmt.Println(s3item.Bucket, s3item.Key)
		content, err := s3.GetWithContext(ctx, s3item.Bucket, s3item.Key)
		errtype.OnErrorPanic(err)
		dash, err := entity.NewManifestFromReader(content, entity.Dash)
		errtype.OnErrorPanic(err)
		dashCollection = append(dashCollection, dash)
	}

	for _, dashItem := range dashCollection {
		_, err := function.LiveToVod(dashItem)
		errtype.OnErrorPanic(err)
		dashItem.Save(os.Stdout)
	}

	// test list objects with pagination
	segmentCollection, err := builder.NewManifestSegmentsFromS3(ctx, env.S3Bucket, "live-quiz/7a73dfbb-9049-4680-8af6-3f11850e5af5/7a73dfbb-9049-4680-8af6-3f11850e5af5/stream-1", 20)
	errtype.OnErrorPanic(err)
	videoDuration, err := function.CalculateVideoDurationFromSegmentsModifiedTime(segmentCollection)
	errtype.OnErrorPanic(err)
	log.Printf("video-duration=%v", videoDuration)
	targetDuration, err := time.ParseDuration("0h7m27s")
	errtype.OnErrorPanic(err)
	log.Printf("targetDuration=%v", targetDuration.String())
	seekSegmentNum := function.GetStartSegment(segmentCollection, targetDuration)
	log.Printf("seekSegmentNum=%v", seekSegmentNum)
}
