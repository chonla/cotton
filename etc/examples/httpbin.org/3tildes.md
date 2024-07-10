# Test GET on httpbin.org with three-tilded code block

Test getting data from httpbin.org using multiple http requests.

## Test setting up

* [Post some data to host](etc/examples/httpbin.org/post.md)

~~~http
GET https://httpbin.org/get?key1=value1&key2=value2 HTTP/1.1
~~~

* [Patch some data to host](etc/examples/httpbin.org/patch.md)

## Assertions

* `Body.args.key1`==`"value1"`
* `Body.args.key2`==`"value2"`