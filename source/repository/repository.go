package repository

import (
	"github.com/angelmotta/covidAlerts-PE/source/model"
)

type Repository interface {
	Create(report *model.CasesReport) error
	// GetByDate
}
