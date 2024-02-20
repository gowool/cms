package cms

import (
	"context"

	"github.com/gowool/cr"
)

var _ Storage[any] = (*EventStorage[any])(nil)

type Storage[T any] interface {
	Create(ctx context.Context, model *T) error
	Update(ctx context.Context, model *T) error
	Delete(ctx context.Context, ids ...int64) error
	FindByID(ctx context.Context, id int64) (T, error)
	Find(ctx context.Context, criteria *cr.Criteria) ([]T, error)
	FindAndCount(ctx context.Context, criteria *cr.Criteria) ([]T, int, error)
}

type EventStorage[T any] struct {
	Storage[T]

	SaveEventTags   func(*T) []string
	DeleteEventTags func(ids ...int64) []string

	BeforeCreate       Hook[*SaveEvent[*T]]
	AfterCreate        Hook[*SaveEvent[*T]]
	BeforeUpdate       Hook[*SaveEvent[*T]]
	AfterUpdate        Hook[*SaveEvent[*T]]
	BeforeDelete       Hook[*DeleteEvent]
	AfterDelete        Hook[*DeleteEvent]
	BeforeFindByID     Hook[*FindByIDEvent[T]]
	AfterFindByID      Hook[*FindByIDEvent[T]]
	BeforeFind         Hook[*FindEvent[T]]
	AfterFind          Hook[*FindEvent[T]]
	BeforeFindAndCount Hook[*FindAndCountEvent[T]]
	AfterFindAndCount  Hook[*FindAndCountEvent[T]]
}

func NewEventStorage[T any](storage Storage[T], saveEventTags func(*T) []string, deleteEventTags func(ids ...int64) []string) *EventStorage[T] {
	return &EventStorage[T]{
		Storage:            storage,
		SaveEventTags:      saveEventTags,
		DeleteEventTags:    deleteEventTags,
		BeforeCreate:       StubHook[*SaveEvent[*T]]{},
		AfterCreate:        StubHook[*SaveEvent[*T]]{},
		BeforeUpdate:       StubHook[*SaveEvent[*T]]{},
		AfterUpdate:        StubHook[*SaveEvent[*T]]{},
		BeforeDelete:       StubHook[*DeleteEvent]{},
		AfterDelete:        StubHook[*DeleteEvent]{},
		BeforeFindByID:     StubHook[*FindByIDEvent[T]]{},
		AfterFindByID:      StubHook[*FindByIDEvent[T]]{},
		BeforeFind:         StubHook[*FindEvent[T]]{},
		AfterFind:          StubHook[*FindEvent[T]]{},
		BeforeFindAndCount: StubHook[*FindAndCountEvent[T]]{},
		AfterFindAndCount:  StubHook[*FindAndCountEvent[T]]{},
	}
}

func (s *EventStorage[T]) Create(ctx context.Context, model *T) error {
	event := &SaveEvent[*T]{
		Model: model,
		tags:  s.SaveEventTags(model),
	}

	return s.BeforeCreate.Trigger(ctx, event, func(ctx context.Context, event *SaveEvent[*T]) error {
		if err := s.Storage.Create(ctx, event.Model); err != nil {
			return err
		}
		return s.AfterCreate.Trigger(ctx, event)
	})
}

func (s *EventStorage[T]) Update(ctx context.Context, model *T) error {
	event := &SaveEvent[*T]{
		Model: model,
		tags:  s.SaveEventTags(model),
	}

	return s.BeforeUpdate.Trigger(ctx, event, func(ctx context.Context, event *SaveEvent[*T]) error {
		if err := s.Storage.Update(ctx, event.Model); err != nil {
			return err
		}
		return s.AfterUpdate.Trigger(ctx, event)
	})
}

func (s *EventStorage[T]) Delete(ctx context.Context, ids ...int64) error {
	event := &DeleteEvent{IDs: ids, tags: s.DeleteEventTags(ids...)}

	return s.BeforeDelete.Trigger(ctx, event, func(ctx context.Context, event *DeleteEvent) error {
		if err := s.Storage.Delete(ctx, event.IDs...); err != nil {
			return err
		}
		return s.AfterDelete.Trigger(ctx, event)
	})
}

func (s *EventStorage[T]) FindByID(ctx context.Context, id int64) (T, error) {
	event := &FindByIDEvent[T]{ID: id}

	err := s.BeforeFindByID.Trigger(ctx, event, func(ctx context.Context, event *FindByIDEvent[T]) (err error) {
		if event.Model, err = s.Storage.FindByID(ctx, event.ID); err != nil {
			return err
		}
		return s.AfterFindByID.Trigger(ctx, event)
	})

	return event.Model, err
}

func (s *EventStorage[T]) Find(ctx context.Context, criteria *cr.Criteria) ([]T, error) {
	event := &FindEvent[T]{Criteria: criteria}

	err := s.BeforeFind.Trigger(ctx, event, func(ctx context.Context, event *FindEvent[T]) (err error) {
		if event.Data, err = s.Storage.Find(ctx, event.Criteria); err != nil {
			return err
		}
		return s.AfterFind.Trigger(ctx, event)
	})

	return event.Data, err
}

func (s *EventStorage[T]) FindAndCount(ctx context.Context, criteria *cr.Criteria) ([]T, int, error) {
	event := &FindAndCountEvent[T]{Criteria: criteria}

	err := s.BeforeFindAndCount.Trigger(ctx, event, func(ctx context.Context, event *FindAndCountEvent[T]) (err error) {
		if event.Data, event.Count, err = s.Storage.FindAndCount(ctx, event.Criteria); err != nil {
			return err
		}
		return s.AfterFindAndCount.Trigger(ctx, event)
	})

	return event.Data, event.Count, err
}
