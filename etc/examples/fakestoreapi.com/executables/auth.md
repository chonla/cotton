* username:"mor_2314"
* password:"83r5^_"

```http
POST https://fakestoreapi.com/auth/login HTTP/1.1
Content-Type: application/json
Content-Length: 43

{"username":"{{username}}","password":"{{password}}"}
```

* access_token:`Body.token`