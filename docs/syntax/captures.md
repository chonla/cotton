---
title: Captures
layout: default
nav_order: 7
permalink: /syntax/captures
parent: Syntax
---

# Captures

Capturing a value is one way to define a [variable](./variables.md) based on a JSON response.

In Cotton, captures are declared the same way as variables — as list items before the request — but with the value set to a JSON path (starting with `$`). The actual value will be resolved from the most recent JSON response after the request runs.

Syntax:

{% highlight markdown %}
* `variable_name` = `$.path.to.value`
{% endhighlight %}

Example:

{% highlight markdown %}
* `token` = `$.data.access_token`

```http
POST https://api.example.com/auth HTTP/1.1
Content-Type: application/json

{"username":"user","password":"pass"}
```
{% endhighlight %}

After the request completes, `token` will be captured from the JSON response and becomes available for subsequent steps (e.g., teardowns or later specs).

{: .info }
- Captures are resolved from the current spec’s response. If a previous step (e.g., a setup) produced a response, JSON-path variables may also resolve before the main request.
- The response body must be valid JSON for captures to work.