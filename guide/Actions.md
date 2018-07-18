# Actions

Action is a request sent to RESTful API application. Use `##` to define an action in a test case.

## Request Body

Some action may need to send a body to the server, like `POST`. Use ~```~ to declare the request body.

## Request Headers

Header of the request can be defined in a 2-columned table with title of `Header` and `Value`.

## Example

~~~
## POST /todos

| Header | Value |
| - | - |
| Content-Type | application/json |
| Authorization | Bearer some-token |

```
{
    "title": "Buy milk"
}
```
~~~

The above example is equivalent to the following HTTP request:

```
POST /todos HTTP/1.0
Content-Type: application/json
Authorization: Bearer some-token

{
    "title": "Buy milk"
}
```

## Supported methods

Cotton provides the following HTTP methods:

* GET
* POST
* PUT
* PATCH
* DELETE
* OPTION

## Assertions of the action

Test cases require [assertions](Assertions.md) to verify the response from the API, but not all actions requires assertions.

## Actions without assertions

All actions without assertions will not be executed. Action without assertions can be treated as [Test Setups](TestSetups.md) or [Test Teardowns](TestTeardowns.md) actions.

| Previous | Index | Next |
| :-: | :-: | :-: |
| [Testcases](Testcases.md) | [Top](README.md) | [Assertions](Assertions.md) |