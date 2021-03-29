package repository

import (
	"github.com/angelmotta/covidAlerts-PE/source/model"
)

type NewCasesRepo interface {
	Create(report *model.CasesReport) error
	// GetByDate
}

type DeceasedCasesRepo interface {
	Create(report *model.DeceasedReport) error
	// GetByDate
}