package usecase

import "kakeibodb/internal/model"

type EventPresenter interface {
	Present(events []*model.Event)
}

type AnalysisPresenter interface {
	Present(report *TimeSeriesReport)
}
