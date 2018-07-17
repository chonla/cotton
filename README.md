# Cotton

Markdown Test Specification Runner. This project is originally inspired by [silk test](https://github.com/matryer/silk).

* Executable RESTful API Markdown-based Document Specification.

## Usage

```
cotton -u <base-api-url> [directory or file]
```

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
* Expected value will be treated as regular expression when surrounded by /.

## Test Setup

Test setup section starts by using H2 (##) with title "Precondition".

Steps in test setup are written using bullets with links. Title of link are just a text explaining what to be done. The link is an actual task to be performed.

```
* [Login](./library/login.md)
* [Get User Data](./library/get-user-data.md)
```

You can put anything you want surrounding the link to make the step clearer.

```
* [Login](./library/login.md) with valid credential and collect access token.
* Then [Get User Data](./library/get-user-data.md) with the token above.
```

## Test Teardown

Test setup section starts by using H2 (##) with title "Finally".

Declaring steps in test teardown is like doing so in test setup.

## Tasks

Tasks are used in Setup and Teardown. Task does not contain any assertion. Capture table can be declared in Task for future use.

## Capture Table

Capture Table is a table telling what to be collected. It has 2 columns with header "Name" and "Value". Name is referrable name. Value is what to be kept. Value syntax is identical to [Assertion](#assertion).

## Comment

All lines not mentioned above are comment.

## License

[License](LICENSE.txt)