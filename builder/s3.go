package builder

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/kennykarnama/convert-live-to-vod/entity"
	"github.com/kennykarnama/convert-live-to-vod/errtype"
	"github.com/kennykarnama/convert-live-to-vod/s3"
	"io"
	"log"
	"os"
)

func ManifestPathsFromCsv(csvin string) ([]*entity.S3Storage, error) {
	var result []*entity.S3Storage
	f, err := os.Open(csvin)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	//read headers
	header, err := csvReader.Read()
	if err != nil {
		return nil, err
	}
	log.Printf("action=ManifestPathsFromCsv header=%v", header)
	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		if len(record) == 0 {
			return nil, fmt.Errorf("action=ManifestPathsFromCsv err=%v", errtype.ErrEmptyCsvRecord)
		}
		s3storage, err := entity.NewS3StorageFromUrl(record[0])
		if err != nil {
			return nil, err
		}
		result = append(result, s3storage)
	}
	return result, nil
}

func NewManifestSegmentsFromS3(ctx context.Context, bucket, prefix string, pageSize int64) ([]*entity.ManifestSegment, error) {
	segmentCollection, err := s3.ListObjectPagesWithContext(ctx, bucket, prefix, pageSize)
	if err != nil {
		return nil, err
	}
	if len(segmentCollection) == 0 {
		return nil, nil
	}
	var manifestSegments []*entity.ManifestSegment
	for _, segment := range segmentCollection {
		manifestSegments = append(manifestSegments, &entity.ManifestSegment{
			S3Storage: entity.S3Storage{
				Key:    *segment.Key,
				Bucket: bucket,
			},
			ModifiedTime: *segment.LastModified,
		})
	}
	return manifestSegments, nil
}
