package engine

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// Evaluator evaluates expectations against HTTP responses
type Evaluator struct{}

// NewEvaluator creates a new evaluator
func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

// EvaluateExpectations evaluates all expectations against a response
func (e *Evaluator) EvaluateExpectations(resp *http.Response, expectations []Expectation) ([]string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Verify it's valid JSON
	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("response is not valid JSON: %w", err)
	}

	failures := []string{}
	jsonBody := string(body)

	for _, exp := range expectations {
		result, err := e.evaluateExpectation(exp, jsonBody)
		if err != nil {
			failures = append(failures, fmt.Sprintf("Error evaluating %s %s %s: %v", exp.Field, exp.Operator, exp.Value, err))
		} else if !result {
			failures = append(failures, fmt.Sprintf("%s %s %s failed", exp.Field, exp.Operator, exp.Value))
		}
	}

	return failures, nil
}

// evaluateExpectation evaluates a single expectation
func (e *Evaluator) evaluateExpectation(exp Expectation, jsonBody string) (bool, error) {
	// Extract field value using GJSON. Spec paths may start with "$" or "$."; strip for GJSON.
	path := exp.Field
	if strings.HasPrefix(path, "$.") {
		path = strings.TrimPrefix(path, "$.")
	} else if strings.HasPrefix(path, "$") {
		path = strings.TrimPrefix(path, "$")
	}
	fieldValue := gjson.Get(jsonBody, path)

	// Handle the value based on type
	switch exp.Value {
	case "null":
		return !fieldValue.Exists(), nil
	case "undefined":
		return !fieldValue.Exists(), nil
	default:
		return e.compare(fieldValue, exp.Value, exp.Operator)
	}
}

// compare compares two values using the given operator
func (e *Evaluator) compare(fieldValue gjson.Result, value string, operator string) (bool, error) {
	switch operator {
	case "==":
		return e.equalsValue(fieldValue, value), nil
	case "!=":
		return !e.equalsValue(fieldValue, value), nil
	case "<":
		return e.lessThan(fieldValue, value)
	case ">":
		return e.greaterThan(fieldValue, value)
	case "<=":
		return e.lessOrEqual(fieldValue, value)
	case ">=":
		return e.greaterOrEqual(fieldValue, value)
	case "=~":
		return e.matchesRegex(fieldValue, value)
	case "!~":
		matched, err := e.matchesRegex(fieldValue, value)
		if err != nil {
			return false, err
		}
		return !matched, nil
	default:
		return false, fmt.Errorf("unsupported operator: %s", operator)
	}
}

// matchesRegex checks if the field value matches the given regex pattern.
func (e *Evaluator) matchesRegex(fieldValue gjson.Result, value string) (bool, error) {
	pattern := strings.TrimSpace(value)
	if strings.HasPrefix(pattern, "/") && strings.HasSuffix(pattern, "/") && len(pattern) >= 2 {
		pattern = strings.TrimPrefix(pattern, "/")
		pattern = strings.TrimSuffix(pattern, "/")
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, fmt.Errorf("invalid regex pattern: %w", err)
	}

	return re.MatchString(fieldValue.String()), nil
}

// equalsValue checks if the field value equals the given value
func (e *Evaluator) equalsValue(fieldValue gjson.Result, value string) bool {
	// Handle string values (quoted)
	if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
		expectedStr := strings.Trim(value, `"`)
		return fieldValue.String() == expectedStr
	}

	// Handle null
	if value == "null" {
		return !fieldValue.Exists() || fieldValue.Type.String() == "null"
	}

	// Handle numbers
	if num, err := strconv.ParseFloat(value, 64); err == nil {
		return fieldValue.Num == num
	}

	// Handle booleans
	if value == "true" || value == "false" {
		return fieldValue.Bool() == (value == "true")
	}

	return fieldValue.String() == value
}

// lessThan checks if field < value (numeric comparison)
func (e *Evaluator) lessThan(fieldValue gjson.Result, value string) (bool, error) {
	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return false, fmt.Errorf("invalid number: %s", value)
	}
	return fieldValue.Num < num, nil
}

// greaterThan checks if field > value (numeric comparison)
func (e *Evaluator) greaterThan(fieldValue gjson.Result, value string) (bool, error) {
	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return false, fmt.Errorf("invalid number: %s", value)
	}
	return fieldValue.Num > num, nil
}

// lessOrEqual checks if field <= value (numeric comparison)
func (e *Evaluator) lessOrEqual(fieldValue gjson.Result, value string) (bool, error) {
	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return false, fmt.Errorf("invalid number: %s", value)
	}
	return fieldValue.Num <= num, nil
}

// greaterOrEqual checks if field >= value (numeric comparison)
func (e *Evaluator) greaterOrEqual(fieldValue gjson.Result, value string) (bool, error) {
	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return false, fmt.Errorf("invalid number: %s", value)
	}
	return fieldValue.Num >= num, nil
}
