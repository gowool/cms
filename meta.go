package cms

type MetaType string

const (
	MetaName     MetaType = "name"
	MetaEquiv    MetaType = "http-equiv"
	MetaProperty MetaType = "property"
)

func (t MetaType) String() string {
	return string(t)
}

type Meta struct {
	Type    MetaType `json:"type,omitempty"`
	Key     string   `json:"key,omitempty"`
	Content string   `json:"content,omitempty"`
}

func (m Meta) Equals(another Meta) bool {
	return m.Type == another.Type && m.Key == another.Key
}
