package entity

import (
	"fmt"
	"github.com/beevik/etree"
	"github.com/kennykarnama/convert-live-to-vod/errtype"
	"io"
)

type ManifestType int

const (
	Dash ManifestType = iota
	Hls
)

type SegmentType int

const (
	Video SegmentType = iota
	Audio
)

type Manifest struct {
	Content      *etree.Document
	ManifestType ManifestType
}

func NewManifestFromReader(source io.Reader, mt ManifestType) (*Manifest, error) {
	doc := etree.NewDocument()
	n, err := doc.ReadFrom(source)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, fmt.Errorf("action=NewManifestFromReader err=%v", errtype.ErrEmptyContent)
	}
	return &Manifest{
		Content:      doc,
		ManifestType: mt,
	}, nil
}

func (m *Manifest) GetType() ManifestType {
	return m.ManifestType
}

func (m *Manifest) GetContent() *etree.Document {
	return m.Content
}

func (m *Manifest) Save(w io.Writer) (n int64, err error) {
	return m.Content.WriteTo(w)
}

func (m *Manifest) GetSegmentPaths() (map[SegmentType]string, error) {
	//todo: fetch segments from manifest xml
	segments := map[SegmentType]string{}
	segments[Video] = "stream-0"
	segments[Audio] = "stream-1"
	return segments, nil
}
