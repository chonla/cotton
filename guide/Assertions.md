# Assertions

Assertions is about to assert the response with the expected values. Use `## Expectation` or `## Expectations` to tell Cotton this section is assertions.

The expected values are defined in a 2-columned table with title of `Assert` and `Expected`.

## Asserted Variables

All variable names are case-insensitive.

| Variable Name | Description |
| - | - |
| StatusCode | HTTP Response Status Code, e.g. `200`, `404` |
| Status | HTTP Response Status, e.g. `200 OK` |
| Header.&lt;header-name&gt; | Value in header corresponding to the `header-name` given |
| Data.&lt;data-name&gt; | Data in response body, available only when data is JSON (data must be JSON parsable and content-type must be `application/json`) |

## Example

```
## Expectations

| Assert | Expected |
| - | - |
| StatusCode | 200 |
| Header.Content-Type | application/json |
| Data.title | Buy milk |
```

| Previous | Index | Next |
| :-: | :-: | :-: |
| [Actions](Actions.md) | [Top](README.md) | [Test Setups](TestSetups.md) |