package reporter

import (
	"cotton/internal/result"
	"html/template"
	"os"
	"time"
)

var HTMLLayout string = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Test Report</title>
<style type="text/css">
.summary-table,.result-table{border-collapse: collapse;}
.summary-table th,.summary-table td,.result-table th,.result-table td{padding: 4px 8px;border: 1px solid #9fb8c4;}
.summary-table th{background-color:rgb(152, 203, 245);text-align:left;}
.summary-table td{background-color: rgb(238, 247, 255);}
.result-table th{background-color:rgb(191, 202, 213);}
.result-table tr:nth-child(even){background-color:rgb(238, 245, 251);}
.result-table tr:nth-child(odd){background-color:#fff;}
.result-table td.test-result-passed,.result-table td.test-result-failed{text-transform:uppercase;}
.result-table td.test-result-passed{background-color:rgb(15, 197, 130);color:#fff;}
.result-table td.test-result-failed{background-color:rgb(197, 54, 15);color:#fff;}
</style>
</head>
<body>
<h1>Test Summary</h1>
<table class="summary-table">
<tr><th>Testcases executed by</th><td>{{.Tool.Name}} {{.Tool.Version}}</td></tr>
<tr><th>Testcase executed</th><td>{{.Summary.Tests}}</td></tr>
<tr><th>Passed</th><td>{{.Summary.Passed}}</td></tr>
<tr><th>Failed</th><td>{{.Summary.Failed}}</td></tr>
<tr><th>Pending</th><td>{{.Summary.Pending}}</td></tr>
<tr><th>Start - End</th><td>{{.Summary.Start}} - {{.Summary.Stop}}</td></tr>
<tr><th>Ellapsed Time</th><td>{{.Summary.EllapsedTime}}</td></tr>
</table>
<h1>Test Results</h1>
<table class="result-table">
<thead>
<tr>
<th>Title</th>
<th>Status</th>
<th>Message</th>
<th>Start - End</th>
<th>Ellapsed Time</th>
</tr>
</thead>
<tbody>
{{range $testResult := .Tests}}
<tr>
<td>{{$testResult.Name}}</td>
<td class="test-result-{{$testResult.Status}}">{{$testResult.Status}}</td>
<td>{{$testResult.ErrorMessage}}</td>
<td>{{$testResult.Start}} - {{$testResult.Stop}}</td>
<td>{{$testResult.EllapsedTime}}</td>
</tr>
{{end}}
</tbody>
</table>
</body>
</html>
`

type HTMLReport struct {
	Tool    HTMLTool
	Summary HTMLSummary
	Tests   []HTMLTest
}

type HTMLTool struct {
	Name    string
	Version string
}

type HTMLSummary struct {
	Tests        int
	Passed       int
	Failed       int
	Pending      int
	Skipped      int
	Other        int
	Start        string
	Stop         string
	EllapsedTime string
}

type HTMLTest struct {
	Name         string
	Status       string
	EllapsedTime string
	Start        string
	Stop         string
	ErrorMessage string
}

var HTMLTimeFormat string = "Jan 2, 2006 15:04:05"

type HTMLReporter struct{}

func NewHTMLReporter() Reporter {
	return &HTMLReporter{}
}

func (r *HTMLReporter) Save(testsuiteResult *result.TestsuiteResult) error {
	reportData := &HTMLReport{
		Tool: HTMLTool{
			Name:    "cotton",
			Version: testsuiteResult.AppVersion,
		},
		Summary: HTMLSummary{
			Tests:        testsuiteResult.TestCount,
			Passed:       testsuiteResult.PassedCount,
			Failed:       testsuiteResult.FailedCount,
			Skipped:      0,
			Pending:      testsuiteResult.SkippedCount,
			Other:        0,
			Start:        time.Unix(testsuiteResult.Start, 0).Format(HTMLTimeFormat),
			Stop:         time.Unix(testsuiteResult.Stop, 0).Format(HTMLTimeFormat),
			EllapsedTime: testsuiteResult.EllapsedTime.String(),
		},
	}
	for _, testResult := range testsuiteResult.TestResults {
		resultData := HTMLTest{
			Name:         testResult.Title,
			EllapsedTime: testResult.EllapsedTime.String(),
			Status:       "failed",
			ErrorMessage: "",
			Start:        time.Unix(testResult.Start, 0).Format(HTMLTimeFormat),
			Stop:         time.Unix(testResult.Stop, 0).Format(HTMLTimeFormat),
		}
		if testResult.Passed {
			resultData.Status = "passed"
		} else {
			if testResult.Error != nil {
				resultData.ErrorMessage = testResult.Error.Error()
			}
		}
		reportData.Tests = append(reportData.Tests, resultData)
	}

	htmlTemplate, err := template.New("report").Parse(HTMLLayout)
	if err != nil {
		return err
	}

	htmlFile, err := os.Create("./result.html")
	if err != nil {
		return err
	}
	defer htmlFile.Close()

	err = htmlTemplate.Execute(htmlFile, reportData)
	return err
}
