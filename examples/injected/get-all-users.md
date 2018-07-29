# Get All Users

To pass this test, token must be passed from command line. Try `cotton -p token=<token>`.

## GET /users

| Header | Value |
| - | - |
| Authorization | Bearer {token} |

## Expectation

Here is the expectation table

| Assert | Expected |
| - | - |
| StatusCode | 200 |
| Header.Content-Type | application/json;charset=utf-8 |
| Data.page_count | 1 |
| Data.page | 1 |

