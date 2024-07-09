---
title: Selectors
layout: default
nav_order: 3
permalink: /guide/selectors
parent: Guide
---

# Selectors

Cotton utilizes selector from [GJSON](https://github.com/tidwall/gjson/blob/master/SYNTAX.md).

## GJSON Playground

You can try GJSON online at [https://gjson.dev/](https://gjson.dev/).

## Quick Syntax

* Path separator with `.`.
* Wildcards (`?`, `*`) are supported.
* Characters excaping with `\` then `.`, `?`, and `*` are escaped as `\.`, `\?` and `\*`.
* Getting array length with `#`, e.g. `friends.#` to get length of friends
* Getting nth element with `.<index>`, e.g. `friends.0.age` to get the first friend's age.

For more detail, see [GJSON Syntax](https://github.com/tidwall/gjson/blob/master/SYNTAX.md).