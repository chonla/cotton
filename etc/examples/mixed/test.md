# Test getting version from opensource.org and post to HTTPBin.org

* [Get copyleft version from opensource.org](<rootDir>/etc/examples/mixed/before1.md)
* [Get "world" translation](<rootDir>/etc/examples/mixed/before2.md)

```http
POST https://httpbin.org/post HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 25

license={{license}}&word={{word}}
```

* `Body.form.license`==`"GPL-3.0"`
* `Body.form.word`==`"world"`

* [Get age of Sam](<rootDir>/etc/examples/mixed/after.md)