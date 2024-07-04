---
title: Requests
layout: default
nav_order: 2
permalink: /syntax/requests
parent: Syntax
---

# Requests

A request is a HTTP request in the code block annotated language with `http`.

**Example**

{% highlight markdown %}
```http
POST https://fakestoreapi.com/products HTTP/1.1
Content-Length: 125
Content-Type: application/json

{"title":"test product","price":13.5,"description":"lorem ipsum set","image":"https://i.pravatar.cc","category":"electronic"}
```
{% endhighlight %}