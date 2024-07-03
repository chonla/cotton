---
title: Selectors
layout: default
nav_order: 8
permalink: /syntax/selectors
parent: Syntax
---

# Selectors

A selector is a reference to a value in a response. Once a response has been received, the values in the response are referrable.

## Response Status Line

The values in the response status line (the very first line of the response) can be referred by the following values:

| Selector | Refer to |
|---|---|
| StatusCode | Response status code, e.g. 200 |
| StatusText | Response status text, e.g. OK |
| Status | Response status, e.g. 200 OK |
| Protocol | Response protocol, e.g. HTTP |
| Version | Response protocol version, e.g. 1.1 |

## Response Headers

The values in the response headers can be referred by prefixing the header name with `Headers.`, for example, `Headers.Content-Type`.

## Response Body

The response body must be in JSON. The values can be referred by prefixed the json selector with `Body.`, for example, `Body.data.id`.

**Example**

Givien the responded JSON:

```json
{
    "data": {
        "id": 3
    }
}
```

The id in the response can be referred with `Body.data.id`.

**Example**

Givien the responded JSON:

```json
[
    {
        "data": {
            "id": 3
        }
    },
    {
        "data": {
            "id": 5
        }
    },
]
```

The id in the second element in the response can be referred with `Body[1].data.id`.

