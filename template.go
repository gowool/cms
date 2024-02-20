package cms

import (
	"time"

	"github.com/gowool/cms/internal"
)

type TemplateType int

const (
	TemplateDB TemplateType = iota
	TemplateFS
)

type Template struct {
	ID        int64        `json:"id,omitempty"`
	Code      string       `json:"code,omitempty"`
	Content   string       `json:"content,omitempty"`
	Type      TemplateType `json:"type,omitempty"`
	Created   time.Time    `json:"created,omitempty"`
	Updated   time.Time    `json:"updated,omitempty"`
	Published *time.Time   `json:"published,omitempty"`
	Expired   *time.Time   `json:"expired,omitempty"`
}

func (t Template) String() string {
	if t.Code == "" {
		return "n/a"
	}
	return t.Code
}

func (t Template) Enabled(now time.Time) bool {
	return t.Published != nil &&
		!t.Published.IsZero() &&
		(t.Published.Before(now) || t.Published.Equal(now)) &&
		(t.Expired == nil || t.Expired.IsZero() || t.Expired.After(now))
}

func (t Template) ContentBytes() []byte {
	return internal.Bytes(t.Content)
}
