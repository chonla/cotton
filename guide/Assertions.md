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
| Cookies.&lt;cookie-name&gt; | Value in cookie corresponding to the `cookie-name` given |
| Data.&lt;data-name&gt; | Data in response body, available only when data is JSON (data must be JSON parsable and content-type must be `application/json`) |

## Pattern Matching

Cotton use regular expression to do complicated assertion. Regular expression pattern is surrounded by `/`.

### Pattern Matching Example

The assertion below tests the returned status code must be `201` and content type must be started with `application/json`.

| Assert | Expected |
| - | - |
| StatusCode | 201 |
| Header.Content-Type | /^application/json($|;)/ |

## Predefined Expectation

Cotton provide some predefined expectation to help assertion. The following expected value can be used as to assert some response value.

| Predefined value | Description |
| - | - |
| `*should exist*` | The asserted variable should present in the response. |
| `*should not exist*` | The asserted variable should not present in the response. |
| `*should be null*` | The asserted variable should be null. |
| `*should not be null*` | The asserted variable should not be null. |
| `*should be true*` | The asserted variable should be boolean with value TRUE. |
| `*should be false*` | The asserted variable should be boolean with value FALSE. |

Underscores can also be used instead of stars, e.g.: `_should exist_`

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