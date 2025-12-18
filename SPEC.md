# Summary

A test runner that can execute API spec documents written in Markdown format.

## Versioning

This app uses semantic versioning convention, SemVer.

## Markdown spec convention

User can write the API document in Markdown format with the following convention.

* The first heading (#) will be test case title.
* The content in the first code fence will be HTTP request body.
* The links to markdown file somewhere before the code fence will be considered test setup and will be executed in order sequentially.
* The links to markdown file somewhere after the code fence will be considered test teardown and will be executed in order sequentially.
* Expectations are listed as bullets, somewhere after code fence. See [Expectation](#expectation) for details.
* Variable declarations are listed as bullets, somewhere before code fence. See [Variables](#variables) for details.
* The markdown file which no expectation detected will be considered an executable, not a test case, and will be executed only if it is linked from a test case.
* The variables declared in setups, test case, teardowns are shared among markdowns. If a test case has 2 setup markdown files, called file A and B, and file A declares a variable var_a and file B declares a variable var_b, file B will be able to refer to var_a, and test case can refer to both vars.
* Supported code fence are triple backticks with or without language annotation, triple tildes with or without language annotation.

## Internal architecture

* JSON paths in this document comply with GJSON syntax, https://github.com/tidwall/gjson.
* JSON path MUST start with `$`.
* Project structure MUST emphasize on modularity, focusing on high cohesion and low coupling.
* Modular packages are required to separate the concern between main.go file (orchestrator) and other packages.
* Corresponding unit testing test cases must be add to covered the crucial and high risk test scenarios.

## Expectation

The expectation format syntax:

```
`<field>` <binary comparison operator> `<value>`
```

* Field is a JSON path, surrounded by backticks. For example, `$.data.access_token` or `$.data.records[0].name`.
* Value is an immediate value, surrounded by backticks. For example, `0` for 0 as a number, `"John"` for "John" as a string literal, `null` for null, `undefined` for undefined, `/The*/` for regular expression.
* Comparison operators:
  * Number operators: < (less than), > (greater than), <= (less than or equal to), >= (greater than or equal to), == (equal to), != (not equal to)
  * String operators: == (equal to), != (not equal to)
  * Regular expression operators: =~ (match), !~ (not match)
* Examples:
  * `$.data.records[0].name` == `"John"`
  * `$.data.records` != undefined
  * `$.data.records[0].position` =~ `/^Executive/`

## Variables

The variable declaration syntax:

```
`<variable name>` = `<value>`
```

* Variable name is a string.
* Value can be immediate value or JSON path (prefixed with $).
  * If value is immediate value, the value is assigned to the variable.
  * If value is JSON path, the value will be retrieved from the most recent HTTP response.
* Duplicate variable declaration will overwrite the existing one.
* Examples:
  * `product_id` = `3`
  * `created_by` = `"employee"`
  * `product_sku` = `$.data.sku`

## Use of variable

The variable can be referred by `{{ }}`, double curly braces, and reused many times in request, or other variable declaration. For example:

* `product_id` = `3`
* `path` = `"/products/{{ product_id }}"` 
* In request,
  ```
    POST {{ path }} HTTP/1.1
    Content-Type: application/json

    product_id={{ product_id }}&value=test
  ```

## CLI Usage

### Execute a test

```bash
cotton <filename>
```

For example: `cotton ./examples/example.spec.md`

### Execute all tests in a path

```bash
cotton <path>
```

For example: `cotton ./tests`

### Execute test with wildcard

```bash
cotton <file-spec>
```

For example: `cotton ./test/*.spec.md`

### Fail fast, Stop test execution on first failure

Use `-f` flag for fail fast execution.

### Debug log

Default items shown in standard log:

* Test case title.
* Test result.
* Tests execution summary.

To turn on debug logging, use `-d` flag. Items added into the standard log:

* Setup and teardown titles.
* Request/response in all operations.
* Variables generations.
* Expectation matchings.

### Silence execution

Use `-q` flag for silence execution. Exit code is 0 for success, otherwise 1. This is for CI/CD. In silence execution, debug logging and monochrome logging will be ignored.

### Monochrome logging

Use `-m` flag for monochrome logging.

### Custom timeout

Use `--timeout <timeout in millisecond>` for custom timeout duration.

### Report Generation

* Use `--html-report` for html reporting.
* Use `--json-report` for json reporting.

## Report

### HTML Report

TODO: Will be implemented later

### JSON Report

Items to be included in report:

* Started time: Epoch, milliseconds
* Finished time: Epoch, milliseconds
* Time ellapsed: Milliseconds
* Number of test executed: number
* Number of test passed: number
* Number of test failed: number
* List of test cases
* Each test case in list composes of:
  * File path
  * Test title
  * Test result
  * Started time: Epoch, milliseconds
  * Finished time: Epoch, milliseconds
  * Time ellapsed: Milliseconds

## Multiple flags execution

The flags arguments can be used together, as long as they are not conflict. For example:

```bash
cotton --html-report --json-report -m -d ./examples
```

The command above shows debug logging in monochrome and yield both HTML report and JSON report.

## Showing app version

Use `-v` flags. This flag is priority. If this flag is set, just show the app version in the following format and exit the app with code 0.

```
cotton v1.0.0
```