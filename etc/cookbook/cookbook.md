# Cotton API reference

## Minimum testcase

A testcase requires one http request and at least one assertion. The following testcase is an example of an untitled testcase with minimum requirement.

~~~markdown
```http
POST https://httpbin.org/post HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 25

key1=value1&key2=value2
```

* `Body.form` is defined
~~~

## A descriptive testcase

The heading level 1 and the following paragraph will be considered as the testcase title and description.

~~~markdown
# Sending post to httpbin.org

This request is an example for sending a request to httpbin.org.

```http
POST https://httpbin.org/post HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 25

key1=value1&key2=value2
```

* `Body.form` is defined
~~~

