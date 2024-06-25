
```markdown
# This is the test title which is an optional

This is the test description which is also an optional.

```http
POST https://httpbin.org/post HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 25

key1=value1&key2=value2
```

* `Body.form` is defined
```

## HTTP Request

## Testcases

Testcase is a markdown file, containing one http request and at least one assertion.

The first heading level 1 is used as the testcase title. The paragraph right after the title is used as the testcase description.

### Example of testcase markdown

~~~markdown
# Test title

This is test description.

```http
POST https://httpbin.org/post HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 25

key1=value1&key2=value2
```

* `Body.form` is defined
* `Body.form.key1`==`"value1"`
* `Body.form.key2`==`"value2"`
~~~

## Assertions

Assertion is unordered list item, written in the syntax `<selector><binary-assertion-operator><expected-value>` or `<selector><unary-assertion-operator>`.

### Binary assertion operators

Binary assertion operator is assertion operator that compares value from the given `selector` to the `expected value`.

#### Supported binary opertors

* `==`, equality assertion, compatible with number, string, regular expression.
* `!=`, inequality assertion, compatible with number, string, regular expression.
* `>`, `>=`, `<`, `<=`, inequality assertion, compatible with number.

### Unary assertion operators

Unary assertion operator is assertion operator that verifies value from the given `selector` is in the specific state.

#### Supported unary operators

* `is defined`, `is undefined` to see if the value presents or not.

## Executables

Executable requires only one http request. Other than that are optional. Executable is used for making a request before or after the test (See [test setups](#test-setups) and [test teardowns](#test-teardowns)).

**_Note that an executable with assertions will be considered a testcase._**

## Test setups

Test setups is a list of links to executable markdown. The list can be placed anywhere in the testcase but must be BEFORE the testcase request. Those executables that will be executed before the request of the testcase is made.

The list will be written in the following syntax.

```markdown
* [<executable-title>](<link-to-executable-markdown-file>)
```

## Test teardowns

Test teardown looks like the test setup list, but is placed AFTER the testcase request. The executables in test teardown list will be execute after the request of the testcase has been made.

## Captures

Values in responses from testcase and executable can be captured into variables for further reuses. The presence of variable will be last long until the testcase ended.