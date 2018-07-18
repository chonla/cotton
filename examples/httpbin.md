# Get data from API 

This can be run with `cotton -u http://httpbin.org httpbin.md`

## GET /get?key1=value1

## Expectation

| Assert | Expected |
| - | - |
| StatusCode | 200 |
| Data.args.key1 | value1 |


