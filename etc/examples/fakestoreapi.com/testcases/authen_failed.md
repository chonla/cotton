# Authentication with invalid user

* username:"invalid"
* password:"invalid"

```http
POST https://fakestoreapi.com/auth/login HTTP/1.1
Content-Type: application/json
Content-Length: 43

{"username":"{{username}}","password":"{{password}}"}
```

* `StatusCode` == `401`
* `Body.token` is undefined