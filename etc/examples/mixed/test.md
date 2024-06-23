# Test getting version from opensource.org and post to HTTPBin.org

* [Get copyleft version from opensource.org](<rootDir>/etc/examples/mixed/before.md)

```http
POST https://httpbin.org/post HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 25

license={{license}}
```

* `Body.form.license`==`"GPL-3.0"`