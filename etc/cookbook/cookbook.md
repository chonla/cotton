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