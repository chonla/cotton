# Actions

Action is a request sent to RESTful API application. Use `##` to define an action in a test case.

## Request Body

Some action may need to send a body to the server, like `POST`. Use ~```~ to declare the request body.

## Request Headers

Header of the request can be defined in a 2-columned table with title of `Header` and `Value`.

## Files Upload

Use bullet list for perform files upload. File name must be stated in anchor link. Field name is corresponding label.

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

## File Upload Example

~~~
## POST /media/photo

* [photo[]](./examples/upload/3x3.gif)
* [photo[]](./examples/upload/lamps.png)
~~~

The above upload example is equivalent to the following HTTP request:

```
POST /media/photo HTTP/1.1
Content-Length: 23064
Content-Type: multipart/form-data; boundary=5a7129a40f523a00b69b870968079715a91d48c8cab50d938df7b9dae25e

--5a7129a40f523a00b69b870968079715a91d48c8cab50d938df7b9dae25e
Content-Disposition: form-data; name="photo[]"; filename="3x3.gif"
Content-Type: application/octet-stream

... 3x3.gif binary content ...
--5a7129a40f523a00b69b870968079715a91d48c8cab50d938df7b9dae25e--
--5a7129a40f523a00b69b870968079715a91d48c8cab50d938df7b9dae25e
Content-Disposition: form-data; name="photo[]"; filename="lamps.png"
Content-Type: application/octet-stream

... lamps binary content ...
--5a7129a40f523a00b69b870968079715a91d48c8cab50d938df7b9dae25e--
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