package repository

import (
	"github.com/angelmotta/covidAlerts-PE/source/model"
)

type NewCasesRepo interface {
	Create(report *model.CasesReport) (int, error)
	// GetByDate
}

type DeceasedCasesRepo interface {
	Create(report *model.DeceasedReport) (int, error)
	// GetByDate
}