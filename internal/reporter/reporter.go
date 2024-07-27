package reporter

import "cotton/internal/result"

type ReporterType string

const (
	CTRF ReporterType = "ctrf"
	HTML ReporterType = "html"
)

type Reporter interface {
	Save(testResult *result.TestsuiteResult) error
}

func NewReporter(reporterType ReporterType) Reporter {
	if reporterType == CTRF {
		return NewCTRFReporter()
	}
	if reporterType == HTML {
		return NewHTMLReporter()
	}
	return NewNoReporter()
}
