# Get All Users

## Preconditions

* [Login With Valid Credential](common/login.md)

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

