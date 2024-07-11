package reporter

import "cotton/internal/result"

type NoReporter struct{}

func NewNoReporter() Reporter {
	return &NoReporter{}
}

func (r *NoReporter) Save(testResult *result.TestsuiteResult) error {
	return nil
}
