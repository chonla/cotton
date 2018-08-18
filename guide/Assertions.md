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

## Predefined Expectation

Cotton provide some predefined expectation to help assertion. The following expected value can be used as to assert some response value.

| Predefined value | Description |
| - | - |
| `*should exist*` | The asserted variable should present in the response. |
| `*should not exist*` | The asserted variable should not present in the response. |

## Example

```
## Expectations

| Assert | Expected |
| - | - |
| StatusCode | 200 |
| Header.Content-Type | application/json |
| Data.title | Buy milk |
| Data.credential | *should not exist* |
```

| Previous | Index | Next |
| :-: | :-: | :-: |
| [Actions](Actions.md) | [Top](README.md) | [Test Setups](TestSetups.md) |