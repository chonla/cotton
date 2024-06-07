# Test POST on httpbin.org

```http
POST https://httpbin.org/post HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 25

secret=thisIsASecretValue
```

## Capture something from result

* secret:$.form.secret