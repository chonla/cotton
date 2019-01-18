# Create Room

## Preconditions

* [Get into lobby](./common/get-into-lobby.md)

## POST /room

```
{
    "password": "A"
}
```

## Expectations

| Assert | Expected |
| - | - |
| StatusCode | 201 |
| Header.Locattion | /.+/ |