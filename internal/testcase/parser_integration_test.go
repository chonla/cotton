//go:build integration
// +build integration

package testcase_test

import (
	"cotton/internal/assertion"
	"cotton/internal/capture"
	"cotton/internal/config"
	"cotton/internal/executable"
	"cotton/internal/httphelper"
	"cotton/internal/logger"
	"cotton/internal/reader"
	"cotton/internal/testcase"
	"cotton/internal/variable"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsingCompleteTestcaseMarkdownFile(t *testing.T) {
	curdir, _ := os.Getwd()
	config := &config.Config{
		BaseDir: curdir + "/../../etc/examples/fakestoreapi.com/testcases/",
	}

	executableParserOptions := &executable.ParserOptions{
		Configurator:  config,
		FileReader:    reader.New(os.ReadFile),
		RequestParser: &httphelper.HTTPRequestParser{},
		Logger:        logger.NewNilLogger(logger.Compact),
	}
	executableParser := executable.NewParser(executableParserOptions)

	parserOptions := &testcase.ParserOptions{
		Configurator:     config,
		FileReader:       reader.New(os.ReadFile),
		RequestParser:    &httphelper.HTTPRequestParser{},
		Logger:           logger.NewNilLogger(logger.Compact),
		ExecutableParser: executableParser,
	}
	parser := testcase.NewParser(parserOptions)

	testcaseOptions := &testcase.TestcaseOptions{
		RequestParser: &httphelper.HTTPRequestParser{},
		Logger:        logger.NewNilLogger(logger.Compact),
	}
	result, err := parser.FromMarkdownFile("../../etc/examples/fakestoreapi.com/testcases/full_documented.md")

	executableOptions := &executable.ExecutableOptions{
		RequestParser: &httphelper.HTTPRequestParser{},
		Logger:        logger.NewNilLogger(logger.Compact),
	}
	expectedSetup1 := executable.New("Create a new user for authentication", `POST https://fakestoreapi.com/users HTTP/1.1
Content-Type: application/json
Content-Length: 277

{"email":"John@gmail.com","username":"{{username}}","password":"{{password}}","name":{"firstname":"John","lastname":"Doe"},"address":{"city":"kilcoole","street":"7835 new road","number":3,"zipcode":"12926-3874","geolocation":{"lat":"-37.3159","long":"81.1496"}},"phone":"1-570-236-7033"}`, executableOptions)
	expectedSetup1.AddCapture(capture.New("new_user_id", "Body.id"))
	expectedSetup1.AddVariable(&variable.Variable{
		Name:  "username",
		Value: "mor_2314",
	})
	expectedSetup1.AddVariable(&variable.Variable{
		Name:  "password",
		Value: "83r5^_",
	})

	expectedSetup2 := executable.New("Authentication with the new user", `POST https://fakestoreapi.com/auth/login HTTP/1.1
Content-Type: application/json
Content-Length: 43

{"username":"{{username}}","password":"{{password}}"}`, executableOptions)
	expectedSetup2.AddCapture(capture.New("access_token", "Body.token"))
	expectedSetup2.AddVariable(&variable.Variable{
		Name:  "username",
		Value: "mor_2314",
	})
	expectedSetup2.AddVariable(&variable.Variable{
		Name:  "password",
		Value: "83r5^_",
	})

	expectedTeardown := executable.New("Delete test user", `DELETE https://fakestoreapi.com/users/{{new_user_id}} HTTP/1.1
Content-Type: application/json`, executableOptions)
	expectedTestcase := testcase.NewTestcase("Full documented testcase", "The testcase is described by providing paragraphs right after the test case title.", `GET https://fakestoreapi.com/products HTTP/1.1
Authorization: Bearer {{access_token}}`, testcaseOptions)
	expectedTestcase.AddSetup(expectedSetup1)
	expectedTestcase.AddSetup(expectedSetup2)
	expectedTestcase.AddTeardown(expectedTeardown)

	eqOp, _ := assertion.NewOp("==")
	gtOp, _ := assertion.NewOp(">")
	expectedTestcase.AddAssertion(assertion.New("StatusCode", eqOp, float64(200)))
	expectedTestcase.AddAssertion(assertion.New("Body.#", gtOp, float64(0)))
	expectedTestcase.AddAssertion(assertion.New("Body.0.id", eqOp, float64(1)))

	assert.NoError(t, err)
	assert.Equal(t, expectedTestcase, result)
}
