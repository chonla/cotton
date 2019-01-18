# Get into lobby

## POST /lobby

| Header | Value |
| - | - |
| Content-Type | application/json |

```
{
    "name": "chonla"
}
```

## Expectations

| Assert | Expected |
| - | - |
| StatusCode | 200 |
| Header.Set-Cookie | *should exist* |
| Header.Set-Cookie | /\bsession_id=.+;/ |