package cms

import (
	"context"
	"time"
)

var _ PageStorage = (*EventPageStorage)(nil)

type partialPageStorage interface {
	FindByURL(ctx context.Context, siteID int64, url string, date time.Time) (Page, error)
	FindByRouteName(ctx context.Context, siteID int64, routeName string, date time.Time) (Page, error)
	FindByPageAlias(ctx context.Context, siteID int64, pageAlias string, date time.Time) (Page, error)
	FindByParentID(ctx context.Context, parentID int64, date time.Time) ([]Page, error)
}

type PageStorage interface {
	Storage[Page]
	partialPageStorage
}

type EventPageStorage struct {
	partialPageStorage
	*EventStorage[Page]
	BeforeFindByURL       Hook[*FindA3Event[int64, string, time.Time, Page]]
	AfterFindByURL        Hook[*FindA3Event[int64, string, time.Time, Page]]
	BeforeFindByRouteName Hook[*FindA3Event[int64, string, time.Time, Page]]
	AfterFindByRouteName  Hook[*FindA3Event[int64, string, time.Time, Page]]
	BeforeFindByPageAlias Hook[*FindA3Event[int64, string, time.Time, Page]]
	AfterFindByPageAlias  Hook[*FindA3Event[int64, string, time.Time, Page]]
	BeforeFindByParentID  Hook[*FindA2Event[int64, time.Time, []Page]]
	AfterFindByParentID   Hook[*FindA2Event[int64, time.Time, []Page]]
}

func NewEventPageStorage(storage PageStorage) *EventPageStorage {
	return &EventPageStorage{
		partialPageStorage: storage,
		EventStorage: NewEventStorage[Page](storage, func(p *Page) []string {
			return buildTags("page", p.ID, p.Name, p.RouteName)
		}, func(ids ...int64) []string {
			return buildIDs("page", ids...)
		}),
		BeforeFindByURL:       StubHook[*FindA3Event[int64, string, time.Time, Page]]{},
		AfterFindByURL:        StubHook[*FindA3Event[int64, string, time.Time, Page]]{},
		BeforeFindByRouteName: StubHook[*FindA3Event[int64, string, time.Time, Page]]{},
		AfterFindByRouteName:  StubHook[*FindA3Event[int64, string, time.Time, Page]]{},
		BeforeFindByPageAlias: StubHook[*FindA3Event[int64, string, time.Time, Page]]{},
		AfterFindByPageAlias:  StubHook[*FindA3Event[int64, string, time.Time, Page]]{},
		BeforeFindByParentID:  StubHook[*FindA2Event[int64, time.Time, []Page]]{},
		AfterFindByParentID:   StubHook[*FindA2Event[int64, time.Time, []Page]]{},
	}
}

func (s *EventPageStorage) FindByURL(ctx context.Context, siteID int64, url string, date time.Time) (Page, error) {
	event := &FindA3Event[int64, string, time.Time, Page]{Arg1: siteID, Arg2: url, Arg3: date}

	err := s.BeforeFindByURL.Trigger(ctx, event, func(ctx context.Context, event *FindA3Event[int64, string, time.Time, Page]) (err error) {
		if event.Result, err = s.partialPageStorage.FindByURL(ctx, event.Arg1, event.Arg2, event.Arg3); err != nil {
			return err
		}
		return s.AfterFindByURL.Trigger(ctx, event)
	})

	return event.Result, err
}

func (s *EventPageStorage) FindByRouteName(ctx context.Context, siteID int64, routeName string, date time.Time) (Page, error) {
	event := &FindA3Event[int64, string, time.Time, Page]{Arg1: siteID, Arg2: routeName, Arg3: date}

	err := s.BeforeFindByRouteName.Trigger(ctx, event, func(ctx context.Context, event *FindA3Event[int64, string, time.Time, Page]) (err error) {
		if event.Result, err = s.partialPageStorage.FindByRouteName(ctx, event.Arg1, event.Arg2, event.Arg3); err != nil {
			return err
		}
		return s.AfterFindByRouteName.Trigger(ctx, event)
	})

	return event.Result, err
}

func (s *EventPageStorage) FindByPageAlias(ctx context.Context, siteID int64, pageAlias string, date time.Time) (Page, error) {
	event := &FindA3Event[int64, string, time.Time, Page]{Arg1: siteID, Arg2: pageAlias, Arg3: date}

	err := s.BeforeFindByPageAlias.Trigger(ctx, event, func(ctx context.Context, event *FindA3Event[int64, string, time.Time, Page]) (err error) {
		if event.Result, err = s.partialPageStorage.FindByPageAlias(ctx, event.Arg1, event.Arg2, event.Arg3); err != nil {
			return err
		}
		return s.AfterFindByPageAlias.Trigger(ctx, event)
	})

	return event.Result, err
}

func (s *EventPageStorage) FindByParentID(ctx context.Context, parentID int64, date time.Time) ([]Page, error) {
	event := &FindA2Event[int64, time.Time, []Page]{Arg1: parentID, Arg2: date}

	err := s.BeforeFindByParentID.Trigger(ctx, event, func(ctx context.Context, event *FindA2Event[int64, time.Time, []Page]) (err error) {
		if event.Result, err = s.partialPageStorage.FindByParentID(ctx, event.Arg1, event.Arg2); err != nil {
			return err
		}
		return s.AfterFindByParentID.Trigger(ctx, event)
	})

	return event.Result, err
}

func PageBeforeSave(storage PageStorage) func(context.Context, *SaveEvent[*Page]) error {
	return func(ctx context.Context, event *SaveEvent[*Page]) error {
		if event.Model.IsHybrid() {
			return nil
		}

		if event.Model.Parent == nil && event.Model.ParentID != nil {
			parent, err := storage.FindByID(ctx, *event.Model.ParentID)
			if err != nil {
				return err
			}
			event.Model.Parent = &parent
		}

		event.Model.FixURL()

		return nil
	}
}

func PageAfterSave(storage PageStorage) func(context.Context, *SaveEvent[*Page]) error {
	return func(ctx context.Context, event *SaveEvent[*Page]) error {
		if event.Model.IsHybrid() || event.Model.ParentID == nil && event.Model.URL == "/" {
			return nil
		}

		if children, err := storage.FindByParentID(ctx, event.Model.ID, time.Time{}); err != nil {
			for _, child := range children {
				p := &child
				p.Parent = event.Model

				_ = storage.Update(ctx, p)
			}
		}

		return nil
	}
}
