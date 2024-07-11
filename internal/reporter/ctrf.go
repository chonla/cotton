package reporter

import (
	"cotton/internal/result"
	"encoding/json"
	"os"
	"runtime"
	"strings"
)

type CTRFReport struct {
	Results CTRFResults `json:"result"`
}

type CTRFResults struct {
	Tool        CTRFTool        `json:"tool"`
	Summary     CTRFSummary     `json:"summary"`
	Tests       []CTRFTest      `json:"tests"`
	Environment CTRFEnvironment `json:"environment"`
}

type CTRFTool struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type CTRFSummary struct {
	Tests   int   `json:"tests"`
	Passed  int   `json:"passed"`
	Failed  int   `json:"failed"`
	Pending int   `json:"pending"`
	Skipped int   `json:"skipped"`
	Other   int   `json:"other"`
	Start   int64 `json:"start"`
	Stop    int64 `json:"stop"`
}

type CTRFTest struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Duration int64  `json:"duration"`
}

type CTRFEnvironment struct {
	OSPlatform string `json:"osPlatform"`
}

type CTRFReporter struct{}

func NewCTRFReporter() Reporter {
	return &CTRFReporter{}
}

func (r *CTRFReporter) Save(testsuiteResult *result.TestsuiteResult) error {
	reportData := &CTRFReport{
		Results: CTRFResults{
			Tool: CTRFTool{
				Name:    "cotton",
				Version: testsuiteResult.AppVersion,
			},
			Summary: CTRFSummary{
				Tests:   testsuiteResult.TestCount,
				Passed:  testsuiteResult.PassedCount,
				Failed:  testsuiteResult.FailedCount,
				Skipped: 0,
				Pending: testsuiteResult.SkippedCount,
				Other:   0,
				Start:   testsuiteResult.Start,
				Stop:    testsuiteResult.Stop,
			},
			Tests: []CTRFTest{},
			Environment: CTRFEnvironment{
				OSPlatform: runtime.GOOS,
			},
		},
	}
	for _, testResult := range testsuiteResult.TestResults {
		resultData := CTRFTest{
			Name:     testResult.Title,
			Duration: testResult.EllapsedTime.Duration(),
			Status:   "failed",
		}
		if testResult.Passed {
			resultData.Status = "passed"
		}
		reportData.Results.Tests = append(reportData.Results.Tests, resultData)
	}
	jsonStr, err := json.MarshalIndent(reportData, "", strings.Repeat(" ", 4))
	if err != nil {
		return err
	}
	err = os.WriteFile("./result.json", []byte(jsonStr), 0644)
	return err
}
