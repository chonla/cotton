# Get data from API

This can be run with `cotton -u http://httpbin.org httpbin.md`

## GET /get?key1=value1

## Expectation

| Assert | Expected |
| - | - |
| StatusCode | 200 |
| Data.args.key1 | value1 |

# Post data to API

## POST /post

Data in header

| Header | Value |
| - | - |
| Content-Type | application/json |

Data in body/payload of HTTP Request

```
{
    "key1": "value1"
}
```
## Expectation

| Assert | Expected |
| - | - |
| StatusCode | 200 |
| Data.json.key1 | value1 |
