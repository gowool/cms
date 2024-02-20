package cms

import (
	"context"
	"time"
)

type partialTemplateStorage interface {
	FindByCode(ctx context.Context, code string, date time.Time) (Template, error)
}

type TemplateStorage interface {
	Storage[Template]
	partialTemplateStorage
}

type EventTemplateStorage struct {
	partialTemplateStorage
	*EventStorage[Template]
	BeforeFindByCode Hook[*FindA2Event[string, time.Time, Template]]
	AfterFindByCode  Hook[*FindA2Event[string, time.Time, Template]]
}

func NewEventTemplateStorage(storage TemplateStorage) *EventTemplateStorage {
	return &EventTemplateStorage{
		partialTemplateStorage: storage,
		EventStorage: NewEventStorage[Template](storage, func(t *Template) []string {
			return buildTags("template", t.ID, t.Code)
		}, func(ids ...int64) []string {
			return buildIDs("template", ids...)
		}),
		BeforeFindByCode: StubHook[*FindA2Event[string, time.Time, Template]]{},
		AfterFindByCode:  StubHook[*FindA2Event[string, time.Time, Template]]{},
	}
}

func (s *EventTemplateStorage) FindByCode(ctx context.Context, code string, date time.Time) (Template, error) {
	event := &FindA2Event[string, time.Time, Template]{Arg1: code, Arg2: date}

	err := s.BeforeFindByCode.Trigger(ctx, event, func(ctx context.Context, event *FindA2Event[string, time.Time, Template]) (err error) {
		if event.Result, err = s.partialTemplateStorage.FindByCode(ctx, event.Arg1, event.Arg2); err != nil {
			return err
		}
		return s.AfterFindByCode.Trigger(ctx, event)
	})

	return event.Result, err
}
