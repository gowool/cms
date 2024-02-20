package cms

import (
	"fmt"

	"github.com/gowool/cr"
)

type SaveEvent[T any] struct {
	Model T
	tags  []string
}

func NewSaveEvent[T any](model T, tags ...string) *SaveEvent[T] {
	return &SaveEvent[T]{
		Model: model,
		tags:  tags,
	}
}

func (e *SaveEvent[T]) Tags() []string {
	return e.tags
}

type DeleteEvent struct {
	IDs  []int64
	tags []string
}

func NewDeleteEvent(ids []int64, tags ...string) *DeleteEvent {
	return &DeleteEvent{
		IDs:  ids,
		tags: tags,
	}
}

func (e *DeleteEvent) Tags() []string {
	return e.tags
}

type FindByIDEvent[T any] struct {
	ID    int64
	Model T
}

type FindEvent[T any] struct {
	Criteria *cr.Criteria
	Data     []T
}

type FindAndCountEvent[T any] struct {
	Criteria *cr.Criteria
	Data     []T
	Count    int
}

type FindA0Event[R any] struct {
	Result R
}

type FindA1Event[A1, R any] struct {
	Arg1   A1
	Result R
}

type FindA2Event[A1, A2, R any] struct {
	Arg1   A1
	Arg2   A2
	Result R
}

type FindA3Event[A1, A2, A3, R any] struct {
	Arg1   A1
	Arg2   A2
	Arg3   A3
	Result R
}

func BuildTags(name string, id int64, data ...string) []string {
	tags := make([]string, 0, len(data)+2)
	tags = append(tags, name)

	if id != 0 {
		tags = append(tags, fmt.Sprintf("%s:%d", name, id))
	}

	for _, item := range data {
		if item != "" {
			tags = append(tags, fmt.Sprintf("%s:%s", name, item))
		}
	}

	return tags
}

func BuildTagsByIDs(name string, data ...int64) []string {
	tags := make([]string, 0, len(data)+1)
	tags = append(tags, name)

	for _, item := range data {
		if item > 0 {
			tags = append(tags, fmt.Sprintf("%s:%d", name, item))
		}
	}

	return tags
}
