# Yas

Yet Another Silk Test Tool. This project is originally inspired by [silk test](https://github.com/matryer/silk).

* Markdown based document-driven RESTful API testing.
* This tool is not compatible with silk markdown.

## Test Suite

Test Suite is a markdown file. Test Suite can contain several test cases.

## Test Case

Test Case uses H1 (#) to identify beginning of test case and its name.

```
# Test case name
```

## Method and URI

Method and URI uses H2 (##) to identify method and URI.

```
## POST /login
```

## Request headers

Request headers will be in a table with 2 columns named "Header" and "Value".

```
| Header | Value |
| - | - |
| Content-Type | application/json |
```

## Request body

Request body uses code block (3 backticks).

~~~
```
{
    "login": "admin",
    "password": "secret"
}
```
~~~

## Expectation

Expectation section starts by using H2 (##) with title "Expectation".

```
## Expectation
```

## Assertion

Assertion uses table with 2 columns. The header must be "Assert" and "Expected".

* Value in response header can be addressed by using "Header" object, e.g. `Header.Content-Type`.
* Value in response body can be addressd by using "Data" object, e.g. `Data.token`.
* StatusCode is response status code.

## Comment

All lines not mentioned above are comment.

## License

[License](LICENSE.txt)