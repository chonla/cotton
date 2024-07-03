---
title: Requests
layout: default
nav_order: 2
permalink: /syntax/requests
parent: Syntax
---

# Requests

Request is a HTTP request in the code block annotated language with `http`.

{% highlight markdown %}
```http
POST https://somedomain.com/resource HTTP/1.1
Content-Length: 17
Content-Type: application/json

{"resource":8839}
```
{% endhighlight %}