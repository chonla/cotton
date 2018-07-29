# Login Should Return Token

## POST /login

Login with a valid credential.

| Header | Value |
| - | - |
| Content-Type | application/json |

```
{
    "login": "admin",
    "pwd": "admin"
}
```

## Expectation

Here is the expectation table

### Example Response

```
{
    "token": "1234"
}
```

| Assert | Expected |
| - | - |
| StatusCode | 200 |
| header.Content-Type | /^application/json($|;)/ |
| Data.token | /.+/ |