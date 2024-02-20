package cms

import (
	"fmt"
	"time"
)

type Site struct {
	ID           int64          `json:"id,omitempty"`
	Name         string         `json:"name,omitempty"`
	Title        string         `json:"title,omitempty"`
	Separator    string         `json:"separator,omitempty"`
	Host         string         `json:"host,omitempty"`
	Locale       string         `json:"locale,omitempty"`
	RelativePath string         `json:"relative_path,omitempty"`
	IsDefault    bool           `json:"is_default,omitempty"`
	Metas        []Meta         `json:"metas,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	Created      time.Time      `json:"created,omitempty"`
	Updated      time.Time      `json:"updated,omitempty"`
	Published    *time.Time     `json:"published,omitempty"`
	Expired      *time.Time     `json:"expired,omitempty"`
}

func (s Site) String() string {
	if s.Name == "" {
		return "n/a"
	}
	return s.Name
}

func (s Site) IsEnabled(now time.Time) bool {
	return s.Published != nil &&
		!s.Published.IsZero() &&
		(s.Published.Before(now) || s.Published.Equal(now)) &&
		(s.Expired == nil || s.Expired.IsZero() || s.Expired.After(now))
}

func (s Site) IsLocalhost() bool {
	return s.Host == "localhost"
}

func (s Site) URL() string {
	if s.IsLocalhost() {
		return s.RelativePath
	}

	return fmt.Sprintf("//%s%s", s.Host, s.RelativePath)
}
