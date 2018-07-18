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

| Previous | Index | Next |
| :-: | :-: | :-: |
| [Testcases](Testcases.md) | [Top](README.md) | [Assertions](Assertions.md) |