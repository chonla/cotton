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

## Code fences

Both triple backticks and triple tildes are supported, with or without a language tag:

{% highlight markdown %}
```http
GET https://api.example.com/ping HTTP/1.1
```

~~~http
GET https://api.example.com/ping HTTP/1.1
~~~
{% endhighlight %}

## Variable substitution

You can insert variables using `{{ name }}` placeholders. Missing variables cause a ParseError and abort execution.

{% highlight markdown %}
* `base_url` = `"https://api.example.com"`

```http
GET {{ base_url }}/products/42 HTTP/1.1
Accept: application/json
```
{% endhighlight %}