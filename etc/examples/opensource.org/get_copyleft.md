# Test getting copyleft from opensource.org

```http
GET https://api.opensource.org/licenses/copyleft HTTP/1.1
```

* `Body.1.id`==`"GPL-3.0"`
* `Body.1.superseded_by`==`null`
* `Body.1.superseded_date` is undefined