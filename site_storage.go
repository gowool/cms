package cms

import (
	"context"
	"time"
)

type partialSiteStorage interface {
	FindByHosts(ctx context.Context, hosts []string, date time.Time) ([]Site, error)
	FindAll(ctx context.Context) ([]Site, error)
}

type SiteStorage interface {
	Storage[Site]
	partialSiteStorage
}

type EventSiteStorage struct {
	partialSiteStorage
	*EventStorage[Site]
	BeforeFindByHosts Hook[*FindA2Event[[]string, time.Time, []Site]]
	AfterFindByHosts  Hook[*FindA2Event[[]string, time.Time, []Site]]
	BeforeFindAll     Hook[*FindA0Event[[]Site]]
	AfterFindAll      Hook[*FindA0Event[[]Site]]
}

func NewEventSiteStorage(storage SiteStorage) *EventSiteStorage {
	return &EventSiteStorage{
		partialSiteStorage: storage,
		EventStorage: NewEventStorage[Site](storage, func(s *Site) []string {
			return buildTags("site", s.ID, s.Name)
		}, func(ids ...int64) []string {
			return buildIDs("site", ids...)
		}),
		BeforeFindByHosts: StubHook[*FindA2Event[[]string, time.Time, []Site]]{},
		AfterFindByHosts:  StubHook[*FindA2Event[[]string, time.Time, []Site]]{},
		BeforeFindAll:     StubHook[*FindA0Event[[]Site]]{},
		AfterFindAll:      StubHook[*FindA0Event[[]Site]]{},
	}
}

func (s *EventSiteStorage) FindByHosts(ctx context.Context, hosts []string, date time.Time) ([]Site, error) {
	event := &FindA2Event[[]string, time.Time, []Site]{Arg1: hosts, Arg2: date}

	err := s.BeforeFindByHosts.Trigger(ctx, event, func(ctx context.Context, event *FindA2Event[[]string, time.Time, []Site]) (err error) {
		if event.Result, err = s.partialSiteStorage.FindByHosts(ctx, event.Arg1, event.Arg2); err != nil {
			return err
		}
		return s.AfterFindByHosts.Trigger(ctx, event)
	})

	return event.Result, err
}

func (s *EventSiteStorage) FindAll(ctx context.Context) ([]Site, error) {
	event := &FindA0Event[[]Site]{}

	err := s.BeforeFindAll.Trigger(ctx, event, func(ctx context.Context, event *FindA0Event[[]Site]) (err error) {
		if event.Result, err = s.partialSiteStorage.FindAll(ctx); err != nil {
			return err
		}
		return s.AfterFindAll.Trigger(ctx, event)
	})

	return event.Result, err
}
