package reporter

import "cotton/internal/result"

type ReporterType string

const (
	CTRF ReporterType = "ctrf"
)

type Reporter interface {
	Save(testResult *result.TestsuiteResult) error
}

func NewReporter(reporterType ReporterType) Reporter {
	if reporterType == CTRF {
		return NewCTRFReporter()
	}
	return NewNoReporter()
}
