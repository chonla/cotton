---
title: Selectors
layout: default
nav_order: 8
permalink: /syntax/selectors
parent: Syntax
---

# Selectors

A selector identifies a value in the JSON response body using a JSON path. Cotton uses [GJSON](https://github.com/tidwall/gjson) syntax.

## Basics

- Paths must start with `$` (e.g., `$.data.id`).
- Dot notation selects object fields.
- Array indices are zero-based (e.g., `$.items.0.name`).

Example response:

```json
{
  "status": "ok",
  "data": {
    "id": 3,
    "items": [
      { "name": "A" },
      { "name": "B" }
    ]
  }
}
```

Selectable values:

- `$.status` → "ok"
- `$.data.id` → 3
- `$.data.items.0.name` → "A"

Advanced examples (GJSON):

- `$.data.items.#` → array length
- `$.data.items.#.name` → all names
- `$.data.items.1` → second element

{: .warning }
Assertions and captures operate on the JSON body. Status line and headers are not available as selectors.

