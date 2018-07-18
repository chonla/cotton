# Google search

This can be run with `cotton -u https://www.google.co.th google.md`

## GET /search?q=cotton

Perform a search for the word `cotton`.

## Expectation

| Assert | Expected |
| - | - |
| StatusCode | 200 |
| Header.Server | gws |