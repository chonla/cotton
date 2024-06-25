# Cotton API reference

## Minimum testcase

Testcase markdown is a markdown file that contains a http request to test if the API satisfies requirements.

A minimum testcase requires one http request and at least one assertion.

~~~markdown
```http
POST https://httpbin.org/post HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 25

key1=value1&key2=value2
```

* `Body.form` is defined
* `Body.form.key1`==`"value1"`
~~~

## Minimum executable

Executable markdown is a markdown file that contains a http request to perform some operation.

The minimum requirement of executable is just one http request, NO assertion.

~~~markdown
```http
POST https://httpbin.org/post HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 25

key1=value1&key2=value2
```
~~~

## A descriptive testcase

The heading level 1 and the following paragraph will be considered as the testcase title and description.

~~~markdown
# Sending post to httpbin.org

This request is an example of sending a request to httpbin.org.

```http
POST https://httpbin.org/post HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 25

key1=value1&key2=value2
```

* `Body.form` is defined
* `Body.form.key1`==`"value1"`
~~~

## Test setup

Test setup is a list of link to executable markdowns and placed anywhere before http request in the testcase.

~~~markdown
# Sending post to httpbin.org

This request is an example of calling setting ups and sending a request to httpbin.org.

* [Authenticate client](./common/auth.md)
* [Get session reference id](./common/get_ref_id.md)

```http
POST https://httpbin.org/post HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 25

key1=value1&key2=value2
```

* `Body.form` is defined
* `Body.form.key1`==`"value1"`
~~~

## Test teardown

Test teardown is a list of link to executable markdowns like test setup but placed anywhere after http request in the testcase.

~~~markdown
# Sending post to httpbin.org

This request is an example of calling tearing downs and sending a request to httpbin.org.

```http
POST https://httpbin.org/post HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 25

key1=value1&key2=value2
```

* [Reset database](./common/reset_db.md)
* [Clear cache](./common/reset_cache.md)

* `Body.form` is defined
* `Body.form.key1`==`"value1"`
~~~

## Capturing result

The result of any request can be captured into variables. The variables can be referred in the following requests once they are captured.

From the following markdown:

~~~markdown
# Sending post to httpbin.org

This request is an example of calling setting ups and sending a request to httpbin.org.

* [Authenticate client](./common/auth.md)
* [Get session reference id](./common/get_ref_id.md)

```http
POST https://httpbin.org/post HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 25

key1=value1&key2=value2
```

* `Body.form` is defined
* `Body.form.key1`==`"value1"`
~~~

Given a url `http://sso.someurl.org/ref` returns the following JSON when received a POST request:

```json
{"ref_id":"1726-ed7a7"}
```

The markdown `get_ref_id.md` will look like this.

~~~markdown
```http
POST http://sso.someurl.org/ref HTTP/1.1
Content-Type: application/json
Content-Length: 20

{"client-id":"1234"}
```

* `refid`:`Body.ref_id`
~~~

where `` `refid` ``:`` `Body.ref_id` `` is capturing syntax. The `refid` is variable name and `Body.ref_id` is selector from result.

The modified testcase will look like this.

~~~markdown
# Sending post to httpbin.org

This request is an example of calling setting ups and sending a request to httpbin.org.

* [Authenticate client](./common/auth.md)
* [Get session reference id](./common/get_ref_id.md)

```http
POST https://httpbin.org/post HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 25

key1=value1&key2=value2&refid={{refid}}
```

* `Body.form` is defined
* `Body.form.key1`==`"value1"`
* `Body.form.refid`==`"1234"`
~~~