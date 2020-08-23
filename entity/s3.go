package entity

import (
	"fmt"
	"github.com/kennykarnama/convert-live-to-vod/errtype"
	"net/url"
	"strings"
)

type S3Storage struct {
	Bucket string
	Key    string
}

func NewS3StorageFromUrl(rawURL string) (*S3Storage, error) {
	unescaped, err := url.QueryUnescape(rawURL)

	if err != nil {
		return nil, err
	}
	q, err := url.Parse(unescaped)
	if err != nil {
		return nil, err
	}
	if q.Scheme == "s3" {
		return &S3Storage{
			Bucket: q.Host,
			Key:    strings.Replace(q.Path, "/", "", 1),
		}, nil
	}
	if q.Scheme == "https" || q.Scheme == "http" {
		s := strings.Split(q.Host, ".")
		if len(s) < 3 {
			return nil, errtype.ErrNotS3
		}
		if s[0] != "s3-ap-southeast-1" && s[1] != "s3-ap-southeast-1" {
			return nil, errtype.ErrNotS3
		}

		paths := strings.Replace(q.Path, "/", "", 1)

		s3Data := strings.SplitN(paths, "/", 2)

		if len(paths) == 0 || len(s3Data) == 0 {
			return nil, errtype.ErrNotS3
		}

		if len(s) == 4 {
			return &S3Storage{
				Bucket: s[0],
				Key:    strings.Join(s3Data[0:], "/"),
			}, nil
		} else if len(s) == 3 {
			return &S3Storage{
				Bucket: s3Data[0],
				Key:    strings.Join(s3Data[1:], "/"),
			}, nil
		}
		return nil, errtype.ErrNotS3

	}
	return nil, errtype.ErrNotValidScheme
}

func (s *S3Storage) Url() string {
	return fmt.Sprintf("%s://%s.s3-ap-southeast-1.amazonaws.com/%s", "https", s.Bucket, s.Key)
}
