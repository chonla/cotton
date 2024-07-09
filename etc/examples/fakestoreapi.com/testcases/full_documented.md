# Full documented testcase

The testcase is described by providing paragraphs right after the test case title.

## Before Test

The header `Before Test` above is an optional and written there for understandability purpose. All links defined as list item before HTTP request are considered test setting ups.

* [Create a new user for authentication](<rootDir>/etc/examples/fakestoreapi.com/executables/add_user.md)
* [Authentication with the new user](<rootDir>/etc/examples/fakestoreapi.com/executables/auth.md)

## Request

The header `Request` above is also an optional. The HTTP request is described by code block with `http` annotation right after the open code block ` ```http `.

Cotton will consider only the first code block as the testing request.

```http
GET https://fakestoreapi.com/products HTTP/1.1
Authorization: Bearer {{access_token}}
```

## After Test

The header `After Test` above is also an optional. All links defined as list after code block ` ```http` will be executed after test request.

* [Delete test user](<rootDir>/etc/examples/fakestoreapi.com/executables/delete_user.md)

## Assertions

Assertion defined as list in `<value><operator><expected>` format. All LHS variable must be defined in inline code block. All expected values must be written in inline code block. Type must be explicitly defined.

* `StatusCode` == `200`
* `Body.#` > `0`
* `Body.0.id` == `1`

**Note:** You can find more detail on syntax on [Guide](https://chonla.github.io/cotton).