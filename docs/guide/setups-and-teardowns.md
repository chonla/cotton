---
title: Setups & Teardowns
layout: default
nav_order: 2
permalink: /guide/setups-and-teardowns
parent: Guide
---

# Setups and Teardowns

## Multiple setups and teardowns

You can put setups and teardowns as many as you want, as long as the setups are before the request and the teardowns are after the request.

**Example**

{% highlight markdown %}
# List Products

* [Sign in](<rootDir>/etc/examples/fakestoreapi.com/executables/auth.md)
* [Create data session](<rootDir>/etc/examples/fakestoreapi.com/executables/create_session.md)

```http
GET https://fakestoreapi.com/products HTTP/1.1
Authorization: Bearer {{access_token}}
```

* `Body.0.id`==`1`

* [Release data session](<rootDir>/etc/examples/fakestoreapi.com/executables/release_session.md)
* [Sign out](<rootDir>/etc/examples/fakestoreapi.com/executables/auth.md)
{% endhighlight %}

## rootDir Definition

`<rootDir>` is a shortcut variable allowing you to make link more dynamic.

You can redefine `<rootDir>` via [command line](../cli.md) option `-r`.

## Nested setups and teardowns

Setups and teardowns cannot be nested. That means in setup and teardown files, declaration of nested setups and teardowns will be ignored.